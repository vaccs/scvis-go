# SSL/TLS Certificates for SCVis Development

This directory contains SSL certificates for running the SCVis application over HTTPS during development.

## Files

- `localhost.pem` - SSL certificate for localhost
- `localhost-key.pem` - Private key for the certificate

## Generating New Certificates

To generate new certificates (valid for 365 days):

```bash
./scripts/generate-ssl-certs.sh
```

## Installing the Certificate in Your Browser

The certificates are self-signed, so you need to trust them in your browser to avoid SSL warnings.

### Chrome/Chromium on Linux

**Method 1: Using certutil (Automatic)**
```bash
certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n "localhost-dev-scvis" -i ./tls/localhost.pem
```

**Method 2: Using Chrome Settings (Manual)**
1. Open Chrome and go to `chrome://settings/certificates`
2. Click on the **Authorities** tab
3. Click **Import**
4. Select `tls/localhost.pem`
5. Check **Trust this certificate for identifying websites**
6. Click **OK**
7. Restart Chrome

### Firefox

1. Open Firefox and go to `about:preferences#privacy`
2. Scroll down to **Certificates** and click **View Certificates**
3. Click on the **Authorities** tab
4. Click **Import**
5. Select `tls/localhost.pem`
6. Check **Trust this CA to identify websites**
7. Click **OK**
8. Restart Firefox

### Edge on Linux

Same as Chrome (uses the same certificate store).

## Verifying the Certificate

To view certificate details:
```bash
openssl x509 -in tls/localhost.pem -text -noout
```

To check if the certificate has the proper Subject Alternative Name (SAN):
```bash
openssl x509 -in tls/localhost.pem -text -noout | grep -A 3 "Subject Alternative Name"
```

You should see:
```
X509v3 Subject Alternative Name: 
    DNS:localhost, DNS:*.localhost, IP Address:127.0.0.1, IP Address:0:0:0:0:0:0:0:1
```

## Removing the Certificate from Browser Trust Store

If you need to remove the certificate:

### Chrome/Chromium
```bash
certutil -d sql:$HOME/.pki/nssdb -D -n "localhost-dev-scvis"
```

Or use Chrome settings to manually delete it from the Authorities tab.

### Firefox
Use `about:preferences#privacy` → View Certificates → Authorities tab → Find and delete the certificate.

## Certificate Details

- **Common Name (CN)**: localhost
- **Organization (O)**: VACCS
- **Organizational Unit (OU)**: Development
- **Locality (L)**: Kalamazoo
- **State (ST)**: Michigan
- **Country (C)**: US
- **Valid For**: 365 days from generation date
- **Key Size**: 2048 bits RSA
- **Signature Algorithm**: SHA-256

## Troubleshooting

### "Your connection is not private" / NET::ERR_CERT_AUTHORITY_INVALID

This means the certificate is not trusted. Follow the installation steps above.

### Service Worker Registration Fails

Service workers require valid HTTPS. After installing the certificate:
1. Close all browser windows completely
2. Reopen the browser
3. Clear site data for localhost:8082 if needed
4. Try accessing the site again

### Certificate Not Trusted After Installation

- Make sure you restarted the browser completely (close all windows/tabs)
- Verify the certificate is installed: `certutil -d sql:$HOME/.pki/nssdb -L | grep localhost`
- Check that you selected "Trust this certificate for identifying websites"

## Production Certificates

For production deployment on `vaccs.ceas.wmich.edu`, use Let's Encrypt certificates located at:
- Certificate: `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem`
- Private Key: `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem`

These are automatically configured when `LOCAL_SERVER=false` in the `.env` file.
