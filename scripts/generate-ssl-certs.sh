#!/bin/bash

# Script to generate self-signed SSL certificates for localhost
# These certificates will have proper SAN extensions for modern browsers

CERT_DIR="./tls"
DAYS_VALID=365

# Create tls directory if it doesn't exist
mkdir -p "$CERT_DIR"

# Generate OpenSSL configuration file with SAN
cat > "$CERT_DIR/openssl.cnf" << EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = v3_req

[dn]
C = US
ST = Michigan
L = Kalamazoo
O = VACCS
OU = Development
CN = localhost

[v3_req]
subjectAltName = @alt_names
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

echo "Generating private key..."
openssl genrsa -out "$CERT_DIR/localhost-key.pem" 2048

echo "Generating certificate signing request..."
openssl req -new -key "$CERT_DIR/localhost-key.pem" \
    -out "$CERT_DIR/localhost.csr" \
    -config "$CERT_DIR/openssl.cnf"

echo "Generating self-signed certificate with SAN..."
openssl x509 -req -days $DAYS_VALID \
    -in "$CERT_DIR/localhost.csr" \
    -signkey "$CERT_DIR/localhost-key.pem" \
    -out "$CERT_DIR/localhost.pem" \
    -extensions v3_req \
    -extfile "$CERT_DIR/openssl.cnf"

# Clean up temporary files
rm "$CERT_DIR/localhost.csr"
rm "$CERT_DIR/openssl.cnf"

echo ""
echo "SSL certificates generated successfully!"
echo "Certificate: $CERT_DIR/localhost.pem"
echo "Private key: $CERT_DIR/localhost-key.pem"
echo "Valid for: $DAYS_VALID days"
echo ""
echo "To trust this certificate in Chrome/Chromium:"
echo "1. Go to chrome://settings/certificates"
echo "2. Click 'Authorities' tab"
echo "3. Click 'Import'"
echo "4. Select $CERT_DIR/localhost.pem"
echo "5. Check 'Trust this certificate for identifying websites'"
echo "6. Click OK"
echo ""
echo "Alternatively, for Linux with Chrome, run:"
echo "  certutil -d sql:\$HOME/.pki/nssdb -A -t \"P,,\" -n \"localhost-dev\" -i $CERT_DIR/localhost.pem"
echo ""
echo "For Firefox:"
echo "1. Go to about:preferences#privacy"
echo "2. Scroll to 'Certificates' -> View Certificates"
echo "3. Click 'Authorities' tab"
echo "4. Click 'Import'"
echo "5. Select $CERT_DIR/localhost.pem"
echo "6. Check 'Trust this CA to identify websites'"
echo ""
