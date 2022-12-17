select count(*) from subscriptions as s
where s.subscribe_username = $1;