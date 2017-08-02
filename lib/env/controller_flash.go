package env

import (
	"log"
	"net/http"

	"github.com/blue-jay/core/flash"
)

// FlashSuccess saves a success flash.
func (c Service) FlashSuccess(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Success(message))
	sess.Save(r, w)
}

// FlashNotice saves a notice flash.
func (c Service) FlashNotice(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Notice(message))
	sess.Save(r, w)
}

// FlashWarning saves a warning flash.
func (c Service) FlashWarning(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Warning(message))
	sess.Save(r, w)
}

// FlashError saves an error flash and logs the error.
func (c Service) FlashError(w http.ResponseWriter, r *http.Request, err error) {
	sess, _ := c.Sess.Instance(r)
	log.Println(err)
	sess.AddFlash(flash.Danger(err.Error()))
	sess.Save(r, w)
}

// FlashErrorGeneric saves a generic error flash and logs the error.
func (c Service) FlashErrorGeneric(w http.ResponseWriter, r *http.Request, err error) {
	sess, _ := c.Sess.Instance(r)
	log.Println(err)
	sess.AddFlash(flash.Danger("An error occurred on the server. Please try again later."))
	sess.Save(r, w)
}
