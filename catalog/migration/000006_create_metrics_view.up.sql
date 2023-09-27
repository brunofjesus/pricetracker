CREATE VIEW PRODUCT_METRICS AS (
    SELECT p.product_id,
        p.price,
        (p.price - AVG(pp.price)) / p.price as discount_percent,
        (p.price - AVG(pp.price)) as discount,
        AVG(pp.price) average,
        MAX(pp.price) maximum,
        MIN(pp.price) minimum,
        COUNT(pp.price) entries,
        (now() - interval '30 day') metrics_since
    FROM product_price pp
        INNER JOIN public.product p on p.product_id = pp.product_id
    WHERE pp.date_time > now() - interval '30 day'
    GROUP BY p.product_id,
        pp.product_id
)