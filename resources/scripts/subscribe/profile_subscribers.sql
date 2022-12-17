select
    s.profile_username username
from subscriptions as s
where s.subscribe_username = $1;