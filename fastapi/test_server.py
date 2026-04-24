import hashlib
import hmac
import time

from abundera_qr_pro import verify_webhook_signature
from abundera_qr_pro.webhook import WebhookVerificationError


def _signed(secret: str, body: str) -> str:
    ts = str(int(time.time()))
    v1 = hmac.new(secret.encode(), f"{ts}.{body}".encode(), hashlib.sha256).hexdigest()
    return f"t={ts},v1={v1}"


def test_valid_signature():
    sig = _signed("whsec_test", '{"event":"x"}')
    assert verify_webhook_signature(
        signature=sig, body='{"event":"x"}', secret="whsec_test"
    )


def test_tampered_body():
    sig = _signed("whsec_test", '{"event":"x"}')
    try:
        verify_webhook_signature(
            signature=sig, body='{"event":"x"}_tamper', secret="whsec_test"
        )
    except WebhookVerificationError:
        return
    raise AssertionError("expected rejection")
