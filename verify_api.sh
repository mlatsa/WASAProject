#!/usr/bin/env bash
set -euo pipefail
BASE="http://localhost:3000"

ID=$(curl -s -X POST "$BASE/session" -H 'Content-Type: application/json' -d '{"name":"Alex"}' | sed -n 's/.*"identifier":"\([^"]*\)".*/\1/p')
echo "ID=$ID"
test -n "$ID"

AUTH=(-H "Authorization: Bearer $ID")

call(){ curl -s -o /dev/null -w "%{http_code}  $1 $2\n" -X "$1" "$2" "${@:3}" || true; }

call GET  "$BASE/health"
call POST "$BASE/session" -H 'Content-Type: application/json' -d '{"name":"Bob"}'
call PUT  "$BASE/user/username" "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"username":"bob_01"}'
call PUT  "$BASE/user/photo" "${AUTH[@]}"
call GET  "$BASE/conversations" "${AUTH[@]}"
call GET  "$BASE/conversations/c1" "${AUTH[@]}"
MSG=$(curl -s -X POST "$BASE/conversations/c1/messages" "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"content":"hi","type":"text"}')
MID=$(echo "$MSG" | sed -n 's/.*"messageId":"\([^"]*\)".*/\1/p'); echo "MID=$MID"
call POST "$BASE/messages/$MID/forward"   "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"conversationId":"c2"}'
call POST "$BASE/messages/$MID/reactions" "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"emoji":"👍"}'
call DELETE "$BASE/messages/$MID/reactions/any" "${AUTH[@]}"
call DELETE "$BASE/messages/$MID" "${AUTH[@]}"
call POST "$BASE/groups/c1/members" "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"member":"Charlie"}'
call POST "$BASE/groups/c1/leave"   "${AUTH[@]}"
call PUT  "$BASE/groups/c1/name"    "${AUTH[@]}" -H 'Content-Type: application/json' -d '{"name":"Chat Group"}'
call PUT  "$BASE/groups/c1/photo"   "${AUTH[@]}"
