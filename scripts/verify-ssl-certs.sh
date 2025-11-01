#!/bin/bash

# Script to verify SSL certificate installation and configuration

echo "=== SSL Certificate Verification ==="
echo ""

# Check if certificate files exist
if [ ! -f "./tls/localhost.pem" ]; then
    echo "❌ Certificate file not found: ./tls/localhost.pem"
    exit 1
fi

if [ ! -f "./tls/localhost-key.pem" ]; then
    echo "❌ Private key file not found: ./tls/localhost-key.pem"
    exit 1
fi

echo "✅ Certificate files exist"
echo ""

# Check certificate validity
echo "Certificate Information:"
echo "------------------------"
openssl x509 -in ./tls/localhost.pem -noout -subject -issuer -dates
echo ""

# Check for SAN extension
echo "Subject Alternative Names:"
echo "-------------------------"
openssl x509 -in ./tls/localhost.pem -text -noout | grep -A 2 "Subject Alternative Name"
echo ""

# Check if certificate is in Chrome's trust store
echo "Browser Trust Store Status:"
echo "---------------------------"
if command -v certutil &> /dev/null; then
    if certutil -d sql:$HOME/.pki/nssdb -L 2>/dev/null | grep -q "localhost-dev-scvis"; then
        echo "✅ Certificate is installed in Chrome/Chromium trust store"
    else
        echo "⚠️  Certificate is NOT in Chrome/Chromium trust store"
        echo "   Run: certutil -d sql:\$HOME/.pki/nssdb -A -t \"P,,\" -n \"localhost-dev-scvis\" -i ./tls/localhost.pem"
    fi
else
    echo "⚠️  certutil not found - cannot check Chrome trust store"
    echo "   Install with: sudo apt-get install libnss3-tools"
fi
echo ""

# Verify certificate and key match
echo "Certificate and Key Match:"
echo "-------------------------"
cert_modulus=$(openssl x509 -noout -modulus -in ./tls/localhost.pem | openssl md5)
key_modulus=$(openssl rsa -noout -modulus -in ./tls/localhost-key.pem 2>/dev/null | openssl md5)

if [ "$cert_modulus" = "$key_modulus" ]; then
    echo "✅ Certificate and private key match"
else
    echo "❌ Certificate and private key do NOT match!"
    exit 1
fi
echo ""

echo "=== Verification Complete ==="
