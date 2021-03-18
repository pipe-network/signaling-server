# signaling-server

This is the signaling server of the pipe network to connect peers.

It is using [SaltyRTC](https://github.com/saltyrtc/saltyrtc-meta) to establish a secure connection between two clients
over this signalling server.

# Install

Just install the dependencies from the go.mod file:

`go mod download`

# Start

Start the server by running 

```
go run main/main.go
```

Add following flags to configure:

```
http service address:
--address localhost:8080

public key file path:
--public_key_file ./public.key

private key file path:
--private_key_file ./private.key

TLS certificate file path:
--tls_cert_file ./cert.pem

TLS key file path:
--tls_key_file ./key.pem
```
