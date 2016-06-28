package flight

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flash"
	"github.com/blue-jay/blueprint/lib/form"
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/session"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

// Info contains commonly used information
type Info struct {
	Sess   *sessions.Session
	Params httprouter.Params
	UserID string
	W      http.ResponseWriter
	R      *http.Request
}

// Context returns commonly used information
func Context(w http.ResponseWriter, r *http.Request) *Info {
	sess := session.Instance(r)
	params := router.Params(r)

	return &Info{
		Sess:   sess,
		Params: params,
		UserID: fmt.Sprintf("%v", sess.Values["id"]),
		W:      w,
		R:      r,
	}
}

// Param gets the URL parameter
func (c *Info) Param(name string) string {
	return c.Params.ByName(name)
}

// Redirect sends a temporary redirect
func (c *Info) Redirect(urlStr string) {
	http.Redirect(c.W, c.R, urlStr, http.StatusFound)
}

// FormValid determines if the user submitted all the required fields and then
// saves an error flash. Returns true if form is valid.
func (c *Info) FormValid(fields ...string) bool {
	if valid, missingField := form.Required(c.R, fields...); !valid {
		c.Sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Warning})
		c.Sess.Save(c.R, c.W)
		return false
	}

	return true
}

// Repopulate fills the forms on the page after the user submits
func (c *Info) Repopulate(v map[string]interface{}, fields ...string) {
	form.Repopulate(c.R.Form, v, fields...)
}

// FlashSuccess saves a success flash
func (c *Info) FlashSuccess(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Success})
	c.Sess.Save(c.R, c.W)
}

// FlashNotice saves a notice flash
func (c *Info) FlashNotice(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Notice})
	c.Sess.Save(c.R, c.W)
}

// FlashWarning saves a warning flash
func (c *Info) FlashWarning(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Warning})
	c.Sess.Save(c.R, c.W)
}

//FlashError saves an error flash and logs the error
func (c *Info) FlashError(err error) {
	log.Println(err)
	c.Sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
	c.Sess.Save(c.R, c.W)
}
