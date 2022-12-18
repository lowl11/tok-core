select
    s.subscribe_username username,
    u.name,
    u.avatar
from subscriptions as s
    inner join users as u on (u.username = s.subscribe_username)
where s.profile_username = $1;