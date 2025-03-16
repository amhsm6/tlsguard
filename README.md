# tlsguard

Secures http endpoints through tls

## Usage

Expects **root.crt** and **root.key** in /etc/tlsguard.

**tlsguard -cert**  ---- Generates *client.crt* and *client.key*.

**tlsguard *secure_port*:*insecure_port*** ---- Creates a secure tunnel. TLS requests to **secure_port** with valid *client.crt* and *client.key* would be forwarded to *localhost:insecure_port*.
