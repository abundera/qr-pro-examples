"""End-to-end FastAPI example for the Abundera QR Pro API.

Run: ABUNDERA_API_KEY=... ABUNDERA_WEBHOOK_SECRET=... uvicorn server:app
"""

from __future__ import annotations

import os
import time
from collections import OrderedDict
from typing import Any

from abundera_qr_pro import Client
from abundera_qr_pro.webhook import (
    WebhookVerificationError,
    verify_webhook_signature,
)
from fastapi import FastAPI, Header, HTTPException, Request

API_KEY = os.environ["ABUNDERA_API_KEY"]
WEBHOOK_SECRET = os.environ["ABUNDERA_WEBHOOK_SECRET"]

app = FastAPI(title="qr-pro-example-fastapi")
client = Client(api_key=API_KEY)


# Replay protection — drop signatures we've already processed in the last 10 minutes.
_SEEN: OrderedDict[str, float] = OrderedDict()
_SEEN_TTL = 600.0


def _remember(signature: str) -> bool:
    now = time.time()
    # Evict expired entries.
    while _SEEN and next(iter(_SEEN.values())) < now - _SEEN_TTL:
        _SEEN.popitem(last=False)
    if signature in _SEEN:
        return False
    _SEEN[signature] = now
    return True


@app.post("/codes")
def create_code(payload: dict[str, Any]) -> dict:
    c = client.create_code(**payload)
    return {"id": c.id, "short_url": c.short_url, "destination_url": c.destination_url}


@app.patch("/codes/{code_id}")
def update_code(code_id: str, payload: dict[str, Any]) -> dict:
    c = client.update_code(code_id, **payload)
    return {"id": c.id, "destination_url": c.destination_url, "active": c.active}


@app.get("/codes/{code_id}/analytics")
def analytics(code_id: str) -> dict:
    s = client.get_analytics(code_id)
    return {"total_scans": s.total_scans, "by_country": s.by_country}


@app.post("/webhook")
async def webhook(
    request: Request,
    x_abundera_signature: str = Header(default=""),
) -> dict:
    if not x_abundera_signature:
        raise HTTPException(400, "missing signature")

    body = (await request.body()).decode("utf-8")
    try:
        verify_webhook_signature(
            signature=x_abundera_signature,
            body=body,
            secret=WEBHOOK_SECRET,
        )
    except WebhookVerificationError as e:
        raise HTTPException(400, f"bad signature: {e}") from e

    if not _remember(x_abundera_signature):
        return {"status": "duplicate"}

    # Parse after verification.
    import json

    payload = json.loads(body)
    print(f"webhook: {payload.get('event')} {payload}")
    return {"status": "ok"}
