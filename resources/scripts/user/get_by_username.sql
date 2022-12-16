select
    id,
    username,
    password,
    name,
    bio,
    avatar,
    wallpaper
from users
where username = $1;