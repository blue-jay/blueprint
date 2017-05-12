package env

import "context"

// User is the type of value stored in the Contexts.
type User struct {
	ID string
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

// NewUserContext returns a new Context that carries value u.
func NewUserContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserContext returns the User value stored in ctx, if any.
func UserContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
