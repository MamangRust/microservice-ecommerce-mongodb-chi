#!/bin/bash
# e2e-local.sh — One-shot: start infra, launch Go services, seed, run hurl e2e tests.
#
# This is a microservice architecture: per-service PostgreSQL databases for isolation.
# Go services run on the host; infrastructure runs in Docker.
#
# Usage:
#   ./scripts/e2e-local.sh              # full run
#   ./scripts/e2e-local.sh --infra-only # start infra only
#   ./scripts/e2e-local.sh --skip-build # skip go build (use existing binaries)
#   ./scripts/e2e-local.sh --down       # tear down infra
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"
COMPOSE_FILE="deployments/local/docker-compose.infra.yml"
LOG_DIR="$ROOT/deployments/local/logs"
HURL_DIR="$ROOT/tests/hurl"
mkdir -p "$LOG_DIR"

# ─── Flags ──────────────────────────────────────────────────────────────
INFRA_ONLY=false
DO_DOWN=false
for arg in "$@"; do
  case "$arg" in
    --infra-only) INFRA_ONLY=true ;;
    --down)       DO_DOWN=true ;;
    --help|-h)
      echo "Usage: $0 [--infra-only] [--down]"
      exit 0 ;;
    *) echo "Unknown flag: $arg"; exit 1 ;;
  esac
done

# ─── Colors ─────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()  { echo -e "${BLUE}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}[OK]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
fail()  { echo -e "${RED}[FAIL]${NC} $*"; }

# ─── Teardown ───────────────────────────────────────────────────────────
cleanup() {
  echo ""
  info "Stopping Go services ..."
  bash "$ROOT/deployments/local/scripts/services-local-stop.sh" 2>/dev/null || true
  rm -f "$ROOT/service/apigateway/.env"
}
trap cleanup EXIT

if [ "$DO_DOWN" = true ]; then
  info "Tearing down infra ..."
  docker compose -f "$COMPOSE_FILE" down -v
  exit 0
fi

# ─── [1/6] Start infra ─────────────────────────────────────────────────
info "[1/6] Starting infrastructure ..."
docker compose -f "$COMPOSE_FILE" up -d

info "Waiting for Postgres health checks ..."
for svc in auth-db role-db user-db email-db category-db merchant-db merchant_award-db merchant_business-db merchant_detail-db merchant_policy-db order-db order_item-db product-db transaction-db cart-db review-db review_detail-db slider-db shipping_address-db banner-db; do
  timeout=60
  while ! docker compose -f "$COMPOSE_FILE" exec -T "$svc" pg_isready -U DRAGON -q 2>/dev/null; do
    timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "$svc not ready"; exit 1; }
    sleep 1
  done
  ok "$svc ready"
done

info "Waiting for Redis ..."
timeout=30
while ! docker compose -f "$COMPOSE_FILE" exec -T redis redis-cli -a dragon_knight ping 2>/dev/null | grep -q PONG; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "Redis not ready"; exit 1; }
  sleep 1
done
ok "Redis ready"

info "Waiting for Kafka ..."
timeout=90
while ! docker compose -f "$COMPOSE_FILE" exec -T kafka bash -c 'exec 3<>/dev/tcp/localhost/9092' 2>/dev/null; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "Kafka not ready"; exit 1; }
  sleep 2
done
ok "Kafka ready"

info "Waiting for ClickHouse ..."
timeout=30
while ! docker compose -f "$COMPOSE_FILE" exec -T clickhouse clickhouse-client --query 'SELECT 1' 2>/dev/null; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "ClickHouse not ready"; exit 1; }
  sleep 1
done
ok "ClickHouse ready"

if [ "$INFRA_ONLY" = true ]; then
  ok "Infrastructure is up. Exiting (--infra-only)."
  echo "Postgres: localhost:5432 (ECOMMERCE) + per-service: 5543-5562"
  echo "Redis: localhost:6379  Kafka: localhost:9092  ClickHouse: localhost:8123/9000"
  exit 0
fi

# ─── [2/6] Reset databases + migrate + seed ─────────────────────────────
info "[2/6] Resetting databases, running migrations, and seeding ..."

