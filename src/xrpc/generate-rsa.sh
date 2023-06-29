#!/bin/bash
set -ex

openssl genrsa -out certificates/ca.key 2048
openssl req -new -x509 -key certificates/ca.key -out certificates/ca.pem -days 3650 -subj "/CN=localhost"

openssl ecparam -genkey -name secp384r1 -out certificates/server.key
openssl req -new -key certificates/server.key -out certificates/server.csr -config generate-rsa.cnf -extensions SAN
openssl x509 -req -sha256 -CA certificates/ca.pem -CAkey certificates/ca.key -CAcreateserial -days 3650 -in certificates/server.csr -out certificates/server.pem -extfile generate-rsa.cnf -extensions SAN

openssl ecparam -genkey -name secp384r1 -out certificates/client.key
openssl req -new -key certificates/client.key -out certificates/client.csr -config generate-rsa.cnf -extensions SAN
openssl x509 -req -sha256 -CA certificates/ca.pem -CAkey certificates/ca.key -CAcreateserial -days 3650 -in certificates/client.csr -out certificates/client.pem -extfile generate-rsa.cnf -extensions SAN

# Common Name (e.g. server FQDN or YOUR name) []:localhost
