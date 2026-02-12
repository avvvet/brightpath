# BrightTurn - Foreclosure Lead Pipeline

Automated pipeline that finds distressed property owners facing foreclosure, gets their phone numbers via skip tracing, and saves qualified wholesale leads to PostgreSQL.

## Tech Stack

- **Go** - Backend pipeline + HTTP API
- **PostgreSQL + PostGIS** - Database with spatial queries
- **RealEstateAPI** - Property data + skip tracing
- **Twilio** - SMS outreach (Phase 2)

## Quick Start
```bash
# Setup database
docker-compose up -d
docker exec -it postgres_db psql -U postgres -d brightpath -c "CREATE EXTENSION IF NOT EXISTS postgis;"

# Run migrations
go run cmd/migrate/main.go

# Run pipeline for a state
go run cmd/worker/main.go --state=TX

# Start API server
go run cmd/server/main.go
```

## Environment Variables
```bash
REAPI_API_KEY=your_key_here
DATABASE_URL=postgres://postgres:postgres@localhost:5432/brightpath?sslmode=disable
PIPELINE_AUCTION_DAYS_MIN=3
PIPELINE_AUCTION_DAYS_MAX=21
PIPELINE_EQUITY_MIN=25
PIPELINE_DAILY_HOUR=6
```

## How It Works

### 3-Step Pipeline (Credit-Safe)

**Step 1: Discovery (FREE)**
- Uses `ids_only: true` to get property IDs without cost
- Searches for properties with auctions 3-21 days out
- Filters: 25%+ equity, SFR, non-corporate, TX/FL/GA/AZ/OH

**Step 2: Enrich NEW Properties Only (1 credit each)**
- Checks database first (credit guard)
- Only fetches PropertyDetail for new properties
- Applies Go filters: firstName exists, owner type = "Individual", not investor

**Step 3: Skip Trace Qualified Leads (1 credit each)**
- Checks if lead already exists (credit guard)
- Fetches phone numbers for qualified properties
- Filters: mobile only, connected, not DNC
- Fallback: Uses relative/co-owner phones if owner phone unavailable

### Idempotent Design

Running the pipeline multiple times = **0 extra credits**
- Properties: `UNIQUE(source, source_id)` prevents duplicates
- Leads: `UNIQUE(property_id)` prevents duplicates
- Skip trace status tracking prevents retries

## Skip Trace Status Flow

| Status | Meaning | Will Retry? |
|--------|---------|-------------|
| `pending` | Never attempted | ✅ YES |
| `completed` | Lead created successfully | ❌ NO |
| `no_match` | Person not found (404) | ❌ NO |
| `no_phone` | No valid mobile phone | ❌ NO |
| `missing_data` | Missing address fields | ❌ NO |
| `api_error` | API failures (429, 500, etc) | ✅ YES |
| `rate_limit` | Monthly limit exceeded (429) | ❌ NO (until upgrade) |

**After upgrading API plan, reset rate limits:**
```sql
UPDATE properties 
SET skip_trace_status = 'pending', skip_trace_error = NULL 
WHERE skip_trace_status = 'rate_limit';
```

## Database Schema

### Key Tables

**properties** - Property details, equity, foreclosure info
- Dedup: `UNIQUE(source, source_id)`
- Tracks: `skip_trace_status`, `skip_trace_attempted_at`

**leads** - Contact info, phone numbers, demographics
- Dedup: `UNIQUE(property_id)`
- Fields: `phone`, `phone_owner_match`, `days_to_auction`

**pipeline_runs** - Execution logs, stats, credit tracking
- Tracks: IDs found, leads created, credits used

## Query Your Leads
```sql
-- Get top leads by priority
SELECT 
    CASE 
        WHEN p.equity_percent >= 80 THEN '🔥 HOT'
        WHEN p.equity_percent >= 60 THEN '⚡ High'
        ELSE '✓ Good'
    END as priority,
    p.owner1_first_name || ' ' || p.owner1_last_name as owner,
    l.phone,
    CASE WHEN l.phone_owner_match THEN '✓ Direct' ELSE '→ Relative' END as contact,
    l.days_to_auction || ' days' as urgency,
    p.equity_percent || '%' as equity,
    p.address,
    p.city,
    p.auction_date
FROM leads l
JOIN properties p ON l.property_id = p.id
WHERE l.status = 'new'
ORDER BY p.equity_percent DESC, l.days_to_auction ASC
LIMIT 20;
```

## API Endpoints
```bash
# Health check
curl http://localhost:8080/health

# Get leads
curl http://localhost:8080/api/leads?status=new&limit=50

# Get pipeline runs
curl http://localhost:8080/api/pipeline/runs?state=TX&limit=10

# Trigger pipeline run
curl -X POST http://localhost:8080/api/pipeline/run?state=FL
```

## Expected Results

### Per State (Daily Runs)

**Texas:**
- ~250-300 leads/month
- ~1,200-1,500 credits/month

**Florida:**
- ~400-500 leads/month
- ~2,000-2,500 credits/month

### Multi-State (5 States Daily)

- ~1,200-1,500 leads/month
- ~6,000-8,000 credits/month
- Expected: 20-40 deals/month (at 2-3% conversion)

## Cost Breakdown

### Monthly (5 States)

- Property credits: ~1,500 (new properties)
- Skip trace credits: ~1,200 (qualified leads)
- **Total: ~2,700 credits/month**

### Per Lead Economics

- Property data: $0.XX per record
- Skip trace: $0.XX per lookup
- SMS (Twilio): $0.0079 per message
- **Total cost per contactable lead: ~$0.XX**

## Troubleshooting

### "Monthly usage limit exceeded (429)"
- Upgrade your RealEstateAPI plan
- Current properties with `rate_limit` status will be retried after upgrade

### "Property missing required skip trace fields"
- Some properties don't have mail addresses in database
- This is normal (~5-10% of properties)

### "Person not found (404)"
- Skip trace database doesn't have info on this person
- Marked as `no_match`, won't retry (saves credits)

### Multiple runs showing 0 new properties
- Correct behavior! Deduplication working
- Properties only processed once
- New auction notices appear daily

## Production Deployment

### Daily Schedule (Cron)
```bash
# Run each state at 6 AM daily
0 6 * * * /usr/local/bin/worker --state=TX
0 6 * * * /usr/local/bin/worker --state=FL
0 6 * * * /usr/local/bin/worker --state=GA
0 6 * * * /usr/local/bin/worker --state=AZ
0 6 * * * /usr/local/bin/worker --state=OH
```

### Monitoring
```sql
-- Daily stats
SELECT 
    state,
    COUNT(*) as total_runs,
    SUM(leads_created) as total_leads,
    SUM(credits_used_property + credits_used_skiptrace) as total_credits
FROM pipeline_runs
WHERE started_at > NOW() - INTERVAL '7 days'
GROUP BY state
ORDER BY total_leads DESC;
```

## Next Features

- [ ] Twilio SMS integration
- [ ] Lead scoring (equity %, urgency, owner type)
- [ ] Multi-state scheduler
- [ ] Email alerts on pipeline failures
- [ ] Dashboard UI (Svelte)

## Support States

Currently: **TX, FL, GA, AZ, OH**

To add more states:
1. Check RealEstateAPI coverage
2. Verify foreclosure laws (non-judicial preferred)
3. Run: `go run cmd/worker/main.go --state=NEW_STATE`

