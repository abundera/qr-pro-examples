// Smoke test: verify the webhook verifier integrates cleanly.
import { verifyWebhookSignature } from "@abundera/qr-pro";
import { createHmac } from "node:crypto";
import { test } from "node:test";
import { strict as assert } from "node:assert";

test("valid signature verifies", () => {
  const secret = "whsec_test";
  const body = '{"event":"code.scanned"}';
  const ts = String(Math.floor(Date.now() / 1000));
  const v1 = createHmac("sha256", secret).update(`${ts}.${body}`).digest("hex");
  assert.doesNotThrow(() =>
    verifyWebhookSignature({ signature: `t=${ts},v1=${v1}`, body, secret }),
  );
});

test("tampered body rejects", () => {
  const secret = "whsec_test";
  const body = '{"event":"code.scanned"}';
  const ts = String(Math.floor(Date.now() / 1000));
  const v1 = createHmac("sha256", secret).update(`${ts}.${body}`).digest("hex");
  assert.throws(() =>
    verifyWebhookSignature({
      signature: `t=${ts},v1=${v1}`,
      body: body + "x",
      secret,
    }),
  );
});
