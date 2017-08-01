DROP TABLE IF EXISTS sitemap;
CREATE TABLE sitemap (
    id  SERIAL PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL,
    lat DECIMAL,
    lng DECIMAL
);

DROP TABLE IF EXISTS rainfall;
CREATE TABLE rainfall (
    id SERIAL PRIMARY KEY,
    sitemapid INTEGER,
    rain DECIMAL,
    datetime TIMESTAMP
);