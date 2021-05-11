package date

import "time"

// Format time
func Format(d time.Time, format string) string {
	return d.Format(format)
}
