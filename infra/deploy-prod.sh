#!/usr/bin/env bash
set -euo pipefail

PROJECT_ID="${PROJECT_ID:-koala-map-jp}"
REGION="${REGION:-asia-northeast1}"
API_SERVICE="${API_SERVICE:-koala-map-jp-api}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "[1/3] Deploy Cloud Run service: ${API_SERVICE}"
gcloud run deploy "${API_SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --source "${REPO_ROOT}/api" \
  --allow-unauthenticated \
  --set-env-vars "FIREBASE_PROJECT_ID=${PROJECT_ID},SERVE_STATIC=false"

echo "[2/3] Deploy Firestore rules/indexes"
firebase deploy --project "${PROJECT_ID}" --only firestore

echo "[3/3] Deploy Hosting"
firebase deploy --project "${PROJECT_ID}" --only hosting

echo "Done: https://${PROJECT_ID}.web.app"
