#!/usr/bin/env bash
set -euo pipefail

# Test environment deploy wrapper.
# Defaults can be overridden via env vars at runtime.
PROJECT_ID="${PROJECT_ID:-koala-map-jp-test}"
REGION="${REGION:-asia-northeast1}"
API_SERVICE="${API_SERVICE:-koala-map-jp-api-test}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Deploying to TEST"
echo "  PROJECT_ID=${PROJECT_ID}"
echo "  REGION=${REGION}"
echo "  API_SERVICE=${API_SERVICE}"

PROJECT_ID="${PROJECT_ID}" REGION="${REGION}" API_SERVICE="${API_SERVICE}" \
  "${SCRIPT_DIR}/deploy-prod.sh"
