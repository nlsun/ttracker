package format

import (
	"strings"

	"github.com/nlsun/ttracker/tracker"
)

func Title(targets []tracker.Target) string {
	names := []string{"Timestamp"}
	for _, tgt := range targets {
		names = append(names, tgt.Name)
	}
	return strings.Join(names, ",")
}

func Line(res map[tracker.Target]string, targets []tracker.Target, timestamp string) string {
	fields := []string{timestamp}
	for _, tgt := range targets {
		fields = append(fields, res[tgt])
	}
	return strings.Join(fields, ",")
}
