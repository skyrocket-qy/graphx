package middleware_test

import (
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
