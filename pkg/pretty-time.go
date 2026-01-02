package pkg

import (
	"fmt"
	"strings"
	"time"
)

// PrettyTime good time format string
func PrettyTime(t time.Duration) string {
	d := t.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	var s time.Duration
	if m == 0 {
		s = d / time.Second
	}
	return strings.TrimSpace(fmt.Sprintf(
		"%s%s%s",
		prettyTimeSection(h, "hour"),
		prettyTimeSection(m, "minute"),
		prettyTimeSection(s, "second"),
	))
}

// prettyTimeSection section time as string
func prettyTimeSection(value time.Duration, section string) string {
	switch value {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d %s ", value, section)
	default:
		return fmt.Sprintf("%d %ss ", value, section)
	}
}
