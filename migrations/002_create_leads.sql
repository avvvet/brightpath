CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL REFERENCES properties(id),

    -- status flow: new → qualified → sms_sent → replied → hot → forwarded → contracted → closed → dead
    status TEXT NOT NULL DEFAULT 'new',

    -- contact (from skip trace)
    phone TEXT,
    phone_type TEXT,
    phone_connected BOOLEAN,
    phone_dnc BOOLEAN,
    phone_owner_match BOOLEAN,
    alt_phones JSONB,
    email TEXT,
    alt_emails JSONB,

    -- demographics (from skip trace)
    owner_age INT,
    owner_gender TEXT,

    -- scoring (Phase 2)
    lead_score INT,
    score_reasons TEXT[],
    days_to_auction INT,

    -- deal tracking (Phase 3)
    closer_id UUID,
    forwarded_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(property_id)
);

CREATE INDEX idx_leads_status ON leads(status);
CREATE INDEX idx_leads_phone ON leads(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_leads_score ON leads(lead_score DESC) WHERE lead_score IS NOT NULL;