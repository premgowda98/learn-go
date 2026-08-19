# Makefile for mTLS Demo

.PHONY: help certs clean server client test demo

help: ## Show this help message
	@echo "🔐 mTLS Demo Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

certs: ## Generate all certificates for mTLS demo
	@echo "🔐 Generating certificates..."
	./generate-certs.sh

clean: ## Remove all generated certificates and binaries
	@echo "🧹 Cleaning up..."
	rm -rf certs/
	rm -f server client

server: certs ## Run the mTLS server
	@echo "🚀 Starting mTLS server..."
	cd server && go run server.go

client: certs ## Run the mTLS client
	@echo "📡 Starting mTLS client..."
	cd client && go run client.go

test: certs ## Run a quick test of the mTLS connection
	@echo "🧪 Testing mTLS connection..."
	@echo "Starting server in background..."
	@cd server && go run server.go &
	@SERVER_PID=$$!; \
	sleep 2; \
	echo "Running client test..."; \
	cd client && go run client.go; \
	echo "Stopping server..."; \
	kill $$SERVER_PID 2>/dev/null || true

demo: ## Run the full demo (server + client in separate terminals)
	@echo "🎭 Starting mTLS Demo..."
	@echo "This will open the server. Run 'make client' in another terminal."
	@make server

build: ## Build server and client binaries
	@echo "🔨 Building binaries..."
	cd server && go build -o ../server server.go
	cd client && go build -o ../client client.go
	@echo "✅ Built: server, client"

info: certs ## Show certificate information
	@echo "📋 Certificate Information:"
	@echo "================================"
	@echo ""
	@echo "🏛️  Root CA Certificate:"
	@openssl x509 -in certs/ca.crt -text -noout | head -20
	@echo ""
	@echo "🖥️  Server Certificate:"
	@openssl x509 -in certs/server.crt -text -noout | head -20
	@echo ""
	@echo "👤 Client Certificate:"
	@openssl x509 -in certs/client.crt -text -noout | head -20

verify: certs ## Verify certificate chain
	@echo "🔍 Verifying certificate chain..."
	@echo "Server certificate verification:"
	@openssl verify -CAfile certs/ca.crt certs/server.crt
	@echo "Client certificate verification:"
	@openssl verify -CAfile certs/ca.crt certs/client.crt

curl-test: certs ## Test with curl (this will fail due to mTLS requirement)
	@echo "🌐 Testing with curl (should fail - no client cert)..."
	@curl -v --cacert certs/ca.crt https://localhost:8443/hello || true
	@echo ""
	@echo "🌐 Testing with curl + client certificate..."
	@curl -v --cacert certs/ca.crt --cert certs/client.crt --key certs/client.key https://localhost:8443/hello
