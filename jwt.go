package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
)

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

const IdentityKey = "id"

func CreateJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm: "deltaFi",
		// Secret key should be passed by environment variable. But we use fixed value here for demo.
		Key:         []byte("deltaFi demo project"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: LoginHandler,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false, // We don't use https for local demo
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080",         // Use localhost for local demo
		CookieSameSite: http.SameSiteDefaultMode, // Only same site for local demo
	})
}

func LoginHandler(c *gin.Context) (interface{}, error) {
	var login Login
	if err := c.ShouldBind(&login); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := login.Email
	password := login.Password

	if userID == "cuizk" && password == "123" {
		return &User{
			Email:     userID,
			FirstName: "Zikun",
			LastName:  "Cui",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}
