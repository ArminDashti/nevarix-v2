#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CERT_DIR="${ROOT}/certs"
DAYS=365

mkdir -p "${CERT_DIR}"
cd "${CERT_DIR}"

if [[ -f ca.crt ]]; then
  echo "Certificates already exist in ${CERT_DIR}"
  exit 0
fi

# CA
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days "${DAYS}" \
  -subj "/CN=Nevarix Dev CA" -out ca.crt

gen_component() {
  local name="$1"
  openssl genrsa -out "${name}.key" 2048
  openssl req -new -key "${name}.key" \
    -subj "/CN=nvx-${name}" -out "${name}.csr"
  openssl x509 -req -in "${name}.csr" -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out "${name}.crt" -days "${DAYS}" -sha256
  rm "${name}.csr"
}

gen_component hub
gen_component agent
gen_component manager

echo "Dev TLS certificates written to ${CERT_DIR}"
