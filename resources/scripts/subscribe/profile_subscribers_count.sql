select count(*) from subscriptions as s
where s.profile_username = $1;