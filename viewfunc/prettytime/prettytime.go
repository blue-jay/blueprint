// Package prettytime provides a funcmap for html/template to that displays
// time using an easy to read format.
package prettytime

import (
	"html/template"

	"github.com/go-sql-driver/mysql"
)

// Map returns a template.FuncMap for PRETTYTIME which outputs a time in this
// format: 3:04 PM 01/02/2006.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["PRETTYTIME"] = func(createdAt mysql.NullTime, updatedAt mysql.NullTime) string {
		if updatedAt.Valid {
			return updatedAt.Time.Format("3:04 PM 01/02/2006")
		}

		return createdAt.Time.Format("3:04 PM 01/02/2006")
	}

	return f
}
