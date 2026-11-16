#!/usr/bin/env bash
#
# Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT
#
# Demonstrates the bugs fixed by fabrica commit
#   24238888a76741b4e883ecf2f1f0ff199d72eed0
#   "Fixed UpdatedAt timestamp for patch and timestamps for postgres"
# against user-service-postgres generated from fabrica main.
#
#   1. POST returns createdAt/updatedAt with nanosecond precision (and local
#      TZ); PostgreSQL stores microseconds, so a GET of the same resource
#      returns different timestamps.
#   2. PATCH returns an updatedAt that differs from the one a later GET returns.
#   3. PATCH that removes a spec field (merge-patch null / JSON Patch remove)
#      is silently ignored: the patched JSON is unmarshalled on top of the
#      existing spec, so the removed field keeps its previous value.
#
# Prereqs: postgres + the service running, e.g.
#   ./bin/user-service serve --port 8080 \
#     --database-url "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable"
#
# Usage: BASE=http://localhost:8080 ./demo-timestamp-patch-issue.sh

set -u
BASE="${BASE:-http://localhost:8080}"
CT_JSON='Content-Type: application/json'
CT_MERGE='Content-Type: application/merge-patch+json'
CT_JSONPATCH='Content-Type: application/json-patch+json'
SUFFIX="$(date +%s)"

compare() { # label, write-value, read-value
	if [[ "$2" == "$3" ]]; then
		echo "  OK       $1: write=$2 read=$3"
	else
		echo "  MISMATCH $1: write=$2 read=$3"
	fi
}

echo "== 1. POST timestamps vs GET timestamps =="
set -x
CREATE=$(curl -s -X POST "$BASE/users" -H "$CT_JSON" -d '{
  "metadata": { "name": "alice'"$SUFFIX"'" },
  "spec": {
    "username": "alice'"$SUFFIX"'",
    "email": "alice'"$SUFFIX"'@example.com",
    "password": "s3cretpass",
    "fullName": "Alice Example",
    "role": "admin",
    "active": true
  }
}')
{ set +x; } 2>/dev/null
RESOURCE_UID=$(jq -r '.metadata.uid' <<<"$CREATE")
set -x
GET1=$(curl -s "$BASE/users/$RESOURCE_UID")
{ set +x; } 2>/dev/null
compare createdAt "$(jq -r .metadata.createdAt <<<"$CREATE")" "$(jq -r .metadata.createdAt <<<"$GET1")"
compare updatedAt "$(jq -r .metadata.updatedAt <<<"$CREATE")" "$(jq -r .metadata.updatedAt <<<"$GET1")"

echo
echo "== 2. PATCH updatedAt vs GET updatedAt =="
set -x
PATCH1=$(curl -s -X PATCH "$BASE/users/$RESOURCE_UID" -H "$CT_MERGE" -d '{"fullName":"Alice Patched"}')
GET2=$(curl -s "$BASE/users/$RESOURCE_UID")
{ set +x; } 2>/dev/null
compare updatedAt "$(jq -r .metadata.updatedAt <<<"$PATCH1")" "$(jq -r .metadata.updatedAt <<<"$GET2")"

echo
echo "== 3a. merge-patch {\"active\": null} should clear active (true -> false) =="
set -x
curl -s -X PATCH "$BASE/users/$RESOURCE_UID" -H "$CT_MERGE" -d '{"active":null}' | jq -c '{active: .spec.active}'
curl -s "$BASE/users/$RESOURCE_UID" | jq -c '{active: .spec.active}'
{ set +x; } 2>/dev/null
echo "  expected: active=false (fixed)   buggy: active=true"

echo
echo "== 3b. JSON Patch remove /fullName (required) should fail validation (400) =="
set -x
curl -s -w '\nHTTP %{http_code}\n' -X PATCH "$BASE/users/$RESOURCE_UID" -H "$CT_JSONPATCH" \
	-d '[{"op":"remove","path":"/fullName"}]'
{ set +x; } 2>/dev/null
echo "  expected: HTTP 400 validation failed (fixed)   buggy: HTTP 200 with fullName unchanged"

echo
echo "== cleanup =="
curl -s -o /dev/null -w 'DELETE HTTP %{http_code}\n' -X DELETE "$BASE/users/$RESOURCE_UID"
