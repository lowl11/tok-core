select
    category.name,
    category.code
from post_categories as category
where code = $1;