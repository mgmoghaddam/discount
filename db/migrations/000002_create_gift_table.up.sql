CREATE TABLE "gift"(
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE NOT NULL,
    gift_amount     DECIMAL(20, 0),
    usage_limit     INT                 NOT NULL DEFAULT 0,
    used_count      INT                          DEFAULT 0,
    expiration_date TIMESTAMPTZ,
    start_date_time TIMESTAMPTZ,
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ         NOT NULL DEFAULT now()
);


CREATE INDEX ON gift (code);
CREATE INDEX ON gift (gift_amount);
CREATE INDEX ON gift (expiration_date);
CREATE INDEX ON gift (start_date_time);