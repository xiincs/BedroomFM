package api

import (
	"bedroomfm/internal/hub"
	"bedroomfm/internal/models"
	"bedroomfm/internal/music"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func randomCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 5)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// POST /api/room/create
func CreateRoom(c *gin.Context) {
	var req struct {
		RoomName string `json:"roomName"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RoomName == "" || req.Nickname == "" {
		c.JSON(400, gin.H{"error": "roomName and nickname required"})
		return
	}

	roomID := uuid.NewString()
	code := randomCode()
	memberID := uuid.NewString()

	room := models.NewRoom(roomID, code, req.RoomName, memberID)
	host := &models.Member{
		ID:       memberID,
		Nickname: req.Nickname,
		Avatar:   avatarURL(req.Nickname),
		Coins:    10,
		IsHost:   true,
		Persona:  "DJ",
	}
	room.Members[memberID] = host

	hub.GlobalHub.CreateRoom(room)

	c.JSON(200, gin.H{
		"roomId":   roomID,
		"code":     code,
		"memberId": memberID,
	})
}

// POST /api/room/join
func JoinRoom(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" || req.Nickname == "" {
		c.JSON(400, gin.H{"error": "code and nickname required"})
		return
	}

	room := hub.GlobalHub.GetRoomByCode(strings.ToUpper(req.Code))
	if room == nil {
		c.JSON(404, gin.H{"error": "room not found"})
		return
	}

	memberID := uuid.NewString()
	member := &models.Member{
		ID:       memberID,
		Nickname: req.Nickname,
		Avatar:   avatarURL(req.Nickname),
		Coins:    10,
		IsHost:   false,
		Persona:  "听众",
	}
	room.Mu.Lock()
	room.Members[memberID] = member
	room.Mu.Unlock()

	c.JSON(200, gin.H{
		"roomId":   room.ID,
		"memberId": memberID,
		"roomName": room.Name,
		"code":     room.Code,
	})
}

// GET /api/music/login/qr/key
func QRKey(c *gin.Context) {
	body, err := music.RawGet("/login/qr/key", nil)
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", body)
}

// GET /api/music/login/qr/create?key=xxx
func QRCreate(c *gin.Context) {
	key := c.Query("key")
	body, err := music.RawGet("/login/qr/create", map[string]string{"key": key, "qrimg": "true"})
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", body)
}

// GET /api/music/login/qr/check?key=xxx
// Returns {code, message, cookie} — cookie populated when code==803
func QRCheck(c *gin.Context) {
	key := c.Query("key")
	body, err := music.RawGet("/login/qr/check", map[string]string{"key": key})
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	// Forward raw response; frontend handles code 801/802/803/800
	c.Data(200, "application/json", body)
}

// GET /api/music/search?q=keyword
func SearchMusic(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(400, gin.H{"error": "q required"})
		return
	}
	results, err := music.Search(q, 20)
	if err != nil {
		c.JSON(502, gin.H{"error": "music api unavailable: " + err.Error()})
		return
	}
	c.JSON(200, results)
}

// GET /api/music/url?id=xxx[&cookie=MUSIC_U%3Dxxx]
func GetMusicURL(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id required"})
		return
	}
	// cookie passed from client (localStorage) takes priority over server env
	cookie := c.Query("cookie")
	u, err := music.GetURL(id, cookie)
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	if u == "" {
		c.JSON(403, gin.H{"error": "vip_required", "message": "该歌曲需要VIP，请在设置中填入你的网易云Cookie"})
		return
	}
	c.JSON(200, gin.H{"url": u})
}

// GET /api/music/lyric?id=xxx
func GetLyric(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id required"})
		return
	}
	lines, err := music.GetLyric(id)
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, lines)
}

// GET /api/room/:id
func GetRoom(c *gin.Context) {
	room := hub.GlobalHub.GetRoom(c.Param("id"))
	if room == nil {
		c.JSON(404, gin.H{"error": "room not found"})
		return
	}
	c.JSON(200, hub.RoomStatePayload(room))
}

// WS /ws/:roomId?memberId=xxx
func WSHandler(c *gin.Context) {
	roomID := c.Param("roomId")
	memberID := c.Query("memberId")

	room := hub.GlobalHub.GetRoom(roomID)
	if room == nil {
		c.JSON(404, gin.H{"error": "room not found"})
		return
	}

	room.Mu.RLock()
	_, ok := room.Members[memberID]
	room.Mu.RUnlock()
	if !ok {
		c.JSON(403, gin.H{"error": "not a member"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := hub.NewClient(hub.GlobalHub, conn, roomID, memberID)
	hub.GlobalHub.Register(client)

	// Send full room state on connect
	state := hub.RoomStatePayload(room)
	stateJSON, _ := json.Marshal(state)
	initMsg, _ := json.Marshal(hub.WSMessage{
		Type:    "room_state",
		Payload: stateJSON,
	})
	client.SendRaw(initMsg)

	// Notify others
	hub.GlobalHub.BroadcastExcept(roomID, client, hub.WSMessage{
		Type:    "member_join",
		Payload: mustJSON(room.Members[memberID]),
	})
	hub.GlobalHub.Broadcast(roomID, hub.WSMessage{
		Type:    "room_state",
		Payload: stateJSON,
	})

	go client.WritePump()
	client.ReadPump(handleWSMessage)
}

func handleWSMessage(c *hub.Client, msg hub.WSMessage) {
	room := hub.GlobalHub.GetRoom(c.RoomID)
	if room == nil {
		return
	}

	switch msg.Type {
	case "chat":
		var p struct {
			Content string `json:"content"`
		}
		json.Unmarshal(msg.Payload, &p)
		if p.Content == "" {
			return
		}
		room.Mu.RLock()
		member := room.Members[c.MemberID]
		room.Mu.RUnlock()
		if member == nil {
			return
		}
		chatMsg := models.ChatMessage{
			ID:       uuid.NewString(),
			MemberID: c.MemberID,
			Nickname: member.Nickname,
			Content:  p.Content,
			Time:     time.Now().UnixMilli(),
		}
		room.Mu.Lock()
		if len(room.Messages) >= 200 {
			room.Messages = room.Messages[50:]
		}
		room.Messages = append(room.Messages, chatMsg)
		room.Mu.Unlock()
		hub.GlobalHub.Broadcast(c.RoomID, hub.WSMessage{
			Type:    "chat",
			Payload: mustJSON(chatMsg),
		})

	case "queue_add":
		var p struct {
			Song models.Song `json:"song"`
		}
		json.Unmarshal(msg.Payload, &p)
		if p.Song.ID == "" {
			return
		}
		room.Mu.RLock()
		member := room.Members[c.MemberID]
		room.Mu.RUnlock()
		if member == nil {
			return
		}
		item := &models.QueueItem{
			QID:     uuid.NewString(),
			Song:    p.Song,
			AddedBy: c.MemberID,
			Adder:   member.Nickname,
			Votes:   0,
			VoteMap: make(map[string]int),
		}
		room.Mu.Lock()
		wasEmpty := len(room.Queue) == 0 && room.Playback.Song == nil
		room.InsertQueued(item)
		if wasEmpty {
			startNext(room)
		}
		room.Mu.Unlock()

		broadcastRoomState(c.RoomID, room)

	case "queue_remove":
		var p struct {
			QID string `json:"qid"`
		}
		json.Unmarshal(msg.Payload, &p)
		room.Mu.Lock()
		for i, qi := range room.Queue {
			if qi.QID == p.QID && (qi.AddedBy == c.MemberID || room.HostID == c.MemberID) {
				room.Queue = append(room.Queue[:i], room.Queue[i+1:]...)
				break
			}
		}
		room.Mu.Unlock()
		broadcastRoomState(c.RoomID, room)

	case "vote_up":
		var p struct {
			QID    string `json:"qid"`
			Amount int    `json:"amount"`
		}
		json.Unmarshal(msg.Payload, &p)
		if p.Amount != 1 && p.Amount != 3 && p.Amount != 5 {
			return
		}
		room.Mu.Lock()
		member := room.Members[c.MemberID]
		if member == nil || member.Coins < p.Amount {
			room.Mu.Unlock()
			return
		}
		for _, qi := range room.Queue {
			if qi.QID == p.QID {
				member.Coins -= p.Amount
				qi.VoteMap[c.MemberID] += p.Amount
				qi.Votes += p.Amount
				break
			}
		}
		// re-sort queue by votes (stable, preserve original order for ties)
		sortQueueByVotes(room.Queue)
		room.Mu.Unlock()
		broadcastRoomState(c.RoomID, room)

	case "vote_skip":
		room.Mu.Lock()
		totalMembers := len(room.Members)
		room.SkipVotes[c.MemberID] = true
		skipCount := len(room.SkipVotes)
		room.Mu.Unlock()

		hub.GlobalHub.Broadcast(c.RoomID, hub.WSMessage{
			Type:    "skip_vote_update",
			Payload: mustJSON(map[string]interface{}{"votes": skipCount, "total": totalMembers}),
		})

		if totalMembers > 0 && float64(skipCount)/float64(totalMembers) > 0.5 {
			room.Mu.Lock()
			room.SkipVotes = make(map[string]bool)
			startNext(room)
			room.Mu.Unlock()
			broadcastRoomState(c.RoomID, room)
		}

	case "reaction":
		var p struct {
			Emoji string `json:"emoji"`
		}
		json.Unmarshal(msg.Payload, &p)
		if p.Emoji == "" {
			return
		}
		room.Mu.RLock()
		member := room.Members[c.MemberID]
		room.Mu.RUnlock()
		hub.GlobalHub.Broadcast(c.RoomID, hub.WSMessage{
			Type: "reaction",
			Payload: mustJSON(map[string]interface{}{
				"emoji":    p.Emoji,
				"memberId": c.MemberID,
				"nickname": member.Nickname,
			}),
		})

	case "playback_sync":
		// Only host can push playback state
		room.Mu.RLock()
		isHost := room.HostID == c.MemberID
		room.Mu.RUnlock()
		if !isHost {
			return
		}
		var p struct {
			IsPlaying bool    `json:"isPlaying"`
			Position  float64 `json:"position"`
		}
		json.Unmarshal(msg.Payload, &p)
		room.Mu.Lock()
		if room.Playback != nil {
			room.Playback.IsPlaying = p.IsPlaying
			room.Playback.Position = p.Position
			room.Playback.StartedAt = time.Now().UnixMilli()
		}
		room.Mu.Unlock()
		hub.GlobalHub.BroadcastExcept(c.RoomID, c, hub.WSMessage{
			Type: "playback_sync",
			Payload: mustJSON(map[string]interface{}{
				"isPlaying": p.IsPlaying,
				"position":  p.Position,
				"startedAt": time.Now().UnixMilli(),
			}),
		})

	case "next_song":
		room.Mu.RLock()
		isHost := room.HostID == c.MemberID
		room.Mu.RUnlock()
		if !isHost {
			return
		}
		room.Mu.Lock()
		room.SkipVotes = make(map[string]bool)
		startNext(room)
		room.Mu.Unlock()
		broadcastRoomState(c.RoomID, room)
	}
}

// startNext advances to the next song. Must be called with room.Mu held (write).
func startNext(room *models.Room) {
	if len(room.Queue) == 0 {
		room.Playback = &models.PlaybackState{}
		return
	}
	next := room.Queue[0]
	room.Queue = room.Queue[1:]
	room.Playback = &models.PlaybackState{
		Song:      &next.Song,
		IsPlaying: true,
		Position:  0,
		StartedAt: time.Now().UnixMilli(),
	}
	// Give score to adder
	if m, ok := room.Members[next.AddedBy]; ok {
		m.Score += 1
	}
}

func sortQueueByVotes(queue []*models.QueueItem) {
	// Insertion sort (stable) descending by votes
	for i := 1; i < len(queue); i++ {
		for j := i; j > 0 && queue[j].Votes > queue[j-1].Votes; j-- {
			queue[j], queue[j-1] = queue[j-1], queue[j]
		}
	}
}

func broadcastRoomState(roomID string, room *models.Room) {
	state := hub.RoomStatePayload(room)
	stateJSON, _ := json.Marshal(state)
	hub.GlobalHub.Broadcast(roomID, hub.WSMessage{
		Type:    "room_state",
		Payload: stateJSON,
	})
}

func mustJSON(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func avatarURL(nickname string) string {
	// DiceBear identicon via color hash
	h := 0
	for _, ch := range nickname {
		h = h*31 + int(ch)
	}
	colors := []string{"7c3aed", "2563eb", "059669", "dc2626", "d97706", "db2777"}
	color := colors[abs(h)%len(colors)]
	_ = color
	initials := string([]rune(nickname)[:min(1, len([]rune(nickname)))])
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=fff&size=64&bold=true", initials, color)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

