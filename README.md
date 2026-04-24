# Abundera QR Pro — example apps

Three minimal but real example apps demonstrating the [Abundera QR Pro API](https://pro.qr.abundera.ai/docs/) and [official SDKs](https://pro.qr.abundera.ai/sdks/).

| Example | Stack | SDK |
|---|---|---|
| [`node-express`](./node-express/) | Node 18+ · Express · TypeScript | [`@abundera/qr-pro`](https://www.npmjs.com/package/@abundera/qr-pro) |
| [`fastapi`](./fastapi/) | Python 3.11+ · FastAPI · uvicorn | [`abundera-qr-pro`](https://pypi.org/project/abundera-qr-pro/) |
| [`gin`](./gin/) | Go 1.21+ · Gin | [`github.com/abundera/qr-pro-go`](https://pkg.go.dev/github.com/abundera/qr-pro-go) |

Each one demonstrates:

- Creating and updating dynamic QR codes
- Reading scan analytics
- Receiving webhooks with HMAC-SHA256 signature verification
- Replay protection (reject signatures seen in the last 10 minutes)
- The raw-body-before-parse rule (the signature is over the exact bytes, so you must not let a body parser touch them first)

All three are CI-tested against the SDK's own webhook verifier so the signature shape stays honest.

## Running any example

Every example reads two environment variables:

```
ABUNDERA_API_KEY        = your Pro API key (abnd_qrpro_...)
ABUNDERA_WEBHOOK_SECRET = the secret returned by POST /api/webhooks
```

Get both at https://pro.qr.abundera.ai/account/ (or via the API).

## License

MIT © Abundera, Inc.
