#!/usr/bin/env sh

set -e

TABLE_NAME="scheduler-service-table"
DYNAMO_ENDPOINT="http://dynamodb:8000"
REGION="us-east-1"

echo "[dynamodb-init] Waiting for DynamoDB Local..."

until aws dynamodb list-tables \
  --endpoint-url "$DYNAMO_ENDPOINT" \
  >/dev/null 2>&1; do
  sleep 2
done

echo "[dynamodb-init] DynamoDB is ready"

echo "[dynamodb-init] Creating table if not exists: $TABLE_NAME"

aws dynamodb create-table \
  --table-name "$TABLE_NAME" \
  --attribute-definitions \
    AttributeName=pk,AttributeType=S \
    AttributeName=sk,AttributeType=S \
  --key-schema \
    AttributeName=pk,KeyType=HASH \
    AttributeName=sk,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --region "$REGION" \
  --endpoint-url "$DYNAMO_ENDPOINT" \
  || echo "[dynamodb-init] Table already exists, skipping"

echo "[dynamodb-init] Done"