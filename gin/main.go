// End-to-end Gin / Go 1.21+ example for the Abundera QR Pro API.
//
// Run: ABUNDERA_API_KEY=... ABUNDERA_WEBHOOK_SECRET=... go run .
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	qrpro "github.com/abundera/qr-pro-go"
	"github.com/gin-gonic/gin"
)

const replayTTL = 10 * time.Minute

type server struct {
	client        *qrpro.Client
	webhookSecret string

	mu   sync.Mutex
	seen map[string]time.Time
}

func main() {
	apiKey := mustenv("ABUNDERA_API_KEY")
	webhookSecret := mustenv("ABUNDERA_WEBHOOK_SECRET")

	client, err := qrpro.New(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		client:        client,
		webhookSecret: webhookSecret,
		seen:          map[string]time.Time{},
	}

	r := gin.Default()
	r.POST("/codes", s.createCode)
	r.PATCH("/codes/:id", s.updateCode)
	r.GET("/codes/:id/analytics", s.analytics)
	r.POST("/webhook", s.webhook)

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func (s *server) createCode(c *gin.Context) {
	var in qrpro.CodeCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := s.client.CreateCode(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, code)
}

func (s *server) updateCode(c *gin.Context) {
	var patch qrpro.CodePatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := s.client.UpdateCode(c.Request.Context(), c.Param("id"), patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, code)
}

func (s *server) analytics(c *gin.Context) {
	stats, err := s.client.GetAnalytics(c.Request.Context(), c.Param("id"), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (s *server) webhook(c *gin.Context) {
	sig := c.GetHeader("X-Abundera-Signature")
	if sig == "" {
		c.String(http.StatusBadRequest, "missing signature")
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "read body")
		return
	}
	if err := qrpro.VerifyWebhookSignature(sig, body, s.webhookSecret, 0); err != nil {
		c.String(http.StatusBadRequest, "bad signature")
		return
	}
	if !s.seenOnce(sig) {
		c.String(http.StatusOK, "duplicate")
		return
	}

	var payload map[string]any
	_ = json.Unmarshal(body, &payload)
	log.Printf("webhook: %v", payload["event"])
	c.String(http.StatusOK, "ok")
}

func (s *server) seenOnce(sig string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, t := range s.seen {
		if now.Sub(t) > replayTTL {
			delete(s.seen, k)
		}
	}
	if _, ok := s.seen[sig]; ok {
		return false
	}
	s.seen[sig] = now
	return true
}

func mustenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s is required", k)
	}
	return v
}

var _ = context.Background // touch context import to avoid unused-import on future refactors
