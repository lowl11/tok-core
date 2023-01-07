select
    username,
    ip_address
from user_ips
where ip_address = $1;