## 🔐 What Is Certificate Pinning?

**Certificate pinning** (also called **SSL pinning**) is a technique used by applications (especially mobile and Go-based apps) to **hardcode trust** for a **specific certificate or public key**, instead of relying solely on the system’s list of trusted Certificate Authorities (CAs).

---

### 🧠 Why Do This?

Normally, HTTPS works like this:

1. A client connects to `https://example.com`
2. The server provides a TLS certificate
3. The client checks:

   * Is the certificate valid?
   * Was it issued by a trusted CA? (from OS cert store)
   * Does the hostname match?

✔️ If all good, TLS handshake succeeds.

---

### 🤯 But with MITM proxies...

When you intercept HTTPS (e.g., using `mitmproxy` or `Squid + SSL Bump`), the proxy **generates a fake certificate** for `example.com` — signed by **your own CA**, not the real one.

If the app uses **system CA trust**, you can install your fake CA and it’ll work fine.

❌ But if the app uses **certificate pinning**, it **knows exactly what the cert (or its public key) should be**, and **refuses any substitute**, even if your proxy cert is technically valid.

---

## 📌 Types of Certificate Pinning

### 1. **Public Key Pinning (Most Common)**

The app stores the **public key hash** of the certificate it expects:

```go
expectedPubKey := "sha256/3hIpM7AoYz+T9dM8qFvgk="
```

When connecting to `example.com`, it compares the server’s actual public key to the stored one.

If it **doesn’t match** → the connection fails.

This is what Go’s `net/http` or `crypto/tls` packages can be configured to do.

---

### 2. **SPKI Pinning (Subject Public Key Info)**

Same as above, but pins the entire SPKI structure (more specific).

---

### 3. **CA Pinning**

The app expects the certificate to be signed by a **specific CA** (e.g., Let’s Encrypt), not just any CA.

Even if your mitmproxy CA is trusted, it’s still rejected because it’s not *that* CA.

---

### 4. **Self-Signed Pinning**

Some apps come with a **self-signed cert** (often embedded), and **only trust that one cert** — common in IoT, internal APIs, or custom services.

---

## ⚠️ What Happens When Certificate Pinning Is Used?

| Attempted MITM              | Result                                                                          |
| --------------------------- | ------------------------------------------------------------------------------- |
| Proxy injects its fake cert | App compares it to pinned cert/key                                              |
| Mismatch found              | TLS fails with error like `x509: certificate is not trusted` or `pinning error` |

---

## 🧰 Real-World Examples

### ✅ Works with MITM:

* Browsers like Chrome, Firefox (with system CA installed)
* `curl` (respects system trust)
* Python, Java, Node.js apps

### ❌ Fails with MITM:

* Mobile banking apps
* Many Go apps that use `tls.Config` with pinned keys
* Telegram, WhatsApp, and some signal apps
* Apps using certificate pinning libraries (e.g., TrustKit, OkHttp pinning)

---

## 🧰 How Pinning Is Done (Go Example)

```go
tlsConfig := &tls.Config{
    InsecureSkipVerify: true, // skip regular validation
    VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
        // Manually verify cert pin
        pubKeyHash := sha256.Sum256(rawCerts[0])
        if base64.StdEncoding.EncodeToString(pubKeyHash[:]) != expectedHash {
            return errors.New("certificate pin mismatch")
        }
        return nil
    },
}
```

---

## 🛠️ Can You Bypass Certificate Pinning?

### 1. **Patch the Binary (advanced)**

* Modify the binary to skip or fake the pin check.
* Use tools like:

  * `frida` (runtime hook)
  * `objection` (for mobile)
  * `gdb` or `radare2` (manual patching)

### 2. **Rebuild the App (if source is available)**

* Remove or disable the pinning logic.

### 3. **Custom Proxy Like Frida-SSL-Unpinning**

* Hook TLS libraries at runtime and bypass checks.

---

## ✅ TL;DR: Certificate Pinning

| Concept             | Meaning                                                             |
| ------------------- | ------------------------------------------------------------------- |
| What is it?         | App hardcodes specific cert or key to trust                         |
| Why use it?         | Prevent MITM, stop trust-on-first-use issues                        |
| What does it block? | Proxies like mitmproxy, even with valid CAs                         |
| Who uses it?        | Secure apps (banking, finance, privacy apps), Go binaries           |
| Can you bypass it?  | Yes, but requires binary patching, runtime hooking, or source edits |