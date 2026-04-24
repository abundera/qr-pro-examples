# fastapi

End-to-end FastAPI / Python 3.11+ example for the Abundera QR Pro API.

Demonstrates:

- Creating and updating dynamic codes
- Reading analytics
- Receiving webhooks with HMAC-SHA256 signature verification
- Replay protection (reject signatures seen in the last 10 minutes)
- Raw-body read before JSON parse (`request.body()` before `request.json()` — verification needs the exact bytes)

## Run

```bash
pip install -e ".[dev]"
ABUNDERA_API_KEY=abnd_qrpro_... \
ABUNDERA_WEBHOOK_SECRET=whsec_... \
  uvicorn server:app --reload
```

## Test

```bash
pytest
```
