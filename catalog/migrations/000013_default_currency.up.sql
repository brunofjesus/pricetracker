BEGIN;

UPDATE PRODUCT SET CURRENCY = 'EUR';
UPDATE PRODUCT_PRICE SET CURRENCY = 'EUR';

COMMIT;