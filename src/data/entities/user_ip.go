package entities

type UserIpGet struct {
	Username  string `db:"username"`
	IpAddress string `db:"ip_address"`
}

type UserIpNew struct {
	Username  string `db:"username"`
	IpAddress string `db:"ip_address"`
}
