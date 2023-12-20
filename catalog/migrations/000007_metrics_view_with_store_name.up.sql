DROP VIEW IF EXISTS product_metrics;

CREATE VIEW product_metrics
            (product_id, store_id, store_name, store_slug, store_website, name, brand, price, available,
             image_url, product_url, discount_percent, diff, average, maximum, minimum, entries, metrics_since)
AS
SELECT p.product_id,
       p.store_id,
       s.name                                    as store_name,
       s.slug                                    as store_slug,
       s.website                                 as store_website,
       p.name,
       p.brand,
       p.price,
       p.available,
       p.image_url,
       p.product_url,
       (avg(pp.price) - p.price) / avg(pp.price) AS discount_percent,
       p.price - avg(pp.price)                   AS diff,
       avg(pp.price)                             AS average,
       max(pp.price)                             AS maximum,
       min(pp.price)                             AS minimum,
       count(pp.price)                           AS entries,
       now() - '30 days'::interval               AS metrics_since
FROM product_price pp
         JOIN product p ON p.product_id = pp.product_id
         JOIN store s ON p.store_id = s.store_id
WHERE pp.date_time > (now() - '30 days'::interval)
GROUP BY p.product_id, pp.product_id, s.name, s.slug, s.website;



