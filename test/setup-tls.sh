# Create a CA key:
openssl genrsa -out ca.key.pem 2048

# Create a CA certificate:
openssl req -x509 -new -nodes -key ca.key.pem -subj "/CN=CARoot" -days 365 -out ca.cert.pem

# Generate the server's private key:
openssl genrsa -out broker.key.pem 2048

# Convert server private key to PKCS#8 format:
openssl pkcs8 -topk8 -inform PEM -outform PEM -in broker.key.pem -out broker.key-pk8.pem -nocrypt

# Generate certificate request for the server:
openssl req -new -config broker.conf -key broker.key.pem -out broker.csr.pem -sha256

# Sign server certificate with the CA:
openssl x509 -req -in broker.csr.pem -CA ca.cert.pem -CAkey ca.key.pem -CAcreateserial -out broker.cert.pem -days 365 -extensions v3_ext -extfile broker.conf -sha256

# Generate the client's private key:
openssl genrsa -out client.key.pem 2048

# Convert client private key to PKCS#8 format:
openssl pkcs8 -topk8 -inform PEM -outform PEM -in client.key.pem -out client.key-pk8.pem -nocrypt

# Generate certificate request for the client:
openssl req -new -subj "/CN=client" -key client.key.pem -out client.csr.pem -sha256

# Sign client certificate with the CA:
openssl x509 -req -in client.csr.pem -CA ca.cert.pem -CAkey ca.key.pem -CAcreateserial -out client.cert.pem -days 365 -sha256

chmod 644 *.pem
