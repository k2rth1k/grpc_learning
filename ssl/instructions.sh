#!/bin/sh

#Output files
#ca.key:       Certificate Authority private key file               (this should not be shared)
#ca.crt:       Certificate Authority trust certificate              (this shpuld be shared with users in real-life)
#server.key:   Server Private key, password protectted              (this shouldn't be shared)
#server.csr:   Server Certificate signed request                    (this should be shared with CA owner)
#server.crt:   Server certificate signed by the CA                  (this would be sent back by the CA owner) - keep on server
#server.pem:   Conversion of server.key into format grpc likes      (this shouldn't be shared)

# Summary
# Private files: ca.key, server.key, server.pem, server.crt
# share files:   ca.crt(needed by client), server.csr (needed by the CA)

# Changes these CN's to match your hosts in your environment if needed.
SERVER_CN=localhost

# Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=$SERVER_CN"

# STEP 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a Certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=$SERVER_CN"

# Step 4: Sign the certificate with the CA we created (it's called self signing) -server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by grpc
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
