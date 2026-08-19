package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverPort     = ":8443"
	serverCertFile = "../certs/server.crt"
	serverKeyFile  = "../certs/server.key"
	caCertFile     = "../certs/ca.crt"
)

type ServerResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	ClientCN  string    `json:"client_cn,omitempty"`
}

type EchoRequest struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Load server certificate and key
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("❌ Failed to load server certificate: %v", err)
	}

	// Load CA certificate for client verification
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("❌ Failed to read CA certificate: %v", err)
	}

	// Create CA certificate pool for client certificate verification
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("❌ Failed to append CA certificate to pool")
	}

	// Create TLS configuration requiring mutual authentication
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert}, // Server certificate
		ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificate
		ClientCAs:    caCertPool,                     // CA to verify client certificates
		MinVersion:   tls.VersionTLS12,               // Minimum TLS version
	}

	// Create HTTP server with TLS configuration
	server := &http.Server{
		Addr:      serverPort,
		TLSConfig: tlsConfig,
		Handler:   createRouter(),
		// Add timeouts for security
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("🚀 Starting mTLS server...")
	fmt.Printf("🔒 Server listening on https://localhost%s\n", serverPort)
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET  /hello - Simple greeting")
	fmt.Println("   POST /echo  - Echo JSON payload")
	fmt.Println("   GET  /info  - Server and client certificate info")
	fmt.Println("🔐 Server requires mutual TLS authentication")
	fmt.Println()

	// Start the HTTPS server
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func createRouter() http.Handler {
	mux := http.NewServeMux()

	// Add middleware to log requests and extract client certificate info
	mux.HandleFunc("/hello", logRequest(handleHello))
	mux.HandleFunc("/echo", logRequest(handleEcho))
	mux.HandleFunc("/info", logRequest(handleInfo))

	return mux
}

// Middleware to log requests and extract client certificate information
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract client certificate information if available
		var clientCN string
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			clientCert := r.TLS.PeerCertificates[0]
			clientCN = clientCert.Subject.CommonName
		}

		fmt.Printf("📥 %s %s from client: %s (CN: %s)\n", r.Method, r.URL.Path, r.RemoteAddr, clientCN)
		
		// Call the next handler
		next(w, r)
	}
}

// Simple hello endpoint
func handleHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clientCN string
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		clientCN = r.TLS.PeerCertificates[0].Subject.CommonName
	}

	response := ServerResponse{
		Message:   fmt.Sprintf("Hello from mTLS server! Welcome, %s", clientCN),
		Timestamp: time.Now(),
		ClientCN:  clientCN,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Echo endpoint that returns the received JSON
func handleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clientCN string
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		clientCN = r.TLS.PeerCertificates[0].Subject.CommonName
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON request
	var echoReq EchoRequest
	if err := json.Unmarshal(body, &echoReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create response with echo data
	response := map[string]interface{}{
		"echo_message":     echoReq.Message,
		"echo_timestamp":   echoReq.Timestamp,
		"server_response":  "Message received successfully via mTLS",
		"server_timestamp": time.Now().Format(time.RFC3339),
		"client_cn":        clientCN,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Info endpoint that provides detailed certificate information
func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info := map[string]interface{}{
		"server": map[string]interface{}{
			"message":      "mTLS Server Information",
			"tls_version":  getTLSVersionString(r.TLS.Version),
			"cipher_suite": getCipherSuiteString(r.TLS.CipherSuite),
		},
	}

	// Add client certificate information if available
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		clientCert := r.TLS.PeerCertificates[0]
		info["client"] = map[string]interface{}{
			"common_name":    clientCert.Subject.CommonName,
			"organization":   clientCert.Subject.Organization,
			"country":        clientCert.Subject.Country,
			"serial_number":  clientCert.SerialNumber.String(),
			"not_before":     clientCert.NotBefore.Format(time.RFC3339),
			"not_after":      clientCert.NotAfter.Format(time.RFC3339),
			"is_ca":          clientCert.IsCA,
			"key_usage":      getKeyUsageString(clientCert.KeyUsage),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(w).Encode(info)
}

// Helper function to get TLS version string
func getTLSVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (%d)", version)
	}
}

// Helper function to get cipher suite string
func getCipherSuiteString(cipherSuite uint16) string {
	switch cipherSuite {
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
	default:
		return fmt.Sprintf("Unknown (%d)", cipherSuite)
	}
}

// Helper function to get key usage string
func getKeyUsageString(keyUsage x509.KeyUsage) []string {
	var usages []string
	
	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "Digital Signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "Content Commitment")
	}
	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "Key Encipherment")
	}
	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "Data Encipherment")
	}
	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "Key Agreement")
	}
	if keyUsage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "Certificate Sign")
	}
	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL Sign")
	}
	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "Encipher Only")
	}
	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "Decipher Only")
	}
	
	return usages
}
