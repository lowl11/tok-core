select
    u.username,
    u.name,
    u.avatar
from users as u
where u.username = any($1);