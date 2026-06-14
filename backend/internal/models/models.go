package models

import (
	"sync"
	"time"
)

type Member struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Coins    int    `json:"coins"`
	Score    int    `json:"score"`
	Persona  string `json:"persona"`
	IsHost   bool   `json:"isHost"`
}

type Song struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Cover    string `json:"cover"`
	Duration int    `json:"duration"` // seconds
}

type QueueItem struct {
	QID     string         `json:"qid"`
	Song    Song           `json:"song"`
	AddedBy string         `json:"addedBy"` // member ID
	Adder   string         `json:"adder"`   // member nickname
	Votes   int            `json:"votes"`
	VoteMap map[string]int `json:"-"` // memberID -> coins spent
}

type PlaybackState struct {
	Song      *Song   `json:"song"`
	IsPlaying bool    `json:"isPlaying"`
	Position  float64 `json:"position"` // seconds
	StartedAt int64   `json:"startedAt"` // unix ms, when position was last set
}

type ChatMessage struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
	Time     int64  `json:"time"`
}

type Room struct {
	Mu sync.RWMutex

	ID       string             `json:"id"`
	Code     string             `json:"code"`
	Name     string             `json:"name"`
	HostID   string             `json:"hostId"`
	Members  map[string]*Member `json:"members"`
	Queue    []*QueueItem       `json:"queue"`
	Playback *PlaybackState     `json:"playback"`
	Messages []ChatMessage      `json:"messages"`

	// skip vote tracking
	SkipVotes   map[string]bool `json:"skipVotes"`
	SkipInitBy  string          `json:"skipInitBy"`
	SkipStartAt time.Time       `json:"-"`

	// round-robin: last adder per slot
	RoundRobin []string `json:"roundRobin"` // ordered member IDs for next turn
}

func NewRoom(id, code, name, hostID string) *Room {
	return &Room{
		ID:        id,
		Code:      code,
		Name:      name,
		HostID:    hostID,
		Members:   make(map[string]*Member),
		Queue:     make([]*QueueItem, 0),
		Playback:  &PlaybackState{},
		Messages:  make([]ChatMessage, 0, 100),
		SkipVotes: make(map[string]bool),
	}
}

func (r *Room) MemberList() []*Member {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	list := make([]*Member, 0, len(r.Members))
	for _, m := range r.Members {
		list = append(list, m)
	}
	return list
}

func (r *Room) MemberCount() int {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return len(r.Members)
}

// InsertQueued inserts a new queue item using round-robin ordering.
// Items from the same user are spread evenly among items from others.
func (r *Room) InsertQueued(item *QueueItem) {
	// Simple round-robin: find the last item by this user, insert after the group
	lastSameUserIdx := -1
	for i := len(r.Queue) - 1; i >= 0; i-- {
		if r.Queue[i].AddedBy == item.AddedBy {
			lastSameUserIdx = i
			break
		}
	}
	if lastSameUserIdx == -1 {
		// User hasn't queued anything, append
		r.Queue = append(r.Queue, item)
		return
	}
	// Insert right after the last item from this user's "round"
	// Find the position after lastSameUserIdx where a different user's block ends
	insertAt := lastSameUserIdx + 1
	// Skip over a single block of other users (one round)
	otherCount := 0
	for insertAt < len(r.Queue) {
		if r.Queue[insertAt].AddedBy != item.AddedBy {
			otherCount++
			if otherCount >= len(r.Members)-1 && len(r.Members) > 1 {
				insertAt++
				break
			}
			insertAt++
		} else {
			break
		}
	}
	// Insert at insertAt
	r.Queue = append(r.Queue, nil)
	copy(r.Queue[insertAt+1:], r.Queue[insertAt:])
	r.Queue[insertAt] = item
}
