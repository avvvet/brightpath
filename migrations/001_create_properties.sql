CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source TEXT NOT NULL DEFAULT 'reapi',
    source_id TEXT NOT NULL,

    -- address (from propertyInfo.address)
    address TEXT NOT NULL,
    street TEXT,
    city TEXT,
    state TEXT NOT NULL,
    zip TEXT,
    county TEXT,
    fips TEXT,
    location GEOGRAPHY(POINT, 4326),

    -- property details (from propertyInfo)
    property_type TEXT,
    bedrooms INT,
    bathrooms NUMERIC(4,1),
    sqft INT,
    year_built INT,
    lot_sqft INT,

    -- valuation
    avm INT,
    assessed_value INT,

    -- mortgage & equity (from root + currentMortgages)
    open_mortgage_balance INT,
    equity_amount INT,
    equity_percent INT,
    loan_type TEXT,
    lender_name TEXT,

    -- foreclosure (from root + auctionInfo)
    pre_foreclosure BOOLEAN DEFAULT false,
    auction BOOLEAN DEFAULT false,
    auction_date DATE,
    auction_time TEXT,
    auction_location TEXT,
    notice_type TEXT,
    recording_date DATE,
    foreclosure BOOLEAN DEFAULT false,

    -- auction enrichment (from auctionInfo)
    bank_estimated_value INT,
    trustee_name TEXT,
    trustee_phone TEXT,
    foreclosure_document_type TEXT,

    -- owner (from ownerInfo)
    owner1_first_name TEXT,
    owner1_last_name TEXT,
    owner1_full_name TEXT,
    owner1_type TEXT,
    owner2_first_name TEXT,
    owner2_last_name TEXT,
    owner_occupied BOOLEAN,
    absentee_owner BOOLEAN,
    investor_buyer BOOLEAN,
    corporate_owned BOOLEAN,
    ownership_length_months INT,
    mail_street TEXT,
    mail_city TEXT,
    mail_state TEXT,
    mail_zip TEXT,

    -- extras for future LLM scoring
    suggested_rent TEXT,
    median_income TEXT,
    flood_zone BOOLEAN,
    hoa BOOLEAN,
    tax_lien BOOLEAN,
    tax_amount TEXT,
    tax_delinquent_year TEXT,
    years_owned INT,
    last_sale_date DATE,

    -- Skip trace tracking (to avoid re-running the skip trace api)
    skip_trace_status TEXT DEFAULT 'pending',
    skip_trace_attempted_at TIMESTAMPTZ,
    skip_trace_error TEXT,

    -- metadata
    raw_data JSONB,
    discovered_at TIMESTAMPTZ DEFAULT NOW(),
    enriched_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(source, source_id)
);

CREATE INDEX idx_properties_state_county ON properties(state, county);
CREATE INDEX idx_properties_equity ON properties(equity_percent) WHERE equity_percent >= 25;
CREATE INDEX idx_properties_auction_date ON properties(auction_date) WHERE auction_date IS NOT NULL;
CREATE INDEX idx_properties_location ON properties USING GIST(location);
CREATE INDEX idx_properties_source_id ON properties(source_id);