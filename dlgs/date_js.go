//go:build js
// +build js

package dlgs

import (
	"time"
)

// Date displays a calendar dialog, returning the date and a bool for success.
func Date(title, text string, defaultDate time.Time) (time.Time, bool, error) {
	return time.Time{}, false, ErrNotImplemented
}
