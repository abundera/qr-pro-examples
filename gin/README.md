# gin

End-to-end Gin / Go 1.21+ example for the Abundera QR Pro API.

Demonstrates:

- Creating and updating dynamic codes
- Reading analytics
- Receiving webhooks with HMAC-SHA256 signature verification
- Replay protection (reject signatures seen in the last 10 minutes)
- Raw body read before JSON unmarshal (verification needs the exact bytes)

## Run

```bash
go mod tidy
ABUNDERA_API_KEY=abnd_qrpro_... \
ABUNDERA_WEBHOOK_SECRET=whsec_... \
  go run .
```

## Test

```bash
go test ./...
```
