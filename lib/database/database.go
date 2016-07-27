// Package database provides a wrapper around the sqlx package.
package database

import (
	"errors"
	"fmt"
	"strings"
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
	// FileStorage is the path to the folder where uploaded files are stored
	FileStorage string
	// Type of database
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
}

// MySQLInfo holds the details for the database connection.
type MySQLInfo struct {
	Username  string
	Password  string
	Database  string
	Charset   string
	Collation string
	Hostname  string
	Port      int
	Parameter string
}

// Connect to the database.
func Connect(i Info, connectDatabase bool) error {
	var err error

	// Store the config
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()

	switch Config().Type {
	case TypeMySQL:
		// Connect to MySQL and ping
		if SQL, err = sqlx.Connect("mysql", dsn(Config().MySQL, connectDatabase)); err != nil {
			return err
		}
	default:
		return errors.New("No registered database in config")
	}

	return err
}

// Disconnect closes the MySQL connection.
func Disconnect() error {
	return SQL.Close()
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
	defer infoMutex.RUnlock()
	return info
}

// DSN returns the Data Source Name.
func dsn(ci MySQLInfo, includeDatabase bool) string {
	// Set defaults
	ci = setDefaults(ci)

	// Build parameters
	param := ci.Parameter

	// If parameter is specified, add a question mark
	// Don't add one if a question mark is already there
	if len(ci.Parameter) > 0 && !strings.HasPrefix(ci.Parameter, "?") {
		param = "?" + ci.Parameter
	}

	// Add collation
	if !strings.Contains(param, "collation") {
		if len(param) > 0 {
			param += "&collation=" + ci.Collation
		} else {
			param = "?collation=" + ci.Collation
		}
	}

	// Add charset
	if !strings.Contains(param, "charset") {
		if len(param) > 0 {
			param += "&charset=" + ci.Charset
		} else {
			param = "?charset=" + ci.Charset
		}
	}

	// Example: root:password@tcp(localhost:3306)/test
	s := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", ci.Username, ci.Password, ci.Hostname, ci.Port, param)

	if includeDatabase {
		s = fmt.Sprintf("%v:%v@tcp(%v:%d)/%v%v", ci.Username, ci.Password, ci.Hostname, ci.Port, ci.Database, param)
	}

	return s
}

// Create a new database with the charset and collation
func Create(ci MySQLInfo) error {
	// Set defaults
	ci = setDefaults(ci)

	// Create the database
	_, err := SQL.Exec(fmt.Sprintf(`CREATE DATABASE %v
				DEFAULT CHARSET = %v
				COLLATE = %v
				;`, ci.Database,
		ci.Charset,
		ci.Collation))
	return err
}

// setDefaults sets the charset and collation if they are not set
func setDefaults(ci MySQLInfo) MySQLInfo {
	if len(ci.Charset) == 0 {
		ci.Charset = "utf8"
	}
	if len(ci.Collation) == 0 {
		ci.Collation = "utf8_unicode_ci"
	}

	return ci
}
