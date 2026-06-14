package music

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// NeteaseAPIBase points to the NeteaseCloudMusicApi proxy.
// Override with env var: NETEASE_API=http://localhost:3000
var NeteaseAPIBase = func() string {
	if v := os.Getenv("NETEASE_API"); v != "" {
		return v
	}
	return "http://iwenwiki.com:3000"
}()

// GlobalCookie is sent with all authenticated requests (VIP playback).
// Set via NETEASE_COOKIE env var, or overridden per-request at runtime.
var GlobalCookie = os.Getenv("NETEASE_COOKIE")

var httpClient = &http.Client{Timeout: 8 * time.Second}

// RawGet exposes the internal get for handler-level proxying.
func RawGet(path string, params map[string]string) ([]byte, error) {
	return get(path, params)
}

func get(path string, params map[string]string) ([]byte, error) {
	u, _ := url.Parse(NeteaseAPIBase + path)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	resp, err := httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type SearchResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Cover    string `json:"cover"`
	Duration int    `json:"duration"`
}

func Search(keyword string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}
	body, err := get("/search", map[string]string{
		"keywords": keyword,
		"type":     "1",
		"limit":    fmt.Sprintf("%d", limit),
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Result struct {
			Songs []struct {
				ID      int64  `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name string `json:"name"`
				} `json:"album"`
				Duration int `json:"duration"`
			} `json:"songs"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	songs := raw.Result.Songs
	if len(songs) == 0 {
		return []SearchResult{}, nil
	}

	// Batch fetch song detail for cover URLs
	idParts := make([]string, len(songs))
	for i, s := range songs {
		idParts[i] = fmt.Sprintf("%d", s.ID)
	}
	coverMap := fetchCovers(strings.Join(idParts, ","))

	results := make([]SearchResult, 0, len(songs))
	for _, s := range songs {
		artist := ""
		if len(s.Artists) > 0 {
			artist = s.Artists[0].Name
		}
		results = append(results, SearchResult{
			ID:       fmt.Sprintf("%d", s.ID),
			Name:     s.Name,
			Artist:   artist,
			Album:    s.Album.Name,
			Cover:    coverMap[s.ID],
			Duration: s.Duration / 1000,
		})
	}
	return results, nil
}

func fetchCovers(ids string) map[int64]string {
	body, err := get("/song/detail", map[string]string{"ids": ids})
	if err != nil {
		return nil
	}
	var raw struct {
		Songs []struct {
			ID int64 `json:"id"`
			Al struct {
				PicURL string `json:"picUrl"`
			} `json:"al"`
		} `json:"songs"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	m := make(map[int64]string, len(raw.Songs))
	for _, s := range raw.Songs {
		m[s.ID] = s.Al.PicURL
	}
	return m
}

func GetURL(id, cookie string) (string, error) {
	params := map[string]string{"id": id, "br": "320000"}
	// Prefer per-request cookie, fall back to global env cookie
	if cookie == "" {
		cookie = GlobalCookie
	}
	if cookie != "" {
		params["cookie"] = cookie
	}
	body, err := get("/song/url", params)
	if err != nil {
		return "", err
	}
	var raw struct {
		Code int `json:"code"`
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}
	if len(raw.Data) == 0 || raw.Data[0].URL == "" {
		return "", fmt.Errorf("no url for song %s", id)
	}
	return raw.Data[0].URL, nil
}

type LyricLine struct {
	Time float64 `json:"time"` // seconds
	Text string  `json:"text"`
}

func GetLyric(id string) ([]LyricLine, error) {
	body, err := get("/lyric", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}
	var raw struct {
		Lrc struct {
			Lyric string `json:"lyric"`
		} `json:"lrc"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return parseLRC(raw.Lrc.Lyric), nil
}

func parseLRC(lrc string) []LyricLine {
	var lines []LyricLine
	for _, line := range splitLines(lrc) {
		t, text, ok := parseLRCLine(line)
		if ok && text != "" {
			lines = append(lines, LyricLine{Time: t, Text: text})
		}
	}
	return lines
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func parseLRCLine(line string) (float64, string, bool) {
	// Format: [mm:ss.xx]text
	if len(line) < 10 || line[0] != '[' {
		return 0, "", false
	}
	end := -1
	for i := 1; i < len(line); i++ {
		if line[i] == ']' {
			end = i
			break
		}
	}
	if end < 0 {
		return 0, "", false
	}
	timeStr := line[1:end]
	text := line[end+1:]

	var min, sec, cs int
	_, err := fmt.Sscanf(timeStr, "%d:%d.%d", &min, &sec, &cs)
	if err != nil {
		return 0, "", false
	}
	t := float64(min)*60 + float64(sec) + float64(cs)/100.0
	return t, text, true
}
