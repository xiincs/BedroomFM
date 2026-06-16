package api

import (
	"bedroomfm/internal/auth"
	"bedroomfm/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/auth/register
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" || req.Nickname == "" {
		c.JSON(400, gin.H{"error": "username、password、nickname 为必填项"})
		return
	}
	if len(req.Username) < 3 || len(req.Username) > 20 {
		c.JSON(400, gin.H{"error": "用户名须在 3–20 个字符之间"})
		return
	}
	if len(req.Password) < 6 {
		c.JSON(400, gin.H{"error": "密码至少 6 位"})
		return
	}
	if models.GlobalUserStore.GetByUsername(req.Username) != nil {
		c.JSON(409, gin.H{"error": "用户名已被占用"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "服务器错误"})
		return
	}
	userID := uuid.NewString()
	user := &models.User{
		ID:           userID,
		Username:     req.Username,
		Nickname:     req.Nickname,
		Avatar:       avatarURL(req.Nickname),
		PasswordHash: string(hash),
		XP:           models.XPRegister,
		Level:        models.LevelFromXP(models.XPRegister),
		CreatedAt:    time.Now().Unix(),
	}
	models.GlobalUserStore.Create(user)
	token, _ := auth.Sign(userID)
	c.JSON(200, buildAuthResponse(token, user))
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "username 和 password 为必填项"})
		return
	}
	user := models.GlobalUserStore.GetByUsername(req.Username)
	if user == nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}
	token, _ := auth.Sign(user.ID)
	c.JSON(200, buildAuthResponse(token, user))
}

// GET /api/auth/me
func GetMe(c *gin.Context) {
	user := getAuthUser(c)
	if user == nil {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}
	tier := models.TierForLevel(user.Level)
	c.JSON(200, gin.H{
		"user":       toPublicUser(user),
		"tier":       tier,
		"xpToNext":   models.XPToNextLevel(user.Level),
		"xpProgress": user.XP - models.XPForLevel(user.Level),
		"allTiers":   models.VIPTiers,
	})
}

// POST /api/user/xp
func AddXP(c *gin.Context) {
	user := getAuthUser(c)
	if user == nil {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}
	var req struct {
		Action string `json:"action"`
	}
	c.ShouldBindJSON(&req)
	valid := map[string]bool{
		"join_room": true, "add_song": true, "song_played": true, "vote_up": true, "chat": true,
	}
	if !valid[req.Action] {
		c.JSON(400, gin.H{"error": "无效的 action"})
		return
	}
	updated, gained := models.GlobalUserStore.GrantXP(user.ID, req.Action)
	if updated == nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}
	tier := models.TierForLevel(updated.Level)
	c.JSON(200, gin.H{
		"xp":         updated.XP,
		"level":      updated.Level,
		"gained":     gained,
		"tier":       tier,
		"xpToNext":   models.XPToNextLevel(updated.Level),
		"xpProgress": updated.XP - models.XPForLevel(updated.Level),
	})
}

func getAuthUser(c *gin.Context) *models.User {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil
	}
	claims, err := auth.Verify(strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		return nil
	}
	return models.GlobalUserStore.GetByID(claims.UserID)
}

func buildAuthResponse(token string, user *models.User) gin.H {
	tier := models.TierForLevel(user.Level)
	return gin.H{
		"token":      token,
		"user":       toPublicUser(user),
		"tier":       tier,
		"xpToNext":   models.XPToNextLevel(user.Level),
		"xpProgress": user.XP - models.XPForLevel(user.Level),
		"allTiers":   models.VIPTiers,
	}
}

type publicUser struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	XP         int    `json:"xp"`
	Level      int    `json:"level"`
	TotalSongs int    `json:"totalSongs"`
	TotalVotes int    `json:"totalVotes"`
	TotalRooms int    `json:"totalRooms"`
	CreatedAt  int64  `json:"createdAt"`
}

func toPublicUser(u *models.User) publicUser {
	return publicUser{
		u.ID, u.Username, u.Nickname, u.Avatar,
		u.XP, u.Level, u.TotalSongs, u.TotalVotes, u.TotalRooms, u.CreatedAt,
	}
}
