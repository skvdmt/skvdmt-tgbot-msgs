package pkg

import (
	"fmt"
	"strings"
	"time"
)

// PrettyTime Форматирует промежуток времени в строку с описанием.
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
		prettyTimeEdit(h, "hour"),
		prettyTimeEdit(m, "minute"),
		prettyTimeEdit(s, "second"),
	))
}

// prettyTimeEdit Редактирует описание форматирования.
func prettyTimeEdit(value time.Duration, section string) string {
	switch value {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d %s ", value, section)
	default:
		return fmt.Sprintf("%d %ss ", value, section)
	}
}
