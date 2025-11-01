# SSL Certificate Setup Complete

## Summary

New SSL certificates have been generated for localhost development with proper Subject Alternative Name (SAN) configuration. These certificates are valid for 365 days and have been automatically installed in your Chrome/Chromium trust store.

## What Was Done

1. ✅ Generated new SSL certificates with proper SAN extensions
   - Certificate: `tls/localhost.pem`
   - Private Key: `tls/localhost-key.pem`
   - Valid for: 365 days (until Oct 31, 2026)

2. ✅ Configured Subject Alternative Names:
   - DNS: localhost
   - DNS: *.localhost
   - IP: 127.0.0.1
   - IP: ::1 (IPv6 localhost)

3. ✅ Installed certificate in Chrome/Chromium trust store
   - Certificate name: `localhost-dev-scvis`

4. ✅ Created helper scripts:
   - `scripts/generate-ssl-certs.sh` - Regenerate certificates
   - `scripts/verify-ssl-certs.sh` - Verify certificate installation

## Next Steps

1. **Restart your browser completely** (close all Chrome windows)
   
2. **Restart the SCVis application:**
   ```bash
   sudo systemctl restart scvisual
   ```
   Or if running manually:
   ```bash
   go build ./cmd/web && ./web
   ```

3. **Clear browser cache for localhost** (optional but recommended):
   - Open DevTools (F12)
   - Right-click the refresh button
   - Select "Empty Cache and Hard Reload"

4. **Visit your application:**
   - Go to https://localhost:8082/scvis/
   - You should now see a secure connection (🔒) without warnings
   - The Flutter service worker should load without SSL errors

## Verification

Run the verification script to check everything is set up correctly:
```bash
./scripts/verify-ssl-certs.sh
```

## Troubleshooting

### Still seeing "Not Secure" or certificate warnings?

1. Make sure you completely closed and reopened Chrome (not just the tab)
2. Verify the certificate is installed:
   ```bash
   certutil -d sql:$HOME/.pki/nssdb -L | grep localhost
   ```
3. Check the certificate details in Chrome:
   - Click the 🔒 or ⚠️ icon in the address bar
   - Click "Certificate is valid" or "Certificate"
   - Verify it shows "localhost-dev-scvis"

### Service worker still failing?

1. Clear all site data:
   - DevTools → Application tab → Storage → Clear site data
2. Completely restart the browser
3. Try accessing the site in an incognito/private window first

### Need to regenerate certificates?

```bash
./scripts/generate-ssl-certs.sh
certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n "localhost-dev-scvis" -i ./tls/localhost.pem
```

## For Other Team Members

If other developers need to set up their environment:

1. Pull the latest code (includes the new certificates)
2. Run the setup:
   ```bash
   ./scripts/generate-ssl-certs.sh
   certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n "localhost-dev-scvis" -i ./tls/localhost.pem
   ```
3. Restart browser completely

## Documentation

See `tls/README.md` for complete documentation on SSL certificate management.

---

**Note:** These are self-signed certificates for development only. The production server at `vaccs.ceas.wmich.edu` uses proper Let's Encrypt certificates.
