
CREATE TABLE pipeline_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state TEXT NOT NULL,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'running',

    -- stats
    ids_found INT DEFAULT 0,
    ids_new INT DEFAULT 0,
    properties_enriched INT DEFAULT 0,
    properties_qualified INT DEFAULT 0,
    leads_created INT DEFAULT 0,
    skip_trace_hits INT DEFAULT 0,
    credits_used_property INT DEFAULT 0,
    credits_used_skiptrace INT DEFAULT 0,

    error TEXT,
    config JSONB,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);