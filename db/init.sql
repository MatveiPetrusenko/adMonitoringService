-- Creation of advertisement table
CREATE TABLE IF NOT EXISTS advertisement (
    id SERIAL,
    ad_id character varying(24),
    name character varying(1024),
    currency character varying(24),
    price character varying(24),
    link character varying(2048)
);

