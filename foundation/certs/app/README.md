# 🔐 Mutual TLS (mTLS) with Certificate Pinning Demo

This is a comprehensive Go application demonstrating **Mutual TLS (mTLS)** authentication with **certificate pinning**. The demo includes a TLS server that requires client certificates and a Go client that verifies server certificates using certificate pinning.

## 🎯 What This Demo Shows

- **Mutual TLS (mTLS)**: Both client and server authenticate each other using certificates
- **Certificate Pinning**: Client verifies server's public key hash to prevent man-in-the-middle attacks
- **Custom PKI**: Complete certificate authority setup with self-signed certificates
- **Secure Communication**: HTTPS API with JSON endpoints over mTLS

## 📁 Project Structure

```
certs/app/
├── server.go              # mTLS server with client certificate verification
├── client.go              # mTLS client with certificate pinning
├── generate-certs.sh      # Certificate generation script
├── Makefile              # Build and run commands
├── go.mod                # Go module file
└── certs/                # Generated certificates (created by script)
    ├── ca.crt            # Root CA certificate
    ├── ca.key            # Root CA private key
    ├── server.crt        # Server certificate
    ├── server.key        # Server private key
    ├── client.crt        # Client certificate
    └── client.key        # Client private key
```

## 🚀 Quick Start

### 1. Generate Certificates
```bash
make certs
```

### 2. Run the Server (Terminal 1)
```bash
make server
```

### 3. Run the Client (Terminal 2)
```bash
make client
```

### 4. Or run automated test
```bash
make test
```

## 🔧 Available Commands

Run `make help` to see all available commands:

- `make certs` - Generate all certificates
- `make server` - Start the mTLS server  
- `make client` - Run the mTLS client
- `make test` - Automated test (server + client)
- `make clean` - Remove certificates and binaries
- `make info` - Show certificate details
- `make verify` - Verify certificate chain
- `make curl-test` - Test with curl

## 🌐 API Endpoints

The server exposes these endpoints:

- `GET /hello` - Simple greeting with client certificate info
- `POST /echo` - Echo JSON payload back to client
- `GET /info` - Detailed TLS and certificate information

## 🔍 Certificate Information

The demo generates:

1. **Root CA Certificate** (`ca.crt`, `ca.key`)
   - Self-signed certificate authority
   - Used to sign both server and client certificates
   - Valid for 10 years

2. **Server Certificate** (`server.crt`, `server.key`)
   - CN: `mtls-server`
   - SAN: `localhost`, `127.0.0.1`
   - Signed by Root CA
   - Valid for 1 year

3. **Client Certificate** (`client.crt`, `client.key`)
   - CN: `mtls-client`
   - Extended Key Usage: Client Authentication
   - Signed by Root CA
   - Valid for 1 year

## 🔐 Security Features

### Server Security
- Requires valid client certificates (mTLS)
- Verifies client certificates against custom CA
- Uses TLS 1.2+ only
- Includes security headers and timeouts

### Client Security
- Verifies server certificate chain
- Implements certificate pinning (SHA-256 public key hash)
- Uses custom CA for server verification
- Prevents man-in-the-middle attacks

## 🧪 Testing Different Scenarios

### Test Certificate Pinning
1. Start the server normally
2. Replace `server.crt` with a different certificate
3. Run the client - it should fail with pinning error

### Test Client Certificate Requirement
```bash
# This will fail - no client certificate
curl -k https://localhost:8443/hello

# This will work - with client certificate
curl --cacert certs/ca.crt --cert certs/client.crt --key certs/client.key https://localhost:8443/hello
```

### Test with Invalid Certificates
1. Generate new certificates: `make clean && make certs`
2. Replace only the client certificate with an old one
3. The server should reject the client

## 🔧 Manual Certificate Generation

If you prefer to generate certificates manually:

```bash
# 1. Generate CA
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt

# 2. Generate Server Certificate
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -out server.crt -days 365

# 3. Generate Client Certificate  
openssl genrsa -out client.key 4096
openssl req -new -key client.key -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -out client.crt -days 365
```

## 🛡️ Security Considerations

### In Production
- Store private keys securely (HSM, key vault)
- Use proper certificate rotation policies
- Implement certificate revocation (CRL/OCSP)
- Monitor certificate expiration
- Use certificates from trusted CAs

### Certificate Pinning
- Pin specific certificates or public keys
- Have a backup pinning strategy
- Plan for certificate rotation
- Consider pinning multiple certificates

### mTLS Best Practices
- Use strong cipher suites
- Implement proper certificate validation
- Use TLS 1.2 or higher
- Regular security audits

## 🔍 Troubleshooting

### Common Issues

1. **Certificate not found errors**
   ```
   Solution: Run `make certs` to generate certificates
   ```

2. **Connection refused**
   ```
   Solution: Make sure server is running on port 8443
   ```

3. **Certificate pinning failures**
   ```
   Solution: Update expectedServerPubKeyHash in client.go
   ```

4. **Permission denied on scripts**
   ```
   Solution: chmod +x generate-certs.sh
   ```

### Debug Commands

```bash
# Check certificate details
openssl x509 -in certs/server.crt -text -noout

# Verify certificate chain
openssl verify -CAfile certs/ca.crt certs/server.crt

# Test TLS connection
openssl s_client -connect localhost:8443 -cert certs/client.crt -key certs/client.key
```

## 📚 Educational Value

This demo teaches:
- TLS/SSL fundamentals
- Certificate-based authentication
- Public Key Infrastructure (PKI)
- Man-in-the-middle attack prevention
- Go crypto/tls package usage
- OpenSSL certificate management

Perfect for security training, penetration testing setup, or understanding enterprise authentication systems.
