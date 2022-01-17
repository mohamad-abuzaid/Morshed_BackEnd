// file: datasource/users.go

package datasource

import (
	"errors"
	"fmt"

	"morshed/data/engine/sql"
	"morshed/helpers"
)

// Engine is from where to fetch the data, in this case the users.
type Engine uint32

const (
	// Memory stands for simple memory location;
	// map[int64] datamodels.User ready to use, it's our source in this example.
	Memory Engine = iota
	// Bolt for boltdb source location.
	Bolt
	// MySQL for mysql-compatible source location.
	MySQL
)

// LoadUsers returns all users(empty map) from the memory, for the sake of simplicty.
func StartMySql(engine Engine) (sql.Database, error) {
	if engine != MySQL {
		return nil, errors.New("Only MySql available")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		helpers.Mgetenv("MYSQL_USER", "root"),
		helpers.Mgetenv("MYSQL_PASSWORD", "BugSquad#2022"),
		helpers.Mgetenv("MYSQL_HOST", "127.0.0.1"),
		helpers.Mgetenv("MYSQL_DATABASE", "morshed-db"),
	)

	db, err := sql.ConnectMySQL(dsn)
	if err != nil {
		return nil, errors.Unwrap(fmt.Errorf("error connecting to the MySQL database:  %w", err))
	}

	return db, nil
}
