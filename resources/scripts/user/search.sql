select
    id,
    username,
    password,
    name,
    bio,
    avatar,
    wallpaper
from users
where
    username ilike concat(cast($1 as varchar), '%') or name ilike concat(cast($1 as varchar), '%');