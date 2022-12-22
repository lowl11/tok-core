select
    username
from user_ips
where ipv4_address = $1;