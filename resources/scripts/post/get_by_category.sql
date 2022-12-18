select
    post.author_username,
    u.name auhtor_name,
    u.avatar auhtor_avatar,

    category.code category_code,
    category.name category_name,

    post.text,
    post.picture,
    post.created_at
from posts as post
         inner join post_categories as category on post.category_code = category.code
         inner join users as u on (u.username = post.author_username)
where post.category_code = $1;