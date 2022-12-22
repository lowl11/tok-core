select
    ipv4_address
from user_ips
where username = $1;