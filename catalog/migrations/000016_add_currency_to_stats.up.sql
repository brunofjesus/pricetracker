DROP VIEW IF EXISTS product_with_stats;

CREATE OR REPLACE VIEW product_with_stats
            (product_id, store_id, store_name, store_slug, store_website, name, brand, price, currency, available,
             image_url, product_url, discount_percent, difference, average, maximum, minimum, entries)
AS
SELECT p.product_id,
       p.store_id,
       s.name                                    as store_name,
       s.slug                                    as store_slug,
       s.website                                 as store_website,
       p.name,
       p.brand,
       p.price,
       p.currency,
       p.available,
       p.image_url,
       p.product_url,
       st.discount_percent,
       st.difference,
       st.average,
       st.maximum,
       st.minimum,
       st.entries
FROM product p
         JOIN store s ON p.store_id = s.store_id
         JOIN product_stats st ON p.product_id = st.product_id
