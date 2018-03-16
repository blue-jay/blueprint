// Package pagination assists with navigating between pages of results.
package pagination

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Info holds the pagination fields.
type Info struct {
	Page       int
	TotalPages int
	PerPage    int
	Offset     int
}

// New returns a pagination struct.
func New(r *http.Request, perPage int) *Info {
	var err error
	info := &Info{
		PerPage: perPage,
	}

	info.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || info.Page < 1 {
		info.Page = 1
	}

	if info.Page > 1 {
		info.Offset = (info.Page - 1) * info.PerPage
	}

	return info
}

// CalculatePages calculates the number of pages by passing in the item total.
func (i *Info) CalculatePages(itemTotal int) {
	i.TotalPages = itemTotal / i.PerPage
	if itemTotal%i.PerPage != 0 {
		i.TotalPages++
	}
}

// Map returns a template.FuncMap PAGINATION which makes it easy to navigate
// between pages of results.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["PAGINATION"] = func(info Info, m map[string]interface{}) template.HTML {

		currentURI, ok := m["CurrentURI"]
		if !ok {
			log.Println("Issue")
			return template.HTML("Pagination could not load because CurrentURI is missing.")
		}

		top := `<nav aria-label="Page navigation"><ul class="pagination">`
		middle := ""
		bottom := `</ul></nav>`

		for i := 1; i <= info.TotalPages; i++ {
			if i == info.Page {
				middle += fmt.Sprintf(`<li class="active"><a href="%v?page=%v">%v</a></li>`, currentURI, i, i)
			} else {
				middle += fmt.Sprintf(`<li><a href="%v?page=%v">%v</a></li>`, currentURI, i, i)
			}

		}

		return template.HTML(top + middle + bottom)
	}

	return f
}
