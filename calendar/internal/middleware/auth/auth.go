package auth

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/api"
	"github.com/Roma7-7-7/workshops/calendar/internal/logging"
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/metrics"
	"github.com/Roma7-7-7/workshops/calendar/internal/services/calendar"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const jwtAliveDuration = 5 * time.Minute

type Middleware struct {
	repo   calendar.Repository
	secret string
}

type Claims struct {
	Timezone string
	jwt.RegisteredClaims
}

const ContextKey = "auth"

type Context struct {
	JWT *Claims
}

func (c *Context) Username() string {
	return c.JWT.Subject
}

func (c *Context) UserTimezone() string {
	return c.JWT.Timezone
}

func GetContext(c *gin.Context) *Context {
	return c.MustGet(ContextKey).(*Context)
}

func (m *Middleware) Login(c *gin.Context) {
	var req api.UserPassword
	if err := c.BindJSON(&req); err != nil {
		api.BadJSONA(c)
		return
	}
	u, err := m.repo.GetUser(req.Username)
	if err != nil {
		logging.Logger.Error("get user", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}
	if u == nil {
		api.NotFoundA(c, fmt.Sprintf("user \"%s\"", u.Name))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err == bcrypt.ErrMismatchedHashAndPassword {
		api.UnauthorizedA(c, "password does not match")
		return
	} else if err != nil {
		logging.Logger.Error("validate password", zap.Error(err))
		api.ServerErrorA(c, err)
		return
	}

	now := time.Now()
	expires := now.Add(jwtAliveDuration)
	claims := &jwt.RegisteredClaims{
		Issuer:    "calendar-app",
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(expires),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(m.secret))
	if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	c.SetCookie("token", tokenS, int(jwtAliveDuration.Seconds()), "/", "calendar-app", false, false)
}

func (m *Middleware) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "calendar-app", false, false)
}

func (m *Middleware) Validate(c *gin.Context) {
	tokenS, err := c.Cookie("token")
	if err == http.ErrNoCookie || (err == nil && tokenS == "") {
		api.UnauthorizedA(c, "request is not authenticated")
		return
	} else if err != nil {
		api.ServerErrorA(c, err)
		return
	}

	cl := &Claims{}
	if _, err = jwt.ParseWithClaims(tokenS, cl, m.keyFunc); err != nil {
		api.ServerErrorA(c, err)
		return
	}
	c.Set(ContextKey, &Context{
		JWT: cl,
	})
	metrics.IncRequest(cl.Subject)
}

func (m *Middleware) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(m.secret), nil
}

func NewMiddleware(repo calendar.Repository, secret string) *Middleware {
	return &Middleware{
		repo:   repo,
		secret: secret,
	}
}
