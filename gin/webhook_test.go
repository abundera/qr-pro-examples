package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	qrpro "github.com/abundera/qr-pro-go"
)

func TestWebhookSignatureRoundTrip(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"event":"code.scanned"}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("."))
	mac.Write(body)
	sig := "t=" + ts + ",v1=" + hex.EncodeToString(mac.Sum(nil))

	if err := qrpro.VerifyWebhookSignature(sig, body, secret, 0); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := qrpro.VerifyWebhookSignature(sig, append(body, 'x'), secret, 0); err == nil {
		t.Fatal("expected rejection for tampered body")
	}
}
