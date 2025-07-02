package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	serverAddress = "localhost:8443"
	// This will be the expected server certificate's public key hash
	// Generated from: openssl x509 -in certs/server.crt -pubkey -noout | openssl pkey -pubin -outform DER | openssl dgst -sha256 -hex
	expectedServerPubKeyHash = "ccce7dfa888af9f9fdb088476fe5c9c96d14d49ab1d5276f4f17977f0914aeef"
)

func main() {
	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	// Load CA certificate for server verification
	caCert, err := ioutil.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	// Create CA certificate pool
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to append CA certificate to pool")
	}

	// Create custom TLS configuration with mTLS and certificate pinning
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert}, // Client certificate for mTLS
		RootCAs:      caCertPool,                    // Custom CA for server verification
		ServerName:   "mtls-server",                 // Must match server certificate CN
		// Custom certificate verification with pinning
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return verifyServerCertificateWithPinning(rawCerts, verifiedChains)
		},
		// We still want to verify the certificate chain normally
		InsecureSkipVerify: false,
	}

	// Create HTTP client with custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 10 * time.Second,
	}

	// Test the connection
	fmt.Println("🚀 Starting mTLS client...")
	fmt.Printf("📡 Connecting to https://%s\n", serverAddress)

	// Test GET request
	testGetRequest(client)

	// Test POST request
	testPostRequest(client)
}

// verifyServerCertificateWithPinning implements certificate pinning
func verifyServerCertificateWithPinning(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(rawCerts) == 0 {
		return fmt.Errorf("no certificates provided")
	}

	// Parse the server certificate
	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return fmt.Errorf("failed to parse server certificate: %v", err)
	}

	// Calculate the SHA-256 hash of the server's public key
	pubKeyDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal server public key: %v", err)
	}

	hash := sha256.Sum256(pubKeyDER)
	actualHash := hex.EncodeToString(hash[:])

	fmt.Printf("🔑 Server certificate public key hash: %s\n", actualHash)

	// For demonstration, we'll print the hash so you can update the constant
	// In production, you would have the expected hash hardcoded
	if expectedServerPubKeyHash == "REPLACE_WITH_ACTUAL_HASH" {
		fmt.Printf("⚠️  Update expectedServerPubKeyHash constant with: %s\n", actualHash)
		fmt.Println("✅ Certificate pinning check passed (development mode)")
		return nil
	}

	// Verify the pinned public key hash
	if actualHash != expectedServerPubKeyHash {
		return fmt.Errorf("🚨 CERTIFICATE PINNING FAILED: expected %s, got %s", expectedServerPubKeyHash, actualHash)
	}

	fmt.Println("✅ Certificate pinning verification passed!")
	return nil
}

func testGetRequest(client *http.Client) {
	fmt.Println("\n📥 Testing GET request...")

	resp, err := client.Get(fmt.Sprintf("https://%s/hello", serverAddress))
	if err != nil {
		log.Printf("❌ GET request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response body: %v", err)
		return
	}

	fmt.Printf("✅ GET Response Status: %s\n", resp.Status)
	fmt.Printf("📄 Response Body: %s\n", string(body))
}

func testPostRequest(client *http.Client) {
	fmt.Println("\n📤 Testing POST request...")

	postData := `{"message": "Hello from mTLS client!", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/echo", serverAddress),
		strings.NewReader(postData))
	if err != nil {
		log.Printf("❌ Failed to create POST request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ POST request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read POST response body: %v", err)
		return
	}

	fmt.Printf("✅ POST Response Status: %s\n", resp.Status)
	fmt.Printf("📄 Response Body: %s\n", string(responseBody))
}
