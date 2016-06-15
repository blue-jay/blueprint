package notepad

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flash"
	"github.com/blue-jay/blueprint/lib/form"
	"github.com/blue-jay/blueprint/lib/middleware/acl"
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
	"github.com/blue-jay/blueprint/model/note"
)

var (
	uri = "/notepad"
)

func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri, Store, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the notes
func Index(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])

	items, err := note.ByUserID(userID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		items = []note.Item{}
	}

	v := view.New("note/index")
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form
func Create(w http.ResponseWriter, r *http.Request) {
	v := view.New("note/create")
	v.Vars["method"] = "POST"
	form.Repopulate(r.Form, v.Vars, "note")
	v.Render(w, r)
}

// Store handles the create form submission
func Store(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])

	if validate, missingField := form.Required(r, "note"); !validate {
		sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Warning})
		sess.Save(r, w)
		Create(w, r)
		return
	}

	_, err := note.Create(r.FormValue("note"), userID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		Create(w, r)
		return
	}

	sess.AddFlash(flash.Info{"Item added.", flash.Success})
	sess.Save(r, w)
	http.Redirect(w, r, uri, http.StatusFound)
}

// Show displays a single note
func Show(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])
	params := router.Params(r)

	item, err := note.ByID(params.ByName("id"), userID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	v := view.New("note/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form
func Edit(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])
	params := router.Params(r)

	item, err := note.ByID(params.ByName("id"), userID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		http.Redirect(w, r, uri, http.StatusFound)
		return
	}

	v := view.New("note/edit")
	v.Vars["method"] = "PATCH"
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission
func Update(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])
	params := router.Params(r)

	if validate, missingField := form.Required(r, "note"); !validate {
		sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Warning})
		sess.Save(r, w)
		Edit(w, r)
		return
	}

	_, err := note.Update(r.FormValue("note"), params.ByName("id"), userID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		Edit(w, r)
		return
	}

	sess.AddFlash(flash.Info{"Item updated.", flash.Success})
	sess.Save(r, w)
	http.Redirect(w, r, uri, http.StatusFound)
}

// Destroy handles the delete form submission
func Destroy(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	userID := fmt.Sprintf("%v", sess.Values["id"])
	params := router.Params(r)

	_, err := note.Delete(params.ByName("id"), userID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
	} else {
		sess.AddFlash(flash.Info{"Item deleted.", flash.Notice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, uri, http.StatusFound)
}
