create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    username text not null unique,
    email text not null unique,
    password_hash text not null,
    display_name text not null,
    role text not null default 'admin',
    language text not null default 'en',
    last_login_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists sessions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text not null unique,
    csrf_token text not null,
    ip_address text not null,
    user_agent text not null,
    remember_me boolean not null default false,
    last_seen_at timestamptz not null default now(),
    expires_at timestamptz not null,
    revoked_at timestamptz,
    created_at timestamptz not null default now()
);

create table if not exists domains (
    id uuid primary key default gen_random_uuid(),
    name text not null unique,
    origin_host text not null,
    origin_port integer not null,
    origin_protocol text not null,
    origin_server_name text not null default '',
    enabled boolean not null default true,
    protection_enabled boolean not null default true,
    protection_mode text not null default 'basic',
    challenge_mode text not null default 'cookie',
    cloudflare_mode boolean not null default false,
    ssl_auto_issue boolean not null default false,
    ssl_enabled boolean not null default false,
    force_https boolean not null default false,
    rate_limit_rps integer not null default 20,
    rate_limit_burst integer not null default 40,
    bad_bot_mode boolean not null default true,
    header_validation boolean not null default true,
    js_challenge_enabled boolean not null default false,
    allowed_methods text[] not null default array['GET','POST','PUT','PATCH','DELETE','HEAD','OPTIONS'],
    notes text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists domain_rules (
    id uuid primary key default gen_random_uuid(),
    domain_id uuid not null references domains(id) on delete cascade,
    name text not null,
    type text not null,
    pattern text not null,
    action text not null,
    enabled boolean not null default true,
    priority integer not null default 100,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists ip_lists (
    id uuid primary key default gen_random_uuid(),
    domain_id uuid references domains(id) on delete cascade,
    list_type text not null,
    ip text not null default '',
    cidr text,
    reason text not null default '',
    expires_at timestamptz,
    created_by uuid references users(id) on delete set null,
    created_at timestamptz not null default now()
);

create table if not exists temporary_bans (
    id uuid primary key default gen_random_uuid(),
    domain_id uuid references domains(id) on delete cascade,
    ip text not null,
    reason text not null,
    source text not null default 'auto_ban',
    expires_at timestamptz not null,
    created_at timestamptz not null default now()
);

create table if not exists request_logs (
    id uuid primary key default gen_random_uuid(),
    domain_id uuid references domains(id) on delete set null,
    domain_name text not null,
    client_ip text not null,
    country_code text not null default '',
    method text not null,
    path text not null,
    query_string text not null default '',
    user_agent text not null default '',
    request_id text not null default '',
    decision text not null,
    status_code integer not null,
    block_reason text not null default '',
    response_time_ms integer not null default 0,
    score integer not null default 0,
    challenge_type text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists traffic_rollups (
    window_start timestamptz not null,
    granularity text not null,
    domain_id uuid not null references domains(id) on delete cascade,
    domain_name text not null,
    decision text not null,
    block_reason text not null default '',
    request_count bigint not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    primary key (window_start, granularity, domain_id, decision, block_reason)
);

create table if not exists system_settings (
    key text primary key,
    value text not null,
    type text not null default 'string',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists audit_logs (
    id uuid primary key default gen_random_uuid(),
    user_id uuid references users(id) on delete set null,
    username text not null default 'system',
    action text not null,
    entity_type text not null default '',
    entity_id text not null default '',
    ip_address text not null default '',
    user_agent text not null default '',
    details text not null default '',
    created_at timestamptz not null default now()
);

create table if not exists ssl_certificates (
    id uuid primary key default gen_random_uuid(),
    domain_id uuid not null unique references domains(id) on delete cascade,
    issuer text not null default '',
    status text not null default 'pending',
    cert_path text not null default '',
    key_path text not null default '',
    expires_at timestamptz,
    last_renew_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_sessions_token_hash on sessions(token_hash);
create index if not exists idx_domains_name on domains(name);
create index if not exists idx_domain_rules_domain_id on domain_rules(domain_id);
create index if not exists idx_ip_lists_lookup on ip_lists(list_type, domain_id, ip);
create index if not exists idx_temporary_bans_ip on temporary_bans(ip, expires_at);
create index if not exists idx_request_logs_created_at on request_logs(created_at desc);
create index if not exists idx_request_logs_domain_id_created_at on request_logs(domain_id, created_at desc);
create index if not exists idx_request_logs_decision_created_at on request_logs(decision, created_at desc);
create index if not exists idx_request_logs_client_ip on request_logs(client_ip, created_at desc);
create index if not exists idx_audit_logs_created_at on audit_logs(created_at desc);
