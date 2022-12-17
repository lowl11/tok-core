select
    s.subscribe_username username
from subscriptions as s
where s.profile_username = $1;