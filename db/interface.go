package db

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

// User represents API user including admin and regular user
type User struct {
	ID         string
	Name       string
	Password   string
	CreatedOn  time.Time
	ModifiedOn mysql.NullTime
}
