package entities

type UserIpNew struct {
	Username  string `db:"username"`
	IpAddress string `db:"ip_address"`
}
