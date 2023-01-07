select
    username,
    ip_address
from user_ips
where username = $1;