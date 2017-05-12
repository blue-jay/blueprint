package controller

import (
	"net/http"

	"github.com/blue-jay/core/view"
)

// IViewService is the interface for HTML templates.
type IViewService interface {
	Base(base string) *view.Info
	New(templateList ...string) *view.Info
	Render(w http.ResponseWriter, r *http.Request) error
}
