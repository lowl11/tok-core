select
    s.profile_username username,
    u.name,
    u.avatar
from subscriptions as s
    inner join users as u on (u.username = s.profile_username)
where s.subscribe_username = $1;