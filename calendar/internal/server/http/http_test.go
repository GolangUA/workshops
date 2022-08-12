package http

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http/httptest"
)

var contextUser = "user"

func authenticatedContext() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(auth.ContextKey, &auth.Context{
		JWT: &auth.Claims{
			Timezone: "Europe/Kiev",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: contextUser,
			},
		},
	})
	return w, c
}
