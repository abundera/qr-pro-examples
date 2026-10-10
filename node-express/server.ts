// End-to-end example: create/update/analytics endpoints + webhook receiver
// with signature verification and replay protection.
//
// Run: ABUNDERA_API_KEY=... ABUNDERA_WEBHOOK_SECRET=... npm start
import express, { Request, Response } from "express";
import rateLimit from "express-rate-limit";
import { AbunderaQRProClient, verifyWebhookSignature } from "@abundera/qr-pro";

const port = Number(process.env.PORT ?? 3000);
const apiKey = must("ABUNDERA_API_KEY");
const webhookSecret = must("ABUNDERA_WEBHOOK_SECRET");

const client = new AbunderaQRProClient({ apiKey });
const app = express();

// Cap requests per client before any route runs, so unauthenticated callers
// cannot hammer the signature check or the API proxy routes.
const RATE_WINDOW_MS = 60 * 1000;
const RATE_MAX_REQUESTS = 120;
app.use(rateLimit({ windowMs: RATE_WINDOW_MS, limit: RATE_MAX_REQUESTS }));

// Parse JSON for normal routes, keep raw for the webhook.
app.use("/webhook", express.raw({ type: "*/*" }));
app.use(express.json());

// Replay protection: remember recently-seen signatures for 10 minutes.
const SEEN_TTL_MS = 10 * 60 * 1000;
const seen = new Map<string, number>();
setInterval(() => {
  const cutoff = Date.now() - SEEN_TTL_MS;
  for (const [k, v] of seen) if (v < cutoff) seen.delete(k);
}, 60_000).unref();

app.post("/codes", async (req: Request, res: Response) => {
  try {
    const code = await client.createCode(req.body);
    res.status(201).json(code);
  } catch (e) {
    res.status(500).json({ error: String(e) });
  }
});

app.patch("/codes/:id", async (req: Request, res: Response) => {
  const code = await client.updateCode(req.params.id, req.body);
  res.json(code);
});

app.get("/codes/:id/analytics", async (req: Request, res: Response) => {
  const stats = await client.getAnalytics(req.params.id);
  res.json(stats);
});

app.post("/webhook", (req: Request, res: Response) => {
  const signature = req.header("x-abundera-signature");
  if (!signature) return res.status(400).send("missing signature");

  try {
    verifyWebhookSignature({
      signature,
      body: (req.body as Buffer).toString("utf-8"),
      secret: webhookSecret,
    });
  } catch {
    return res.status(400).send("bad signature");
  }

  // Replay check — reject signatures we've already processed.
  if (seen.has(signature)) return res.status(200).send("duplicate");
  seen.set(signature, Date.now());

  const payload = JSON.parse((req.body as Buffer).toString("utf-8"));
  console.log("webhook:", payload.event, payload);
  res.status(200).send("ok");
});

app.listen(port, () => console.log(`listening on :${port}`));

function must(name: string): string {
  const v = process.env[name];
  if (!v) throw new Error(`${name} is required`);
  return v;
}
