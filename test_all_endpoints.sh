#!/bin/bash

BASE="http://localhost:3001/api"
PASS=0
FAIL=0
TOTAL=0
ADMIN_TOKEN=""
USER_TOKEN=""
USER_REFRESH=""
PACKAGE_ID=""
ORDER_ID=""
VIDEO_ID=""
VIDEO_ORDER_ID=""
SUPPORT_ID=""
CONTACT_ID=""
TESTIMONIAL_ID=""
SESSION_ID=""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

log() {
    local label="$1" status="$2" detail="$3"
    TOTAL=$((TOTAL+1))
    if [ "$status" = "PASS" ]; then
        PASS=$((PASS+1))
        echo -e "  ${GREEN}✓ PASS${NC} $label ${CYAN}$detail${NC}"
    else
        FAIL=$((FAIL+1))
        echo -e "  ${RED}✗ FAIL${NC} $label ${RED}$detail${NC}"
    fi
}

section() {
    echo ""
    echo -e "${BOLD}${YELLOW}═══════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${YELLOW}  $1${NC}"
    echo -e "${BOLD}${YELLOW}═══════════════════════════════════════════════════${NC}"
}

extract_json() {
    echo "$1" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('$2',''))" 2>/dev/null
}

extract_data_json() {
    echo "$1" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('$2',''))" 2>/dev/null
}

echo -e "${BOLD}${CYAN}╔═══════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${CYAN}║   KPM ACADEMY API - COMPLETE ENDPOINT TEST       ║${NC}"
echo -e "${BOLD}${CYAN}║   $(date '+%Y-%m-%d %H:%M:%S')                          ║${NC}"
echo -e "${BOLD}${CYAN}╚═══════════════════════════════════════════════════╝${NC}"

# Kill any existing server on port 3001
fuser -k 3001/tcp 2>/dev/null
sleep 1

# Start server on port 3001
echo -e "\n${BOLD}Starting server on port 3001...${NC}"
RATE_LIMIT=500 RATE_LIMIT_EXPIRATION=60 APP_PORT=3001 go run . &>/tmp/server.log &
SERVER_PID=$!
sleep 4

# Verify server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}Server failed to start!${NC}"
    cat /tmp/server.log
    exit 1
fi
echo -e "${GREEN}Server started (PID: $SERVER_PID)${NC}"

# Verify DB was re-created with admin
ADMIN_EMAIL=$(mysql -u root -pFakhrizal2005 -N -e "SELECT email FROM kpm_academy.users WHERE role='admin' LIMIT 1" 2>/dev/null)
echo -e "Admin in DB: ${GREEN}$ADMIN_EMAIL${NC}"

cleanup() {
    echo -e "\n${BOLD}Shutting down server...${NC}"
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    echo -e "${GREEN}Server stopped.${NC}"
}
trap cleanup EXIT

# ════════════════════════════════════════════════════════
# 1. HEALTH CHECK
# ════════════════════════════════════════════════════════
section "1. HEALTH CHECK"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/health")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /health" "PASS" "[HTTP $HTTP]"
else
    log "GET /health" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 2. AUTH - REGISTER
# ════════════════════════════════════════════════════════
section "2. AUTH - REGISTER"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test User","email":"testuser@test.com","password":"Test1234!","phone":"081234567890"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "201" ]; then
    log "POST /auth/register [create user]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/register [create user]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Register duplicate email
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Dup User","email":"testuser@test.com","password":"Test1234!","phone":"081234567890"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "409" ]; then
    log "POST /auth/register [duplicate email]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/register [duplicate email]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Register with missing fields
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"email":"x@x.com"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /auth/register [missing fields validation]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/register [missing fields validation]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 3. AUTH - LOGIN
# ════════════════════════════════════════════════════════
section "3. AUTH - LOGIN"

# Verify user first (auto-verify for testing)
mysql -u root -pFakhrizal2005 -e "UPDATE kpm_academy.users SET is_verified=1 WHERE email='testuser@test.com'" 2>/dev/null

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"testuser@test.com","password":"Test1234!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    USER_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['access_token'])" 2>/dev/null)
    USER_REFRESH=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['refresh_token'])" 2>/dev/null)
    log "POST /auth/login [user login]" "PASS" "[HTTP $HTTP] token_ok=$([ -n "$USER_TOKEN" ] && echo 'yes' || echo 'no')"
else
    log "POST /auth/login [user login]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Admin login
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@kpmacademy.com","password":"Admin123!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    ADMIN_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['access_token'])" 2>/dev/null)
    log "POST /auth/login [admin login]" "PASS" "[HTTP $HTTP] token_ok=$([ -n "$ADMIN_TOKEN" ] && echo 'yes' || echo 'no')"
