package hub

import (
	"bedroomfm/internal/models"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	hub      *RoomHub
	conn     *websocket.Conn
	send     chan []byte
	RoomID   string
	MemberID string
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type RoomHub struct {
	mu      sync.RWMutex
	rooms   map[string]*models.Room
	clients map[string]map[*Client]bool // roomID -> set of clients
}

var GlobalHub = &RoomHub{
	rooms:   make(map[string]*models.Room),
	clients: make(map[string]map[*Client]bool),
}

func (h *RoomHub) GetRoom(id string) *models.Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[id]
}

func (h *RoomHub) GetRoomByCode(code string) *models.Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, r := range h.rooms {
		if r.Code == code {
			return r
		}
	}
	return nil
}

func (h *RoomHub) CreateRoom(room *models.Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.rooms[room.ID] = room
	h.clients[room.ID] = make(map[*Client]bool)
}

func (h *RoomHub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.RoomID] == nil {
		h.clients[c.RoomID] = make(map[*Client]bool)
	}
	h.clients[c.RoomID][c] = true
}

func (h *RoomHub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.RoomID]; ok {
		delete(set, c)
	}
}

func (h *RoomHub) Broadcast(roomID string, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("broadcast marshal error:", err)
		return
	}
	h.mu.RLock()
	clients := h.clients[roomID]
	h.mu.RUnlock()

	for c := range clients {
		select {
		case c.send <- data:
		default:
			// slow client, drop
		}
	}
}

func (h *RoomHub) BroadcastExcept(roomID string, except *Client, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	clients := h.clients[roomID]
	h.mu.RUnlock()

	for c := range clients {
		if c == except {
			continue
		}
		select {
		case c.send <- data:
		default:
		}
	}
}

func (c *Client) SendRaw(data []byte) {
	select {
	case c.send <- data:
	default:
	}
}

func NewClient(h *RoomHub, conn *websocket.Conn, roomID, memberID string) *Client {
	return &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		RoomID:   roomID,
		MemberID: memberID,
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(onMessage func(c *Client, msg WSMessage)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg WSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		onMessage(c, msg)
	}
}

func mustJSON(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func roomState(r *models.Room) map[string]interface{} {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	members := make([]*models.Member, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(members, m)
	}
	return map[string]interface{}{
		"id":        r.ID,
		"code":      r.Code,
		"name":      r.Name,
		"hostId":    r.HostID,
		"members":   members,
		"queue":     r.Queue,
		"playback":  r.Playback,
		"skipVotes": r.SkipVotes,
		"messages":  r.Messages,
	}
}

func RoomStatePayload(r *models.Room) map[string]interface{} {
	return roomState(r)
}
