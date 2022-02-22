package checklist

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/ansiterm"
)

func allReady(s []ListItemState, waitAllReady bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, state := range s {
		if !state.isFinal(waitAllReady) {
			return false
		}
	}

	return true
}

func printToTabWriter(w *ansiterm.TabWriter, items []string) error {
	format := strings.Repeat("%s\t", len(items)-1)
	format += "%s\n"

	values := make([]interface{}, len(items))
	for i := range items {
		values[i] = items[i]
	}

	if _, err := fmt.Fprintf(w, format, values...); err != nil {
		return err
	}

	return nil
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripANSI(str string) string {
	return re.ReplaceAllString(str, "")
}
