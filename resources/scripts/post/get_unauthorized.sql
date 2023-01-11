select
    post.author_username,
    u.name author_name,
    u.avatar author_avatar,

    category.code category_code,
    category.name category_name,

    post.code,
    post.text,
    post.picture,
    post.created_at
from posts as post
         inner join post_categories as category on post.category_code = category.code
         inner join users as u on (u.username = post.author_username)
order by post.created_at desc
offset 0
limit 1000;