#!/bin/bash
# Inspired from https://github.com/grpc/grpc-java/tree/master/examples#generating-self-signed-certificates-for-use-with-grpc

# Output files
# ca.key: certificate Authority private key file (this shouldn't be shared in real-life)
# ca.crt: certificate Authority trust certificate (this shouldn't be shared with users in real-life)
# server.key: Server private key, password protected (this shouldn't be shared)
# server.csr: server certificate signing request (this should be shared with CA owner)
# server.crt: Server certificate signed by the CA (this would be sent back by the CA owner) - keep on server.
# server.pem: Conversion of server.key into a format gRPC likes (this shouldn't be shared).

# summary:
# Private files: ca.key, server.key, server.pem, server.crt
# 'Share' files: ca.crt (needed by the client), server.csr (needed by the CA)

# Changes these CN's to match your hosts in your environment if needed
SERVER_CN=localhost

# step 1: Generate Certificate Authority + Trust Certificate (cs.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN={SERVER_CN}"

# step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN={SERVER_CN}"
# openssl req -passin pass:1111 -new -key server.key -reqexts SAN -config cert.conf <(printf "[SAN]\nsubjectAltName=DNS:localhost")) -out server.csr -subj "/CN={SERVER_CN}"

# step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt
# openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
# use the below statement if encountered the following issue 
# x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0".   
openssl x509 -req -passin pass:1111 -days 365 -extfile <(printf "subjectAltName=DNS:localhost") -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem 