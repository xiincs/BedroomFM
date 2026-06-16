package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const jwtSecret = "bedroomfm-s3cr3t-2024-xk9"

type Claims struct {
	UserID string `json:"uid"`
	Exp    int64  `json:"exp"`
}

func b64(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Sign(userID string) (string, error) {
	header := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims := Claims{UserID: userID, Exp: time.Now().Add(30 * 24 * time.Hour).Unix()}
	payload, _ := json.Marshal(claims)
	data := header + "." + b64(payload)
	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(data))
	return data + "." + b64(mac.Sum(nil)), nil
}

func Verify(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	data := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(data))
	if !hmac.Equal([]byte(b64(mac.Sum(nil))), []byte(parts[2])) {
		return nil, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}
