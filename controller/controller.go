package controller

import (
	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/email"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/generate"
	"github.com/blue-jay/core/server"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/storage/driver/mysql"
	"github.com/blue-jay/core/view"

	"github.com/jmoiron/sqlx"
)

// Service represents all the services that the application uses.
type Service struct {
	Asset      asset.Info
	Email      email.Info
	Form       form.Info
	Generation generate.Info
	MySQL      mysql.Info
	Server     server.Info
	DB         *sqlx.DB
	Router     IRouterService
	Sess       *session.Info
	Template   view.Template
	View       IViewService
}
