CREATE TABLE app_user (
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR(100) UNIQUE NOT NULL,
    name       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE api_key (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  key_hash     BYTEA NOT NULL,
  scopes       TEXT[] NOT NULL DEFAULT '{shorten,read,admin}',
  rate_limit   INT NOT NULL DEFAULT 120,
  quota_month  INT NOT NULL DEFAULT 5000,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  revoked_at   TIMESTAMPTZ
);

CREATE TABLE link (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES app_user(id) ON DELETE SET NULL,
  alias        VARCHAR(64) UNIQUE NOT NULL,    -- slug: [a-zA-Z0-9_-]
  target_url   TEXT NOT NULL CHECK (length(target_url) <= 4096),
  title        TEXT,
  tags         TEXT[] DEFAULT '{}',
  is_active    BOOLEAN NOT NULL DEFAULT true,
  expire_at    TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE click (
  id           BIGSERIAL PRIMARY KEY,
  link_id      BIGINT NOT NULL REFERENCES link(id) ON DELETE CASCADE,
  ts           TIMESTAMPTZ NOT NULL DEFAULT now(),
  ip_hash      BYTEA,
  ua           TEXT,
  referrer     TEXT,
  country      TEXT,
  utm_source   TEXT,
  utm_medium   TEXT,
  utm_campaign TEXT
);

CREATE INDEX idx_click_linkid_ts ON click(link_id, ts DESC);

CREATE TABLE click_daily (
  link_id      BIGINT NOT NULL REFERENCES link(id) ON DELETE CASCADE,
  day          DATE NOT NULL,
  total        BIGINT NOT NULL,
  primary key (link_id, day)
);

CREATE INDEX idx_link_userid ON link(user_id);
CREATE INDEX idx_link_active_expire ON link(is_active, expire_at);