package models

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	PasswordHash string `json:"-"`
	XP           int    `json:"xp"`
	Level        int    `json:"level"`
	TotalSongs   int    `json:"totalSongs"`
	TotalVotes   int    `json:"totalVotes"`
	TotalRooms   int    `json:"totalRooms"`
	CreatedAt    int64  `json:"createdAt"`
}

// XP curve: reaching level n costs XPForLevel(n) total. Going n→n+1 costs n*100.
func LevelFromXP(xp int) int {
	n := 1
	for n < 100 && XPForLevel(n+1) <= xp {
		n++
	}
	return n
}

func XPForLevel(n int) int {
	if n <= 1 {
		return 0
	}
	return n * (n - 1) / 2 * 100
}

func XPToNextLevel(n int) int {
	if n >= 100 {
		return 0
	}
	return n * 100
}

type VIPTier struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	MinLevel int      `json:"minLevel"`
	MaxLevel int      `json:"maxLevel"`
	Color    string   `json:"color"`
	Gradient []string `json:"gradient"`
	Perks    []string `json:"perks"`
}

var VIPTiers = []VIPTier{
	{0, "听众", 1, 4, "#6b7280", []string{"#6b7280", "#9ca3af"},
		[]string{"参与实时音乐房间", "基础点歌与顶歌", "表情实时互动"}},
	{1, "小歌迷", 5, 9, "#3b82f6", []string{"#3b82f6", "#60a5fa"},
		[]string{"蓝色专属昵称", "解锁更多表情包", "优先搜索通道"}},
	{2, "节奏达人", 10, 19, "#06b6d4", []string{"#06b6d4", "#10b981"},
		[]string{"青绿渐变昵称", "炫彩弹幕效果", "专属节奏图标"}},
	{3, "黄金DJ", 20, 34, "#f59e0b", []string{"#f59e0b", "#ef4444"},
		[]string{"金色昵称光效", "顶歌额外加成 +1", "VIP 专属房间权限"}},
	{4, "白金制作人", 35, 49, "#8b5cf6", []string{"#8b5cf6", "#ec4899"},
		[]string{"紫粉渐变昵称", "高级粒子表情效果", "月度 MVP 专属徽章"}},
	{5, "钻石传说", 50, 74, "#06b6d4", []string{"#06b6d4", "#8b5cf6", "#ec4899"},
		[]string{"动态钻石头像框", "钻石专属徽章动画", "双倍顶歌效果"}},
	{6, "传奇之声", 75, 99, "#a855f7", []string{"#f59e0b", "#ef4444", "#8b5cf6", "#06b6d4"},
		[]string{"彩虹流光昵称特效", "全站传奇专属标识", "专属系统公告气泡"}},
	{7, "音乐之神", 100, 100, "#fbbf24", []string{"#fbbf24", "#ef4444", "#8b5cf6", "#06b6d4", "#10b981"},
		[]string{"神话级粒子光效", "永久 DJ 专属头衔", "全部权限永久解锁"}},
}

func TierForLevel(level int) VIPTier {
	for i := len(VIPTiers) - 1; i >= 0; i-- {
		if level >= VIPTiers[i].MinLevel {
			return VIPTiers[i]
		}
	}
	return VIPTiers[0]
}

const (
	XPRegister   = 50
	XPJoinRoom   = 10
	XPAddSong    = 20
	XPSongPlayed = 50
	XPVoteUp     = 5
	XPChat       = 2
)

// persistedUser mirrors User but includes PasswordHash for file storage.
// The main User struct keeps json:"-" on PasswordHash to prevent API leakage.
type persistedUser struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	PasswordHash string `json:"passwordHash"`
	XP           int    `json:"xp"`
	Level        int    `json:"level"`
	TotalSongs   int    `json:"totalSongs"`
	TotalVotes   int    `json:"totalVotes"`
	TotalRooms   int    `json:"totalRooms"`
	CreatedAt    int64  `json:"createdAt"`
}

type UserStore struct {
	mu         sync.RWMutex
	users      map[string]*User
	byUsername map[string]*User
	filePath   string
}

// NewUserStore creates a store backed by filePath.
// Existing data is loaded immediately; pass "" to disable persistence.
func NewUserStore(filePath string) *UserStore {
	s := &UserStore{
		users:      make(map[string]*User),
		byUsername: make(map[string]*User),
		filePath:   filePath,
	}
	if filePath != "" {
		s.load()
	}
	return s
}

func (s *UserStore) load() {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Printf("UserStore: failed to read %s: %v", s.filePath, err)
		return
	}
	var list []persistedUser
	if err := json.Unmarshal(data, &list); err != nil {
		log.Printf("UserStore: failed to parse %s: %v", s.filePath, err)
		return
	}
	for i := range list {
		p := &list[i]
		u := &User{
			ID: p.ID, Username: p.Username, Nickname: p.Nickname,
			Avatar: p.Avatar, PasswordHash: p.PasswordHash,
			XP: p.XP, Level: p.Level,
			TotalSongs: p.TotalSongs, TotalVotes: p.TotalVotes, TotalRooms: p.TotalRooms,
			CreatedAt: p.CreatedAt,
		}
		s.users[u.ID] = u
		s.byUsername[u.Username] = u
	}
	log.Printf("UserStore: loaded %d users from %s", len(list), s.filePath)
}

// save writes the store to disk; must be called with mu write-locked.
func (s *UserStore) save() {
	if s.filePath == "" {
		return
	}
	list := make([]persistedUser, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, persistedUser{
			ID: u.ID, Username: u.Username, Nickname: u.Nickname,
			Avatar: u.Avatar, PasswordHash: u.PasswordHash,
			XP: u.XP, Level: u.Level,
			TotalSongs: u.TotalSongs, TotalVotes: u.TotalVotes, TotalRooms: u.TotalRooms,
			CreatedAt: u.CreatedAt,
		})
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Printf("UserStore: marshal error: %v", err)
		return
	}
	// Write atomically: temp file → rename
	tmp := s.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		log.Printf("UserStore: write error: %v", err)
		return
	}
	if err := os.Rename(tmp, s.filePath); err != nil {
		log.Printf("UserStore: rename error: %v", err)
	}
}

func (s *UserStore) Create(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	s.byUsername[user.Username] = user
	s.save()
}

func (s *UserStore) GetByID(id string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *UserStore) GetByUsername(username string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.byUsername[username]
}

func (s *UserStore) GrantXP(id string, action string) (*User, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := s.users[id]
	if u == nil {
		return nil, 0
	}
	amount := 0
	switch action {
	case "join_room":
		amount = XPJoinRoom
		u.TotalRooms++
	case "add_song":
		amount = XPAddSong
		u.TotalSongs++
	case "song_played":
		amount = XPSongPlayed
	case "vote_up":
		amount = XPVoteUp
		u.TotalVotes++
	case "chat":
		amount = XPChat
	}
	if amount > 0 {
		u.XP += amount
		u.Level = LevelFromXP(u.XP)
		s.save()
	}
	return u, amount
}

// GlobalUserStore is initialized by main via InitUserStore.
var GlobalUserStore *UserStore

func InitUserStore(filePath string) {
	GlobalUserStore = NewUserStore(filePath)
}
