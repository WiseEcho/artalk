package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	UID    int    `json:"-"`
	jwt.StandardClaims
}

func LoginGetUserToken(user entity.User, key string, ttl int) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID: ID2UID(user.ID),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),                                       // 签发时间
			ExpiresAt: time.Now().Add(time.Second * time.Duration(ttl)).Unix(), // 过期时间
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return t, nil
}

var ErrTokenNotProvided = fmt.Errorf("token not provided")
var ErrTokenUserNotFound = fmt.Errorf("user not found")
var ErrTokenInvalidFromDate = fmt.Errorf("token is invalid starting from a certain date")

func GetTokenByReq(c *fiber.Ctx) string {
	token := c.Query("token")
	if token == "" {
		token = c.FormValue("token")
	}
	if token == "" {
		token = c.Get(fiber.HeaderAuthorization)
		token = strings.TrimPrefix(token, "Bearer ")
	}
	return token
}

func GetJwtDataByReq(app *core.App, c *fiber.Ctx) (jwtCustomClaims, error) {
	token := GetTokenByReq(c)
	if token == "" {
		return jwtCustomClaims{}, ErrTokenNotProvided
	}

	jwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(app.Conf().AppKey), nil // 密钥
	})
	if err != nil {
		return jwtCustomClaims{}, err
	}

	claims := jwtCustomClaims{}
	tmp, _ := json.Marshal(jwt.Claims)
	_ = json.Unmarshal(tmp, &claims)
	claims.UID, err = UID2ID(claims.UserID)
	if err != nil {
		return jwtCustomClaims{}, fmt.Errorf("UID2ID userID:%s,err:%v", claims.UserID, err)
	}
	return claims, nil
}

func GetUserByReq(app *core.App, c *fiber.Ctx) (entity.User, error) {
	claims, err := GetJwtDataByReq(app, c)
	if err != nil {
		return entity.User{}, err
	}

	user := app.Dao().FindUserByID(uint(claims.UID))
	if user.IsEmpty() {
		return entity.User{Model: gorm.Model{ID: uint(claims.UID)}}, ErrTokenUserNotFound
	}

	// check tokenValidFrom
	if user.TokenValidFrom.Valid && claims.IssuedAt < user.TokenValidFrom.Time.Unix() {
		return entity.User{}, ErrTokenInvalidFromDate
	}

	return user, nil
}