else
    log "POST /auth/login [admin login]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Wrong password
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"testuser@test.com","password":"wrongpass"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "POST /auth/login [wrong password]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/login [wrong password]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Non-existent email
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"nobody@test.com","password":"Test1234!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "POST /auth/login [non-existent email]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/login [non-existent email]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 4. AUTH - FORGOT/RESET PASSWORD
# ════════════════════════════════════════════════════════
section "4. AUTH - FORGOT/RESET PASSWORD"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/forgot-password" \
    -H "Content-Type: application/json" \
    -d '{"email":"testuser@test.com"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
RESET_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('reset_token',''))" 2>/dev/null)
if [ "$HTTP" = "200" ] && [ -n "$RESET_TOKEN" ]; then
    log "POST /auth/forgot-password [success]" "PASS" "[HTTP $HTTP] got_token=yes"
else
    log "POST /auth/forgot-password [success]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Reset with token
if [ -n "$RESET_TOKEN" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/reset-password" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"testuser@test.com\",\"token\":\"$RESET_TOKEN\",\"new_password\":\"NewPass1234!\"}")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /auth/reset-password [success]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /auth/reset-password [success]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Login with new password
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"testuser@test.com","password":"NewPass1234!"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        USER_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['access_token'])" 2>/dev/null)
        USER_REFRESH=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['refresh_token'])" 2>/dev/null)
        log "POST /auth/login [with new password]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /auth/login [with new password]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# Reset with invalid token
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/reset-password" \
    -H "Content-Type: application/json" \
    -d '{"email":"testuser@test.com","token":"invalidtoken","new_password":"X1234567!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /auth/reset-password [invalid token]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/reset-password [invalid token]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Forgot password with non-existent email (should still return 200)
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/forgot-password" \
    -H "Content-Type: application/json" \
    -d '{"email":"ghost@test.com"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "POST /auth/forgot-password [non-existent email - no leak]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/forgot-password [non-existent email - no leak]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 5. AUTH - REFRESH TOKEN
# ════════════════════════════════════════════════════════
section "5. AUTH - REFRESH TOKEN"

if [ -n "$USER_REFRESH" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/refresh-token" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$USER_REFRESH\"}")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        NEW_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['access_token'])" 2>/dev/null)
        if [ -n "$NEW_TOKEN" ]; then
            USER_TOKEN="$NEW_TOKEN"
            log "POST /auth/refresh-token" "PASS" "[HTTP $HTTP] new_token=yes"
        else
            log "POST /auth/refresh-token" "FAIL" "[HTTP $HTTP] no new token in response"
        fi
    else
        log "POST /auth/refresh-token" "FAIL" "[HTTP $HTTP] $BODY"
    fi
else
    log "POST /auth/refresh-token" "FAIL" "no refresh token available"
fi

# Refresh with invalid token
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/refresh-token" \
    -H "Content-Type: application/json" \
    -d '{"refresh_token":"invalidtoken"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "POST /auth/refresh-token [invalid token]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/refresh-token [invalid token]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 6. PROFILE
# ════════════════════════════════════════════════════════
section "6. PROFILE"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/profile" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /profile" "PASS" "[HTTP $HTTP]"
else
    log "GET /profile" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Update profile
RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/profile" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated User","phone":"08999999999","student_name":"Ahmad","student_class":"XII","student_major":"IPA","school_name":"SMAN 1","gender":"male","religion":"islam"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "PUT /profile [update]" "PASS" "[HTTP $HTTP]"
else
    log "PUT /profile [update]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Profile without token
RESP=$(curl -s -w "\n%{http_code}" "$BASE/profile")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "GET /profile [no token]" "PASS" "[HTTP $HTTP]"
else
    log "GET /profile [no token]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Change password
RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/profile/change-password" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"old_password":"NewPass1234!","new_password":"FinalPass123!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "PUT /profile/change-password [success]" "PASS" "[HTTP $HTTP]"
    # Re-login with new password
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"testuser@test.com","password":"FinalPass123!"}')
    BODY=$(echo "$RESP" | sed '$d')
    USER_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['access_token'])" 2>/dev/null)
    USER_REFRESH=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['tokens']['refresh_token'])" 2>/dev/null)
else
    log "PUT /profile/change-password [success]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Change password with wrong old password
RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/profile/change-password" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"old_password":"wrongold","new_password":"X1234567!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "PUT /profile/change-password [wrong old]" "PASS" "[HTTP $HTTP]"
else
    log "PUT /profile/change-password [wrong old]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 7. PUBLIC ENDPOINTS
# ════════════════════════════════════════════════════════
section "7. PUBLIC ENDPOINTS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/packages")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /packages [public]" "PASS" "[HTTP $HTTP]"
else
    log "GET /packages [public]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/videos")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /videos [public]" "PASS" "[HTTP $HTTP]"
else
    log "GET /videos [public]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/testimonials")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /testimonials [public]" "PASS" "[HTTP $HTTP]"
else
    log "GET /testimonials [public]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Non-existent package
RESP=$(curl -s -w "\n%{http_code}" "$BASE/packages/00000000-0000-0000-0000-000000000000")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "404" ]; then
    log "GET /packages/:id [not found]" "PASS" "[HTTP $HTTP]"
else
    log "GET /packages/:id [not found]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Invalid UUID
RESP=$(curl -s -w "\n%{http_code}" "$BASE/packages/not-a-uuid")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "GET /packages/:id [invalid UUID]" "PASS" "[HTTP $HTTP]"
else
    log "GET /packages/:id [invalid UUID]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Non-existent video
RESP=$(curl -s -w "\n%{http_code}" "$BASE/videos/00000000-0000-0000-0000-000000000000")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "404" ]; then
    log "GET /videos/:id [not found]" "PASS" "[HTTP $HTTP]"
else
    log "GET /videos/:id [not found]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 8. CONTACT FORM
# ════════════════════════════════════════════════════════
section "8. CONTACT FORM"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/kontak-form" \
    -H "Content-Type: application/json" \
    -d '{"name":"Budi","email":"budi@test.com","phone":"081111111111","subject":"Test","message":"Hello from test"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "201" ]; then
    log "POST /kontak-form [create]" "PASS" "[HTTP $HTTP]"
else
    log "POST /kontak-form [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Missing required fields
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/kontak-form" \
    -H "Content-Type: application/json" \
    -d '{"name":"NoEmail"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /kontak-form [missing fields]" "PASS" "[HTTP $HTTP]"
else
    log "POST /kontak-form [missing fields]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 9. SUPPORT
# ════════════════════════════════════════════════════════
section "9. SUPPORT"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/support/submit" \
    -H "Content-Type: application/json" \
    -d '{"name":"Help User","email":"help@test.com","question":"How to reset password?"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "201" ]; then
    log "POST /support/submit [create]" "PASS" "[HTTP $HTTP]"
else
    log "POST /support/submit [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Missing required fields
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/support/submit" \
    -H "Content-Type: application/json" \
    -d '{"name":"NoQ"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /support/submit [missing fields]" "PASS" "[HTTP $HTTP]"
else
    log "POST /support/submit [missing fields]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 10. CHAT
# ════════════════════════════════════════════════════════
section "10. CHAT"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/chat/send" \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello support!"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
CHAT_SESSION=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('session_id',''))" 2>/dev/null)
if [ "$HTTP" = "201" ]; then
    log "POST /chat/send [create]" "PASS" "[HTTP $HTTP]"
else
    log "POST /chat/send [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$CHAT_SESSION" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/chat/history?session_id=$CHAT_SESSION")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /chat/history [with session_id]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /chat/history [with session_id]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# Chat without session_id
RESP=$(curl -s -w "\n%{http_code}" "$BASE/chat/history")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "GET /chat/history [no session_id]" "PASS" "[HTTP $HTTP]"
else
    log "GET /chat/history [no session_id]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 11. USER DASHBOARD
# ════════════════════════════════════════════════════════
section "11. USER DASHBOARD"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/dashboard" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /dashboard [user]" "PASS" "[HTTP $HTTP]"
else
    log "GET /dashboard [user]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 12. ADMIN - PACKAGE CRUD