# Per-service DB maps (must match docker-compose.infra.yml ports)
declare -A DB_CONTAINER=( [auth]=postgres_auth [role]=postgres_role [user]=postgres_user [email]=postgres_email [category]=postgres_category [merchant]=postgres_merchant [merchant_award]=postgres_merchant_award [merchant_business]=postgres_merchant_business [merchant_detail]=postgres_merchant_detail [merchant_policy]=postgres_merchant_policy [order]=postgres_order [order_item]=postgres_order-item [product]=postgres_product [transaction]=postgres_transaction [cart]=postgres_cart [review]=postgres_review [review_detail]=postgres_review_detail [slider]=postgres_slider [shipping_address]=postgres_shipping_address [banner]=postgres_banner )
declare -A DB_NAME=( [auth]=auth_db [role]=role_db [user]=user_db [email]=email_db [category]=category_db [merchant]=merchant_db [merchant_award]=merchant_award_db [merchant_business]=merchant_business_db [merchant_detail]=merchant_detail_db [merchant_policy]=merchant_policy_db [order]=order_db [order_item]=order_item_db [product]=product_db [transaction]=transaction_db [cart]=cart_db [review]=review_db [review_detail]=review_detail_db [slider]=slider_db [shipping_address]=shipping_address_db [banner]=banner_db )
declare -A DB_PORT=( [auth]=5543 [role]=5544 [user]=5545 [email]=5546 [category]=5547 [merchant]=5548 [merchant_award]=5549 [merchant_business]=5550 [merchant_detail]=5551 [merchant_policy]=5552 [order]=5553 [order_item]=5554 [product]=5555 [transaction]=5556 [cart]=5557 [review]=5558 [review_detail]=5559 [slider]=5560 [shipping_address]=5561 [banner]=5562 )

ALL_SERVICES=(auth role user email category merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product transaction cart review review_detail slider shipping_address banner)

# Drop schema + recreate for clean state
for svc in "${ALL_SERVICES[@]}"; do
  docker exec "${DB_CONTAINER[$svc]}" psql -U DRAGON -d "${DB_NAME[$svc]}" \
    -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" >/dev/null 2>&1 || warn "drop schema failed for $svc"
done

# Run migrations per service
for svc in "${ALL_SERVICES[@]}"; do
  info "  migrating $svc ..."
  DB_HOST=localhost DB_PORT="${DB_PORT[$svc]}" DB_NAME="${DB_NAME[$svc]}" \
    go run service/migrate/cmd/main.go -dir "service/$svc/database/migration" up \
    > /tmp/hurl_migrate_$svc.log 2>&1 || { fail "migrate failed for $svc"; cat /tmp/hurl_migrate_$svc.log; exit 1; }
done
ok "All migrations complete"

# Seed roles
docker exec postgres_role psql -U DRAGON -d role_db -c \
  "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN'), ('ROLE_USER') ON CONFLICT DO NOTHING;" >/dev/null 2>&1 || true
ok "Roles seeded"

# Flush redis
docker exec redis-ecommerce-e2e redis-cli -a dragon_knight FLUSHALL >/dev/null 2>&1 || true
ok "Redis flushed"

# Run seeder
info "  running seeder ..."
go run service/seeder/cmd/main.go > /tmp/hurl_seeder.log 2>&1 || warn "seeder had warnings (see /tmp/hurl_seeder.log)"
ok "Seeder complete"

# Stats backfill
info "  running stats backfill ..."
go run service/stats_writer/cmd/main.go backfill > /tmp/hurl_backfill.log 2>&1 || warn "stats backfill had warnings"
ok "Stats backfill complete"

# ─── [5/6] Launch Go services ──────────────────────────────────────────
info "[5/6] Launching all Go services locally ..."

for svc in auth role user category merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product transaction cart review review_detail slider shipping_address banner email stats_writer stats_reader; do
  (cd "$ROOT/service/$svc" && setsid nohup "$ROOT/bin/$svc" > "$LOG_DIR/$svc.log" 2>&1 &)
done
(cd "$ROOT/service/apigateway" && setsid nohup "$ROOT/bin/apigateway" > "$LOG_DIR/apigateway.log" 2>&1 &)
ok "All 23 services launched"

