// Package flight provides an abstraction around commonly used features.
package flight

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/sessions"
)

var (
	assetInfo *asset.Info
	formInfo  *form.Info
	viewInfo  *view.Info
	xsrfInfo  *xsrf.Info
	dbInfo    *sqlx.DB

	mutex sync.RWMutex
)

// StoreConfig safely stores the variables.
func StoreConfig(ai *asset.Info,
	fi *form.Info,
	vi *view.Info,
	xi *xsrf.Info,
	db *sqlx.DB) {
	mutex.Lock()

	assetInfo = ai
	formInfo = fi
	viewInfo = vi
	xsrfInfo = xi
	dbInfo = db

	mutex.Unlock()
}

// Xsrf returns the xsrf configuration.
func Xsrf() *xsrf.Info {
	mutex.RLock()
	x := xsrfInfo
	mutex.RUnlock()
	return x
}

// Info holds the commonly used information.
type Info struct {
	Asset  *asset.Info
	Form   *form.Info
	Sess   *sessions.Session
	UserID string
	W      http.ResponseWriter
	R      *http.Request
	View   *view.Info
	DB     *sqlx.DB
}

// Context returns commonly used information.
func Context(w http.ResponseWriter, r *http.Request) *Info {
	sess, _ := session.Instance(r)

	mutex.RLock()
	i := &Info{
		Asset:  assetInfo,
		Form:   formInfo,
		Sess:   sess,
		UserID: fmt.Sprintf("%v", sess.Values["id"]),
		W:      w,
		R:      r,
		View:   viewInfo,
		DB:     dbInfo,
	}
	mutex.RUnlock()

	return i
}

// Param gets the URL parameter.
func (c *Info) Param(name string) string {
	return router.Param(c.R, name)
}

// Redirect sends a temporary redirect.
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

// Repopulate fills the forms on the page after the user submits.
func (c *Info) Repopulate(v map[string]interface{}, fields ...string) {
	form.Repopulate(c.R.Form, v, fields...)
}

// FlashSuccess saves a success flash.
func (c *Info) FlashSuccess(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Success})
	c.Sess.Save(c.R, c.W)
}

// FlashNotice saves a notice flash.
func (c *Info) FlashNotice(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Notice})
	c.Sess.Save(c.R, c.W)
}

// FlashWarning saves a warning flash.
func (c *Info) FlashWarning(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Warning})
	c.Sess.Save(c.R, c.W)
}

//FlashError saves an error flash and logs the error.
func (c *Info) FlashError(err error) {
	log.Println(err)
	c.Sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
	c.Sess.Save(c.R, c.W)
}
