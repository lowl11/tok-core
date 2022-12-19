select
    category.name,
    category.code
from post_categories as category
where category.name ilike concat(cast($1 as varchar), '%') or category.code ilike concat(cast($1 as varchar), '%');