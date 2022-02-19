package checklist

import (
	"fmt"
	"strings"

	"github.com/juju/ansiterm"
)

func allReady(s []ListItemState) bool {
	if len(s) == 0 {
		return false
	}

	for _, state := range s {
		if !state.isFinal() {
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
