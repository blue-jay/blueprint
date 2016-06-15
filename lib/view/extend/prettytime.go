package extend

import (
	"html/template"

	"github.com/go-sql-driver/mysql"
)

// PrettyTime returns a template.FuncMap
// * PRETTYTIME outputs a nice time format showing the updated time or the created time
func PrettyTime() template.FuncMap {
	f := make(template.FuncMap)

	f["PRETTYTIME"] = func(createdAt mysql.NullTime, updatedAt mysql.NullTime) string {

		if updatedAt.Valid {
			return updatedAt.Time.Format("3:04 PM 01/02/2006")
		}

		return createdAt.Time.Format("3:04 PM 01/02/2006")
	}

	return f
}
