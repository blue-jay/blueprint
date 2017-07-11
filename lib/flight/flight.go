// Package flight provides access to the application settings safely.
package flight

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blue-jay/blueprint/lib/env"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/view"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

var (
	configInfo env.Info
	dbInfo     *sqlx.DB

	mutex sync.RWMutex
)

// StoreConfig stores the application settings so controller functions can
//access them safely.
func StoreConfig(ci env.Info) {
	mutex.Lock()
	configInfo = ci
	mutex.Unlock()
}

// StoreDB stores the database connection settings so controller functions can
// access them safely.
func StoreDB(db *sqlx.DB) {
	mutex.Lock()
	dbInfo = db
	mutex.Unlock()
}

// Info structures the application settings.
type Info struct {
	Config env.Info
	Sess   *sessions.Session
	UserID string
	W      http.ResponseWriter
	R      *http.Request
	View   view.Info
	DB     *sqlx.DB
}

// Context returns the application settings.
func Context(w http.ResponseWriter, r *http.Request) Info {
	var id string

	// Get the session
	sess, err := configInfo.Session.Instance(r)

	// If the session is valid
	if err == nil {
		// Get the user id
		id = fmt.Sprintf("%v", sess.Values["id"])
	}

	mutex.RLock()
	i := Info{
		Config: configInfo,
		Sess:   sess,
		UserID: id,
		W:      w,
		R:      r,
		View:   configInfo.View,
		DB:     dbInfo,
	}
	mutex.RUnlock()

	return i
}

// Reset will delete all package globals
func Reset() {
	mutex.Lock()
	configInfo = env.Info{}
	dbInfo = &sqlx.DB{}
	mutex.Unlock()
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

// FlashError saves an error flash and logs the error.
func (c *Info) FlashError(err error) {
	log.Println(err)
	c.Sess.AddFlash(flash.Info{err.Error(), flash.Error})
	c.Sess.Save(c.R, c.W)
}

// FlashErrorGeneric saves a generic error flash and logs the error.
func (c *Info) FlashErrorGeneric(err error) {
	log.Println(err)
	c.Sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
	c.Sess.Save(c.R, c.W)
}
