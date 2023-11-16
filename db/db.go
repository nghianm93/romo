package db

const (
	DBNAME   = "MONGO_DB_NAME"
	userColl = "users"
	hostColl = "hosts"
)

type Store struct {
	User UserStore
	Host HostStore
}
