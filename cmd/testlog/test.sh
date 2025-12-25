#!/bin/bash

cd "$(dirname "$0")"

echo "=== Direct to terminal (should have colors) ==="
TL_LOG=TRACE go run main.go

echo ""
echo "=== Piped through cat (no colors) ==="
TL_LOG=TRACE go run main.go 2>&1 | cat

echo ""
echo "=== Piped through tee (no colors) ==="
TL_LOG=TRACE go run main.go 2>&1 | tee /tmp/testlog_output.txt

echo ""
echo "=== Redirected to file (no colors) ==="
TL_LOG=TRACE go run main.go > /tmp/testlog_redirect.txt 2>&1
cat /tmp/testlog_redirect.txt

echo ""
echo "=== JSON mode (TL_LOG_JSON=true) ==="
TL_LOG=TRACE TL_LOG_JSON=true go run main.go

echo ""
echo "=== Default log level (INFO) ==="
go run main.go
