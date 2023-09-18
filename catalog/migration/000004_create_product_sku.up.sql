BEGIN;

CREATE TABLE IF NOT EXISTS PRODUCT_SKU
(
    PRODUCT_ID BIGINT NOT NULL,
    SKU VARCHAR(64) NOT NULL,

    CONSTRAINT FK_PRODUCT_ID FOREIGN KEY (PRODUCT_ID) REFERENCES PRODUCT (PRODUCT_ID),
    PRIMARY KEY(PRODUCT_ID, SKU)
);

COMMIT;