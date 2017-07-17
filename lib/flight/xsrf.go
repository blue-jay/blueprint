package flight

import (
	"sync"

	"github.com/blue-jay/core/xsrf"
)

var (
	xsrfInfo  xsrf.Info
	xsrfMutex sync.RWMutex
)

// StoreXsrf sets the csrf configuration.
func StoreXsrf(x xsrf.Info) {
	xsrfMutex.Lock()
	xsrfInfo = x
	xsrfMutex.Unlock()
}

// Xsrf returns the csrf configuration.
func Xsrf() xsrf.Info {
	xsrfMutex.RLock()
	x := xsrfInfo
	xsrfMutex.RUnlock()
	return x
}