# ════════════════════════════════════════════════════════
section "12. ADMIN - PACKAGE CRUD"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/packages" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"Paket UTBK Saintek","description":"Paket lengkap UTBK Saintek","kelas":"XII","jenjang":"SMA","price":150000,"membership_duration_days":30,"cards":"[{\"id\":\"card-1\",\"title\":\"Biologi\"}]","questions":"[{\"question\":\"Apa itu sel?\",\"options\":[\"A\",\"B\",\"C\",\"D\"],\"answer\":\"A\"}]"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
PACKAGE_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
if [ "$HTTP" = "201" ] && [ -n "$PACKAGE_ID" ]; then
    log "POST /admin/packages [create]" "PASS" "[HTTP $HTTP] id=$PACKAGE_ID"
else
    log "POST /admin/packages [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Missing title
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/packages" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"description":"No title"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /admin/packages [missing title]" "PASS" "[HTTP $HTTP]"
else
    log "POST /admin/packages [missing title]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/packages" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/packages [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/packages [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$PACKAGE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/packages/$PACKAGE_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/packages/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/packages/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/packages/$PACKAGE_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"title":"Paket UTBK Saintek Updated","description":"Updated","price":200000}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "PUT /admin/packages/:id [update]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/packages/:id [update]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Add card
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/packages/$PACKAGE_ID/cards" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"cards":"[{\"id\":\"card-1\",\"title\":\"Biologi\"},{\"id\":\"card-2\",\"title\":\"Kimia\"}]"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/packages/:id/cards [add]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/packages/:id/cards [add]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Remove card
    RESP=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE/admin/packages/$PACKAGE_ID/cards/card-1" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "DELETE /admin/packages/:id/cards/:cardId [remove]" "PASS" "[HTTP $HTTP]"
    else
        log "DELETE /admin/packages/:id/cards/:cardId [remove]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Import questions
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/packages/$PACKAGE_ID/import-pdf" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"questions":"[{\"question\":\"Soal baru\",\"options\":[\"A\",\"B\"],\"answer\":\"A\"}]"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/packages/:id/import-pdf [import]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/packages/:id/import-pdf [import]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# User tries admin endpoint
RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/packages" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "403" ]; then
    log "GET /admin/packages [user forbidden]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/packages [user forbidden]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 13. ADMIN - VIDEO CRUD
# ════════════════════════════════════════════════════════
section "13. ADMIN - VIDEO CRUD"

VID_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/videos" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"Video Biologi Sel","description":"Materi tentang sel","video_url":"https://youtube.com/watch?v=abc123","thumbnail":"thumb.jpg","price":50000,"access_duration_days":30}')
HTTP=$(echo "$VID_RESP" | tail -1)
BODY=$(echo "$VID_RESP" | sed '$d')
VIDEO_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
if [ "$HTTP" = "201" ] && [ -n "$VIDEO_ID" ]; then
    log "POST /admin/videos [create]" "PASS" "[HTTP $HTTP] id=$VIDEO_ID"
else
    log "POST /admin/videos [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/videos" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/videos [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/videos [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$VIDEO_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/videos/$VIDEO_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/videos/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/videos/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/videos/$VIDEO_ID/toggle" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/videos/:id/toggle [toggle off]" "PASS" "[HTTP $HTTP]"
        # Toggle back on for later tests
        curl -s -X POST "$BASE/admin/videos/$VIDEO_ID/toggle" -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
    else
        log "POST /admin/videos/:id/toggle [toggle off]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/videos/$VIDEO_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"title":"Video Updated","description":"Updated desc","price":75000}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "PUT /admin/videos/:id [update]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/videos/:id [update]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/videos/$VIDEO_ID/toggle" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/videos/:id/toggle [toggle back on]" "PASS" "[HTTP $HTTP]"
        # Toggle one more time to ensure active for later tests
        curl -s -X POST "$BASE/admin/videos/$VIDEO_ID/toggle" -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
    else
        log "POST /admin/videos/:id/toggle [toggle back on]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 14. ORDERS (User)
# ════════════════════════════════════════════════════════
section "14. ORDERS - USER"

if [ -n "$PACKAGE_ID" ]; then
    ORD_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/orders" \
        -H "Authorization: Bearer $USER_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"package_id\":\"$PACKAGE_ID\"}")
    HTTP=$(echo "$ORD_RESP" | tail -1)
    BODY=$(echo "$ORD_RESP" | sed '$d')
    ORDER_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
    if [ "$HTTP" = "201" ] && [ -n "$ORDER_ID" ]; then
        log "POST /orders [create order]" "PASS" "[HTTP $HTTP] id=$ORDER_ID"
    else
        log "POST /orders [create order]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# Order without package_id
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/orders" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /orders [missing package_id]" "PASS" "[HTTP $HTTP]"
else
    log "POST /orders [missing package_id]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Order with invalid package_id
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/orders" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"package_id":"not-a-uuid"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /orders [invalid package_id]" "PASS" "[HTTP $HTTP]"
else
    log "POST /orders [invalid package_id]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/orders" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /orders [list user orders]" "PASS" "[HTTP $HTTP]"
else
    log "GET /orders [list user orders]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$ORDER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/orders/$ORDER_ID" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /orders/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /orders/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" "$BASE/orders/status" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /orders/status [payment status]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /orders/status [payment status]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 15. PAYMENT SIMULATION
# ════════════════════════════════════════════════════════
section "15. PAYMENT SIMULATION"

if [ -n "$ORDER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/orders/$ORDER_ID/pay" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /orders/:id/pay [simulate payment]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /orders/:id/pay [simulate payment]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Pay again (should fail - already paid)
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/orders/$ORDER_ID/pay" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "400" ]; then
        log "POST /orders/:id/pay [already paid]" "PASS" "[HTTP $HTTP]"
    else
        log "POST /orders/:id/pay [already paid]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 16. PRACTICE SESSION
# ════════════════════════════════════════════════════════
section "16. PRACTICE SESSION"

if [ -n "$PACKAGE_ID" ] && [ -n "$ORDER_ID" ]; then
    PRACT_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/practice/start" \
        -H "Authorization: Bearer $USER_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"package_id\":\"$PACKAGE_ID\",\"card_id\":\"card-1\"}")
    HTTP=$(echo "$PRACT_RESP" | tail -1)
    BODY=$(echo "$PRACT_RESP" | sed '$d')
    SESSION_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
    if [ "$HTTP" = "201" ] && [ -n "$SESSION_ID" ]; then
        log "POST /practice/start [start session]" "PASS" "[HTTP $HTTP] id=$SESSION_ID"
    else
        log "POST /practice/start [start session]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    if [ -n "$SESSION_ID" ]; then
        RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/practice/submit" \
            -H "Authorization: Bearer $USER_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"session_id\":\"$SESSION_ID\",\"answers\":\"[{\\\"question_id\\\":1,\\\"answer\\\":\\\"A\\\"}]\",\"duration\":120}")
        HTTP=$(echo "$RESP" | tail -1)
        BODY=$(echo "$RESP" | sed '$d')
        if [ "$HTTP" = "200" ]; then
            log "POST /practice/submit [submit answers]" "PASS" "[HTTP $HTTP]"
        else
            log "POST /practice/submit [submit answers]" "FAIL" "[HTTP $HTTP] $BODY"
        fi

        # Submit again (should fail - already finished)
        RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/practice/submit" \
            -H "Authorization: Bearer $USER_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"session_id\":\"$SESSION_ID\",\"answers\":\"[]\",\"duration\":60}")
        HTTP=$(echo "$RESP" | tail -1)
        BODY=$(echo "$RESP" | sed '$d')
        if [ "$HTTP" = "400" ]; then
            log "POST /practice/submit [already finished]" "PASS" "[HTTP $HTTP]"
        else
            log "POST /practice/submit [already finished]" "FAIL" "[HTTP $HTTP] $BODY"
        fi

        RESP=$(curl -s -w "\n%{http_code}" "$BASE/practice/$SESSION_ID" \
            -H "Authorization: Bearer $USER_TOKEN")
        HTTP=$(echo "$RESP" | tail -1)
        BODY=$(echo "$RESP" | sed '$d')
        if [ "$HTTP" = "200" ]; then
            log "GET /practice/:id [show session]" "PASS" "[HTTP $HTTP]"
        else
            log "GET /practice/:id [show session]" "FAIL" "[HTTP $HTTP] $BODY"
        fi
    fi

    RESP=$(curl -s -w "\n%{http_code}" "$BASE/practice/history" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /practice/history" "PASS" "[HTTP $HTTP]"
    else
        log "GET /practice/history" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" "$BASE/practice/statistics" \
        -H "Authorization: Bearer $USER_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /practice/statistics [user stats]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /practice/statistics [user stats]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# Start without purchased package (use non-existent package UUID)
FAKE_PKG="00000000-0000-0000-0000-000000000000"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/practice/start" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"package_id\":\"$FAKE_PKG\",\"card_id\":\"card-99\"}")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "404" ]; then
    log "POST /practice/start [no purchase - not found]" "PASS" "[HTTP $HTTP]"
elif [ "$HTTP" = "403" ]; then
    log "POST /practice/start [no purchase - forbidden]" "PASS" "[HTTP $HTTP]"
else
    log "POST /practice/start [no purchase - forbidden]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 17. VIDEO ORDER + PAYMENT
# ════════════════════════════════════════════════════════
section "17. VIDEO ORDER + PAYMENT"

if [ -n "$VIDEO_ID" ]; then
    VORD_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/videos/$VIDEO_ID/order" \
        -H "Authorization: Bearer $USER_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"total_price":50000}')
    HTTP=$(echo "$VORD_RESP" | tail -1)
    BODY=$(echo "$VORD_RESP" | sed '$d')
    VIDEO_ORDER_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
    if [ "$HTTP" = "201" ] && [ -n "$VIDEO_ORDER_ID" ]; then
        log "POST /videos/:id/order [create video order]" "PASS" "[HTTP $HTTP] id=$VIDEO_ORDER_ID"
    else
        log "POST /videos/:id/order [create video order]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    if [ -n "$VIDEO_ORDER_ID" ]; then
        RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/videos/$VIDEO_ID/pay/$VIDEO_ORDER_ID" \
            -H "Authorization: Bearer $USER_TOKEN")
        HTTP=$(echo "$RESP" | tail -1)
        BODY=$(echo "$RESP" | sed '$d')
        if [ "$HTTP" = "200" ]; then
            log "POST /videos/:id/pay/:orderId [simulate video payment]" "PASS" "[HTTP $HTTP]"
        else
            log "POST /videos/:id/pay/:orderId [simulate video payment]" "FAIL" "[HTTP $HTTP] $BODY"
        fi
    fi
fi

# ════════════════════════════════════════════════════════
# 18. TESTIMONIALS
# ════════════════════════════════════════════════════════
section "18. TESTIMONIALS"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/testimonials" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"Sangat membantu!","rating":5}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
TESTIMONIAL_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null)
if [ "$HTTP" = "201" ] && [ -n "$TESTIMONIAL_ID" ]; then
    log "POST /testimonials [create]" "PASS" "[HTTP $HTTP] id=$TESTIMONIAL_ID"
else
    log "POST /testimonials [create]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Empty content
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/testimonials" \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"","rating":5}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /testimonials [empty content]" "PASS" "[HTTP $HTTP]"
else
    log "POST /testimonials [empty content]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/testimonials/my" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /testimonials/my" "PASS" "[HTTP $HTTP]"
else
    log "GET /testimonials/my" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 19. NOTIFICATIONS
# ════════════════════════════════════════════════════════
section "19. NOTIFICATIONS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/notifications" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /notifications [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /notifications [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/notifications/unread-count" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /notifications/unread-count" "PASS" "[HTTP $HTTP]"
else
    log "GET /notifications/unread-count" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/notifications/read-all" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "POST /notifications/read-all" "PASS" "[HTTP $HTTP]"
else
    log "POST /notifications/read-all" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 20. USER KONTAK-FORM (Protected)
# ════════════════════════════════════════════════════════
section "20. USER KONTAK-FORM (Protected)"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/kontak-form" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /kontak-form [user list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /kontak-form [user list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 21. ADMIN DASHBOARD
# ════════════════════════════════════════════════════════
section "21. ADMIN DASHBOARD"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/dashboard" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/dashboard" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/dashboard" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 22. ADMIN - USER MANAGEMENT
# ════════════════════════════════════════════════════════
section "22. ADMIN - USER MANAGEMENT"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/users" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/users [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/users [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Get test user ID
TEST_USER_ID=$(mysql -u root -pFakhrizal2005 -N -e "SELECT id FROM kpm_academy.users WHERE email='testuser@test.com' LIMIT 1" 2>/dev/null)
if [ -n "$TEST_USER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/users/$TEST_USER_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/users/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/users/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/users/$TEST_USER_ID/toggle-active" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/users/:id/toggle-active" "PASS" "[HTTP $HTTP]"
        # Toggle back
        curl -s -X POST "$BASE/admin/users/$TEST_USER_ID/toggle-active" -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
    else
        log "POST /admin/users/:id/toggle-active" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# Non-existent user
RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/users/00000000-0000-0000-0000-000000000000" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "404" ]; then
    log "GET /admin/users/:id [not found]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/users/:id [not found]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/login-logs" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/login-logs" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/login-logs" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 23. ADMIN - ORDER MANAGEMENT
# ════════════════════════════════════════════════════════
section "23. ADMIN - ORDER MANAGEMENT"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/orders" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/orders [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/orders [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$ORDER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/orders/$ORDER_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/orders/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/orders/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/orders/$ORDER_ID/verify" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/orders/:id/verify" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/orders/:id/verify" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 24. ADMIN - TRANSACTIONS
# ════════════════════════════════════════════════════════
section "24. ADMIN - TRANSACTIONS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/transactions" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/transactions [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/transactions [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/transactions/stats" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/transactions/stats" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/transactions/stats" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$ORDER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/transactions/$ORDER_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/transactions/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/transactions/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 25. ADMIN - ENROLL KEYS
# ════════════════════════════════════════════════════════
section "25. ADMIN - ENROLL KEYS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/enroll-keys" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/enroll-keys [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/enroll-keys [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 26. ADMIN - PRACTICE STATISTICS
# ════════════════════════════════════════════════════════
section "26. ADMIN - PRACTICE STATISTICS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/practice-statistics" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/practice-statistics [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/practice-statistics [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$SESSION_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/practice-statistics/$SESSION_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/practice-statistics/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/practice-statistics/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 27. ADMIN - REPORTS
# ════════════════════════════════════════════════════════
section "27. ADMIN - REPORTS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/reports" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/reports" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/reports" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 28. ADMIN - TESTIMONIAL MANAGEMENT
# ════════════════════════════════════════════════════════
section "28. ADMIN - TESTIMONIAL MANAGEMENT"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/testimonials" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/testimonials [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/testimonials [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$TESTIMONIAL_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/testimonials/$TESTIMONIAL_ID/approve" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/testimonials/:id/approve" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/testimonials/:id/approve" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/testimonials/$TESTIMONIAL_ID/toggle-active" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/testimonials/:id/toggle-active" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/testimonials/:id/toggle-active" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 29. ADMIN - SUPPORT MANAGEMENT
# ════════════════════════════════════════════════════════
section "29. ADMIN - SUPPORT MANAGEMENT"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/support" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/support [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/support [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Get support ticket ID from DB
SUPPORT_ID=$(mysql -u root -pFakhrizal2005 -N -e "SELECT id FROM kpm_academy.support_tickets LIMIT 1" 2>/dev/null)
if [ -n "$SUPPORT_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/support/$SUPPORT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/support/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/support/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/support/$SUPPORT_ID/answer" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"answer":"You can reset your password from the profile page."}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/support/:id/answer" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/support/:id/answer" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/support/$SUPPORT_ID/status" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"status":"answered"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "PUT /admin/support/:id/status [update status]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/support/:id/status [update status]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/support/$SUPPORT_ID/status" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"status":"invalid_status"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "400" ]; then
        log "PUT /admin/support/:id/status [invalid status]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/support/:id/status [invalid status]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 30. ADMIN - CONTACT FORM MANAGEMENT
# ════════════════════════════════════════════════════════
section "30. ADMIN - CONTACT FORM MANAGEMENT"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/kontak-form" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/kontak-form [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/kontak-form [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

CONTACT_ID=$(mysql -u root -pFakhrizal2005 -N -e "SELECT id FROM kpm_academy.contact_forms LIMIT 1" 2>/dev/null)
if [ -n "$CONTACT_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/kontak-form/$CONTACT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "GET /admin/kontak-form/:id [show]" "PASS" "[HTTP $HTTP]"
    else
        log "GET /admin/kontak-form/:id [show]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/kontak-form/$CONTACT_ID/reply" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"admin_reply":"We will get back to you soon."}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/kontak-form/:id/reply" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/kontak-form/:id/reply" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/kontak-form/$CONTACT_ID/status" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"status":"replied"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "PUT /admin/kontak-form/:id/status [update status]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/kontak-form/:id/status [update status]" "FAIL" "[HTTP $HTTP] $BODY"
    fi

    # Invalid status
    RESP=$(curl -s -w "\n%{http_code}" -X PUT "$BASE/admin/kontak-form/$CONTACT_ID/status" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"status":"bogus"}')
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "400" ]; then
        log "PUT /admin/kontak-form/:id/status [invalid status]" "PASS" "[HTTP $HTTP]"
    else
        log "PUT /admin/kontak-form/:id/status [invalid status]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 31. ADMIN - VIDEO ORDERS
# ════════════════════════════════════════════════════════
section "31. ADMIN - VIDEO ORDERS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/video-orders" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/video-orders [list]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/video-orders [list]" "FAIL" "[HTTP $HTTP] $BODY"
fi

if [ -n "$VIDEO_ORDER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/admin/video-orders/$VIDEO_ORDER_ID/grant" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "POST /admin/video-orders/:id/grant" "PASS" "[HTTP $HTTP]"
    else
        log "POST /admin/video-orders/:id/grant" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 32. ADMIN - NOTIFICATIONS
# ════════════════════════════════════════════════════════
section "32. ADMIN - NOTIFICATIONS"

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/notifications" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "GET /admin/notifications" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/notifications" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 33. PAYMENT NOTIFICATION WEBHOOK
# ════════════════════════════════════════════════════════
section "33. PAYMENT NOTIFICATION WEBHOOK"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/payment/notification" \
    -H "Content-Type: application/json" \
    -d '{"status":"settlement","order_id":"test"}')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "POST /payment/notification [webhook]" "PASS" "[HTTP $HTTP]"
else
    log "POST /payment/notification [webhook]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 34. AUTH - LOGOUT
# ════════════════════════════════════════════════════════
section "34. AUTH - LOGOUT"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/logout" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "200" ]; then
    log "POST /auth/logout [success]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/logout [success]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# 35. CLEANUP - DELETE TEST DATA
# ════════════════════════════════════════════════════════
section "35. CLEANUP - DELETE TEST DATA"

if [ -n "$TESTIMONIAL_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE/admin/testimonials/$TESTIMONIAL_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "DELETE /admin/testimonials/:id [cleanup]" "PASS" "[HTTP $HTTP]"
    else
        log "DELETE /admin/testimonials/:id [cleanup]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

if [ -n "$VIDEO_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE/admin/videos/$VIDEO_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "DELETE /admin/videos/:id [cleanup]" "PASS" "[HTTP $HTTP]"
    else
        log "DELETE /admin/videos/:id [cleanup]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

if [ -n "$PACKAGE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE/admin/packages/$PACKAGE_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    HTTP=$(echo "$RESP" | tail -1)
    BODY=$(echo "$RESP" | sed '$d')
    if [ "$HTTP" = "200" ]; then
        log "DELETE /admin/packages/:id [cleanup]" "PASS" "[HTTP $HTTP]"
    else
        log "DELETE /admin/packages/:id [cleanup]" "FAIL" "[HTTP $HTTP] $BODY"
    fi
fi

# ════════════════════════════════════════════════════════
# 36. EDGE CASES & ERROR HANDLING
# ════════════════════════════════════════════════════════
section "36. EDGE CASES & ERROR HANDLING"

# Request with invalid JSON body
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d 'not json')
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "400" ]; then
    log "POST /auth/login [invalid JSON]" "PASS" "[HTTP $HTTP]"
else
    log "POST /auth/login [invalid JSON]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Access protected endpoint with invalid token
RESP=$(curl -s -w "\n%{http_code}" "$BASE/profile" \
    -H "Authorization: Bearer invalidtoken123")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "GET /profile [invalid token]" "PASS" "[HTTP $HTTP]"
else
    log "GET /profile [invalid token]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Access admin endpoint without admin role
RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/dashboard" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "403" ]; then
    log "GET /admin/dashboard [user not admin]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/dashboard [user not admin]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Access protected endpoint without auth header
RESP=$(curl -s -w "\n%{http_code}" "$BASE/orders")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "401" ]; then
    log "GET /orders [no auth header]" "PASS" "[HTTP $HTTP]"
else
    log "GET /orders [no auth header]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Non-existent route
RESP=$(curl -s -w "\n%{http_code}" "$BASE/nonexistent/route")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "404" ]; then
    log "GET /nonexistent/route [404]" "PASS" "[HTTP $HTTP]"
else
    log "GET /nonexistent/route [404]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# Admin endpoints with user token
RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/orders" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "403" ]; then
    log "GET /admin/orders [user forbidden]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/orders [user forbidden]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/support" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "403" ]; then
    log "GET /admin/support [user forbidden]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/support [user forbidden]" "FAIL" "[HTTP $HTTP] $BODY"
fi

RESP=$(curl -s -w "\n%{http_code}" "$BASE/admin/kontak-form" \
    -H "Authorization: Bearer $USER_TOKEN")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
if [ "$HTTP" = "403" ]; then
    log "GET /admin/kontak-form [user forbidden]" "PASS" "[HTTP $HTTP]"
else
    log "GET /admin/kontak-form [user forbidden]" "FAIL" "[HTTP $HTTP] $BODY"
fi

# ════════════════════════════════════════════════════════
# RESULTS
# ════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}${CYAN}╔═══════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${CYAN}║                   TEST RESULTS                   ║${NC}"
echo -e "${BOLD}${CYAN}╚═══════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  Total:  ${BOLD}$TOTAL${NC}"
echo -e "  Passed: ${GREEN}${BOLD}$PASS${NC}"
echo -e "  Failed: ${RED}${BOLD}$FAIL${NC}"
echo ""

if [ "$FAIL" -eq 0 ]; then
    echo -e "  ${GREEN}${BOLD}ALL TESTS PASSED ✓${NC}"
else
    echo -e "  ${RED}${BOLD}$FAIL TEST(S) FAILED ✗${NC}"
fi
echo ""
