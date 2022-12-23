select
    post.author_username,
    u.name author_name,
    u.avatar author_avatar,

    category.code category_code,
    category.name category_name,

    post.code,
    post.text,
    post.picture,
    post.picture_width,
    post.picture_height,
    post.created_at
from posts as post
    inner join post_categories as category on post.category_code = category.code
    inner join users as u on (u.username = post.author_username)
where post.author_username = $1
order by post.created_at desc;