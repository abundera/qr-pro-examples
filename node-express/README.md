# node-express

End-to-end Node 18+ / Express / TypeScript example for the Abundera QR Pro API.

Demonstrates:

- Creating and updating dynamic codes
- Reading analytics
- Receiving webhooks with HMAC-SHA256 signature verification
- Replay protection (reject signatures seen in the last 10 minutes)
- Raw-body handling on the webhook route (critical — parsing JSON first breaks verification)

## Run

```bash
npm install
ABUNDERA_API_KEY=abnd_qrpro_... \
ABUNDERA_WEBHOOK_SECRET=whsec_... \
  npm start
```

The server listens on `:3000` by default. Point a webhook at `http://your-host:3000/webhook`.

## Test

```bash
npm test
```
