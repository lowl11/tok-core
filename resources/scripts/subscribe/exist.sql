select id from subscriptions
where profile_username = $1 and subscribe_username = $2;