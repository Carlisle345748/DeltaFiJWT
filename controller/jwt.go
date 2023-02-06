package controller

import (
	"net/http"
	"time"

	"deltaFiJWT/dao"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
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
			if v, ok := data.(*dao.User); ok {
				return jwt.MapClaims{IdentityKey: v.ID}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, token string, time time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":        CodeSuccess,
				"message":     "success",
				"token":       token,
				"expire_time": time,
			})
		},
		Authenticator: LoginHandler,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    CodeLoginFail,
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

	user, err := dao.Authenticate(login.Email, login.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
