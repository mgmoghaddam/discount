CREATE TABLE "discount"
(
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE NOT NULL,
    percent_off     INT,
    discount_amount DECIMAL(20, 0),
    usage_limit     INT                 NOT NULL,
    used_count      INT                          DEFAULT 0,
    expiration_date TIMESTAMPTZ,
    start_date_time TIMESTAMPTZ,
    max_amount      DECIMAL(20, 0),
    min_amount      DECIMAL(20, 0),
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ         NOT NULL DEFAULT now()
);


CREATE INDEX ON discount (code);
CREATE INDEX ON discount (expiration_date);
CREATE INDEX ON discount (start_date_time);
CREATE INDEX ON discount (percent_off);
CREATE INDEX ON discount (discount_amount);