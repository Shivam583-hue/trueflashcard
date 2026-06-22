# Deploying the backend to a VPS via Cloudflare Tunnel

The API runs as a systemd service on the VPS, listening on `localhost:8080`.
Cloudflare Tunnel exposes it at `https://api.nemportfolio.in` (Cloudflare
terminates TLS — no nginx/certbot, no open inbound ports).

```
browser ──HTTPS──▶ Cloudflare ──tunnel──▶ cloudflared ──HTTP──▶ localhost:8080 (server)
```

## 1. Build the binary (on your dev machine)

```
cd server
make build-linux            # set GOARCH=arm64 first if `uname -m` on the VPS is aarch64
```

Produces `server/bin/server-linux`.

## 2. Copy files to the VPS

```
ssh <vps> 'sudo mkdir -p /opt/trueflashcard'
scp server/bin/server-linux        <vps>:/tmp/server
scp deploy/trueflashcard.service   <vps>:/tmp/
ssh <vps> 'sudo mv /tmp/server /opt/trueflashcard/server && sudo chmod +x /opt/trueflashcard/server'
```

## 3. Create the env file on the VPS

Copy `deploy/trueflashcard.env.example` to `/etc/trueflashcard.env`, fill in real
values. If prod uses the **same Neon database** as dev, the schema is already
migrated — otherwise run `make migrate` against the new `DATABASE_URL` first.

```
sudo chmod 600 /etc/trueflashcard.env
```

## 4. Install and start the service

```
sudo mv /tmp/trueflashcard.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now trueflashcard
sudo systemctl status trueflashcard          # should be active (running)
curl -s localhost:8080/flashcard.v1.FolderService/ListFolders \
  -H 'content-type: application/json' -H 'connect-protocol-version: 1' -d '{}'
# expect: {"code":"unauthenticated", ...}  → the server is up
```

## 5. Cloudflare Tunnel (once nemportfolio.in is active on Cloudflare)

```
# install cloudflared on the VPS, then:
cloudflared tunnel login                      # opens a browser link; authorize nemportfolio.in
cloudflared tunnel create flashcards          # prints a TUNNEL_ID + writes a credentials json
cloudflared tunnel route dns flashcards api.nemportfolio.in
```

Put `deploy/cloudflared-config.yml` at `/etc/cloudflared/config.yml` and replace
`<TUNNEL_ID>`. Then run it as a service:

```
sudo cloudflared service install
sudo systemctl enable --now cloudflared
```

Verify: `curl https://api.nemportfolio.in/flashcard.v1.FolderService/ListFolders ...`
should return the same unauthenticated JSON.

## 6. Wire up the two external services

- **Google console** → add Authorized redirect URI: `https://api.nemportfolio.in/auth/google/callback`
- **Vercel** → set `NEXT_PUBLIC_API_URL` and `NEXT_PUBLIC_AUTH_URL` to
  `https://api.nemportfolio.in`, then redeploy the frontend.

## Updating later

```
cd server && make build-linux
scp server/bin/server-linux <vps>:/tmp/server
ssh <vps> 'sudo mv /tmp/server /opt/trueflashcard/server && sudo systemctl restart trueflashcard'
```
