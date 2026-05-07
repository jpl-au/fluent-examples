#!/usr/bin/env bash
# Build WASM artefacts from tether-wasm and copy into static/ for serving.
# By default builds stdlib Go. Pass "tinygo" to build the TinyGo variant.
set -euo pipefail

cd "$(dirname "$0")"

WASM_MODULE="../../tether-wasm"
VARIANT="${1:-go}"

echo "=== Building WASM ($VARIANT) ==="
if [ "$VARIANT" = "tinygo" ]; then
    (cd "$WASM_MODULE" && bash build-tinygo.sh)
    cp "$WASM_MODULE/client.tinygo.wasm" static/client.wasm
else
    (cd "$WASM_MODULE" && bash build-go.sh)
    cp "$WASM_MODULE/client.go.wasm" static/client.wasm
fi

# Copy the matching wasm_exec.js. Each toolchain requires its own.
if [ "$VARIANT" = "tinygo" ]; then
    cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" static/wasm_exec.js
else
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" static/wasm_exec.js
fi

echo
echo "static/ contents:"
ls -lh static/
echo
echo "Ready. Start the server with: go run ."
