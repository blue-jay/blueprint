package env

import (
	"context"
	"net/http"

	"encoding/gob"

	"github.com/blue-jay/core/session"
)

// User is the type of value stored in the Contexts.
type User struct {
	ID        string
	Email     string
	FirstName string
}

// LoggedIn returns true if the user is logged in.
func (u User) LoggedIn() bool {
	return u.ID != ""
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

func init() {
	gob.Register(key(userKey))
	gob.Register(User{})
}

// NewUserContext returns a new Context that carries value u.
func NewUserContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserContext returns the User value stored in ctx, if any.
func UserContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

// StoreUserSession stores the user in the session.
func StoreUserSession(w http.ResponseWriter, r *http.Request, info *session.Info, u *User) error {
	sess, _ := info.Instance(r)
	sess.Values[userKey] = u
	return sess.Save(r, w)
}

// UserSession returns the User value stored in session, if any.
func UserSession(r *http.Request, info *session.Info) (*User, bool) {
	sess, _ := info.Instance(r)
	u, ok := sess.Values[userKey].(User)
	return &u, ok
}
