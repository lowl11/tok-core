select
    id,
    username,
    password,
    name,
    bio
from users
where username = :username and password = :password;