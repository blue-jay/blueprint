// Package notepad provides a simple CRUD application in a web page.
package notepad

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model"
	"github.com/blue-jay/blueprint/model/note"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/notepad"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	// On the template generation, Note needs to be capitalized
	//items, _, err := model.Note.ByUserID(c.UserID)
	if err != nil {
		c.FlashError(err)
		items = []note.Item{}
	}

	v := c.View.New("note/index")
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("note/create")
	c.Repopulate(v.Vars, "name")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Create(w, r)
		return
	}

	_, err := model.Note.Create(r.FormValue("name"), c.UserID)
	if err != nil {
		c.FlashError(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Item added.")
	c.Redirect(uri)
}

// Show displays a single item.
func Show(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := model.Note.ByID(c.Param("id"), c.UserID)
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("note/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := model.Note.ByID(c.Param("id"), c.UserID)
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("note/edit")
	c.Repopulate(v.Vars, "name")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Edit(w, r)
		return
	}

	_, err := model.Note.Update(r.FormValue("name"), c.Param("id"), c.UserID)
	if err != nil {
		c.FlashError(err)
		Edit(w, r)
		return
	}

	c.FlashSuccess("Item updated.")
	c.Redirect(uri)
}

// Destroy handles the delete form submission.
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	_, err := model.Note.DeleteSoft(c.Param("id"), c.UserID)
	if err != nil {
		c.FlashError(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}
