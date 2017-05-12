package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/note"

	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/pagination"
	"github.com/blue-jay/core/router"
)

// Notepad represents the services required for this controller.
type Notepad struct {
	env.Service
}

// LoadNotepad registers the Notepad handlers.
func LoadNotepad(s env.Service) {
	// Create handler.
	h := new(Notepad)
	h.Service = s

	// Load routes.
	c := router.Chain(acl.DisallowAnon)
	h.Router.Get("/notepad", h.Index, c...)
	h.Router.Get("/notepad/create", h.Create, c...)
	h.Router.Post("/notepad/create", h.Store, c...)
	h.Router.Get("/notepad/view/:id", h.Show, c...)
	h.Router.Get("/notepad/edit/:id", h.Edit, c...)
	h.Router.Patch("/notepad/edit/:id", h.Update, c...)
	h.Router.Delete("/notepad/:id", h.Destroy, c...)
}

// Index displays the items.
func (h *Notepad) Index(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	// Create a pagination instance with a max of 10 results.
	p := pagination.New(r, 10)

	items, _, err := note.ByUserIDPaginate(h.DB, id, p.PerPage, p.Offset)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
		items = []note.Item{}
	}

	count, err := note.ByUserIDCount(h.DB, id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
	}

	// Calculate the number of pages.
	p.CalculatePages(count)

	v := h.View.New("note/index")
	v.Vars["items"] = items
	v.Vars["pagination"] = p
	v.Render(w, r)
}

// Create displays the create form.
func (h *Notepad) Create(w http.ResponseWriter, r *http.Request) {
	v := h.View.New("note/create")
	form.Repopulate(r.Form, v.Vars, "name")
	v.Render(w, r)
}

// Store handles the create form submission.
func (h *Notepad) Store(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	if !h.FormValid(w, r, "name") {
		h.Create(w, r)
		return
	}

	_, err := note.Create(h.DB, r.FormValue("name"), id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
		h.Create(w, r)
		return
	}

	h.FlashSuccess(w, r, "Item added.")
	http.Redirect(w, r, "/notepad", http.StatusFound)
}

// Show displays a single item.
func (h *Notepad) Show(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	item, _, err := note.ByID(h.DB, h.Router.Param(r, "id"), id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	v := h.View.New("note/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func (h *Notepad) Edit(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	item, _, err := note.ByID(h.DB, h.Router.Param(r, "id"), id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	v := h.View.New("note/edit")
	form.Repopulate(r.Form, v.Vars, "name")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func (h *Notepad) Update(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	if !h.FormValid(w, r, "name") {
		h.Edit(w, r)
		return
	}

	_, err := note.Update(h.DB, r.FormValue("name"), h.Router.Param(r, "id"), id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
		h.Edit(w, r)
		return
	}

	h.FlashSuccess(w, r, "Item updated.")
	http.Redirect(w, r, "/notepad", http.StatusFound)
}

// Destroy handles the delete form submission.
func (h *Notepad) Destroy(w http.ResponseWriter, r *http.Request) {
	// Get the current user ID.
	id := ""
	if u, ok := env.UserSession(r, h.Sess); ok {
		id = u.ID
	}

	_, err := note.DeleteSoft(h.DB, h.Router.Param(r, "id"), id)
	if err != nil {
		h.FlashErrorGeneric(w, r, err)
	} else {
		h.FlashNotice(w, r, "Item deleted.")
	}

	http.Redirect(w, r, "/notepad", http.StatusFound)
}
