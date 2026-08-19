#!/bin/bash

# Certificate Generation Script for mTLS Demo
# This script creates a complete PKI setup with:
# - Root CA certificate
# - Server certificate (signed by CA)
# - Client certificate (signed by CA)

set -e  # Exit on any error

echo "🔐 Generating certificates for mTLS demo..."

# Create certs directory if it doesn't exist
mkdir -p certs
cd certs

# Clean up any existing certificates
rm -f *.crt *.key *.csr *.srl

echo "📋 Step 1: Generating Root CA private key..."
# Generate CA private key (4096-bit RSA)
openssl genrsa -out ca.key 4096

echo "📋 Step 2: Creating Root CA certificate..."
# Create CA certificate (self-signed, valid for 10 years)
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/C=US/ST=California/L=San Francisco/O=mTLS Demo CA/OU=Security/CN=mTLS-Root-CA/emailAddress=ca@mtls-demo.local"

echo "📋 Step 3: Generating Server private key..."
# Generate server private key
openssl genrsa -out server.key 4096

echo "📋 Step 4: Creating Server certificate signing request..."
# Create server certificate signing request
openssl req -new -key server.key -out server.csr -subj "/C=US/ST=California/L=San Francisco/O=mTLS Demo/OU=Server/CN=mtls-server/emailAddress=server@mtls-demo.local"

echo "📋 Step 5: Creating Server certificate (signed by CA)..."
# Create server certificate extensions file
cat > server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = mtls-server
DNS.2 = localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

# Sign server certificate with CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile server.ext

echo "📋 Step 6: Generating Client private key..."
# Generate client private key
openssl genrsa -out client.key 4096

echo "📋 Step 7: Creating Client certificate signing request..."
# Create client certificate signing request
openssl req -new -key client.key -out client.csr -subj "/C=US/ST=California/L=San Francisco/O=mTLS Demo/OU=Client/CN=mtls-client/emailAddress=client@mtls-demo.local"

echo "📋 Step 8: Creating Client certificate (signed by CA)..."
# Create client certificate extensions file
cat > client.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
EOF

# Sign client certificate with CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -extfile client.ext

# Clean up temporary files
rm -f *.csr *.ext

echo "✅ Certificate generation complete!"
echo ""
echo "📁 Generated files:"
echo "   ca.crt      - Root CA certificate (public)"
echo "   ca.key      - Root CA private key (keep secure!)"
echo "   server.crt  - Server certificate (public)"
echo "   server.key  - Server private key (keep secure!)"
echo "   client.crt  - Client certificate (public)"
echo "   client.key  - Client private key (keep secure!)"
echo ""
echo "🔍 Certificate Information:"
echo ""
echo "📋 Root CA Certificate:"
openssl x509 -in ca.crt -text -noout | grep -A 5 "Subject:"
echo ""
echo "📋 Server Certificate:"
openssl x509 -in server.crt -text -noout | grep -A 5 "Subject:"
echo ""
echo "📋 Client Certificate:"
openssl x509 -in client.crt -text -noout | grep -A 5 "Subject:"
echo ""
echo "🔐 Server Certificate Public Key Hash (for certificate pinning):"
# Calculate the SHA-256 hash of the server's public key
openssl x509 -in server.crt -pubkey -noout | openssl pkey -pubin -outform DER | openssl dgst -sha256 -binary | openssl enc -base64
echo ""
echo "🚀 Ready to run mTLS demo!"
echo "   1. Run server: go run server.go"
echo "   2. Run client: go run client.go"