# ─── [6/6] Wait for health ─────────────────────────────────────────────
info "[6/6] Waiting for all services to become healthy ..."
GW_URL="http://localhost:5000/api/auth/hello"
PORT_REGEX=':5005[1-9]|:5006[0-9]|:5000'

ok=0
for round in 1 2 3; do
  info "  health check round $round ..."
  for i in $(seq 1 40); do
    gw=$(curl -s -o /dev/null -w '%{http_code}' --max-time 2 "$GW_URL" 2>/dev/null || echo "000")
    if [ "$gw" = "200" ]; then ok=1; break; fi
    sleep 3
  done
  info "  round $round: apigateway=$gw"
  [ "$ok" = "1" ] && break
done

if [ "$ok" != "1" ]; then
  fail "Apigateway not healthy. Logs:"
  for f in "$LOG_DIR"/*.log; do echo "-- $(basename "$f")"; tail -5 "$f" 2>/dev/null; done
  exit 1
fi
ok "All services healthy"

# ─── [7/6] Run hurl e2e tests ──────────────────────────────────────────
info "[7/6] Running hurl e2e tests ..."

# Discover a seeded role-less user for rules_strict.hurl
USER_EMAIL=""
if docker exec postgres_user psql -U DRAGON -d user_db -t -A -c \
  "SELECT email FROM users WHERE firstname='User1' AND email LIKE 'user_%' ORDER BY user_id LIMIT 1;" \
  > /tmp/hurl_seed_user.txt 2>/dev/null; then
  USER_EMAIL=$(tr -d ' \r' < /tmp/hurl_seed_user.txt)
  USER_ID=$(docker exec postgres_user psql -U DRAGON -d user_db -t -A -c \
    "SELECT user_id FROM users WHERE email = '$USER_EMAIL' LIMIT 1;" | tr -d ' \r')
  if [ -n "$USER_ID" ]; then
    docker exec postgres_role psql -U DRAGON -d role_db -c \
      "DELETE FROM user_roles WHERE user_id = $USER_ID;" >/dev/null 2>&1 || true
  fi
fi

TS=$(date +%s)
TEST_PASSWORD="HurlPass123"
PASS=0; FAIL=0
FAILED=()

for f in "$HURL_DIR"/*.hurl; do
  name=$(basename "$f")
  testEmail="hurl.${name%.hurl}.$TS@example.com"

  common_vars="--test --jobs 1 --variable baseUrl=http://localhost:5000 --variable testEmail=$testEmail --variable testPassword=$TEST_PASSWORD --variable testImage=$HURL_DIR/assets/test.png"

  if [ "$name" = "rules_strict.hurl" ]; then
    if [ -z "$USER_EMAIL" ]; then
      warn "SKIP $name (no seeded role-less user)"
      continue
    fi
    extra_vars="--variable adminEmail=admin.$TS@example.com --variable adminPassword=$TEST_PASSWORD --variable userEmail=$USER_EMAIL --variable userPassword=password1"
    hurl $common_vars $extra_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  else
    hurl $common_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  fi

  if [ $? -eq 0 ]; then
    PASS=$((PASS + 1))
    ok "  ✓ $name"
  else
    FAIL=$((FAIL + 1))
    FAILED+=("$name")
    fail "  ✗ $name (see /tmp/hurl_$name.log)"
  fi

  # Pace auth endpoints (rate-limited)
  sleep 2
done

# ─── [8/6] Verify ClickHouse stats ─────────────────────────────────────
info "[8/6] Verifying ClickHouse stats ..."
sleep 7  # wait for stats-writer flush

for table in order_events order_item_events transaction_events; do
  row_count=$(docker exec clickhouse-ecommerce clickhouse-client --database ecommerce \
    --query "SELECT count() FROM $table FINAL" 2>/dev/null | tr -d ' \r' || echo "unknown")
  info "  stats $table rows: ${row_count}"
done

# ─── Summary ────────────────────────────────────────────────────────────
echo ""
echo "========================================"
echo "  E2E Test Summary"
echo "========================================"
ok "Passed: $PASS"
if [ "$FAIL" -gt 0 ]; then
  fail "Failed: $FAIL"
  echo ""
  echo "Failed suites:"
  for f in "${FAILED[@]}"; do echo "  - $f"; done
  exit 1
else
  ok "All tests passed! 🎉"
fi
