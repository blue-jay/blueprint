// Package database provides a wrapper around the sqlx package.
package database

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	// SQL wrapper
	SQL *sqlx.DB
	// Database info
	info      Info
	infoMutex sync.RWMutex
)

// Type is the type of database.
type Type string

const (
	// TypeMySQL is MySQL
	TypeMySQL Type = "MySQL"
)

// Info contains the database configurations.
type Info struct {
	// Database type
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
}

// MySQLInfo holds the details for the database connection.
type MySQLInfo struct {
	Username  string
	Password  string
	Database  string
	Hostname  string
	Port      int
	Parameter string
}

// Connect to the database.
func Connect(i Info) error {
	var err error

	// Store the config
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()

	switch Config().Type {
	case TypeMySQL:
		// Connect to MySQL and ping
		if SQL, err = sqlx.Connect("mysql", dsn(Config().MySQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}
	default:
		log.Println("No registered database in config")
	}

	return err
}

// ResetConfig removes the config.
func ResetConfig() {
	infoMutex.Lock()
	info = Info{}
	infoMutex.Unlock()
}

// Config returns the config.
func Config() Info {
	infoMutex.RLock()
	i := info
	infoMutex.RUnlock()
	return i
}

// DSN returns the Data Source Name.
func dsn(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Database + ci.Parameter
}
