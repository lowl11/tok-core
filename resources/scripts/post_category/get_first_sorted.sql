select
    category.name,
    category.code
from post_categories as category
order by category.name
limit $1;