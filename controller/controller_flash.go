package controller

import (
	"log"
	"net/http"

	"github.com/blue-jay/core/flash"
)

// FlashSuccess saves a success flash.
func (c Service) FlashSuccess(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Info{message, flash.Success})
	sess.Save(r, w)
}

// FlashNotice saves a notice flash.
func (c Service) FlashNotice(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Info{message, flash.Notice})
	sess.Save(r, w)
}

// FlashWarning saves a warning flash.
func (c Service) FlashWarning(w http.ResponseWriter, r *http.Request, message string) {
	sess, _ := c.Sess.Instance(r)
	sess.AddFlash(flash.Info{message, flash.Warning})
	sess.Save(r, w)
}

// FlashError saves an error flash and logs the error.
func (c Service) FlashError(w http.ResponseWriter, r *http.Request, err error) {
	sess, _ := c.Sess.Instance(r)
	log.Println(err)
	sess.AddFlash(flash.Info{err.Error(), flash.Error})
	sess.Save(r, w)
}

// FlashErrorGeneric saves a generic error flash and logs the error.
func (c Service) FlashErrorGeneric(w http.ResponseWriter, r *http.Request, err error) {
	sess, _ := c.Sess.Instance(r)
	log.Println(err)
	sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
	sess.Save(r, w)
}
