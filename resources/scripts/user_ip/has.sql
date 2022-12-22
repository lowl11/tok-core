select * from user_ips
where username = $1 and ipv4_address = $2;