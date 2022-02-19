package checklist

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/juju/ansiterm"
)

// ANSI escape codes
const (
	escape    = "\x1b"
	noFormat  = 0
	bold      = 1
	fgBlack   = 30
	fgRed     = 31
	fgGreen   = 32
	fgYellow  = 33
	fgBlue    = 34
	fgMagenta = 35
	fgCyan    = 36
	fgWhite   = 37
	fgDefault = 39
	fgHiBlue  = 94
)

// Defaults
var (
	defaultInterval = 500 * time.Millisecond

	defaultStateIconMap = map[ListItemState]string{
		Waiting: "↻",
		Ready:   "✔",
		Error:   "✖",
	}

	defaultStateColorMap = map[ListItemState]int{
		Waiting: fgCyan,
		Ready:   fgGreen,
		Error:   fgRed,
	}
)

type (
	// CheckList holds the state of a check list
	CheckList struct {
		w       CleanWriter
		tb      *ansiterm.TabWriter
		headers ListItemInfo
		items   []Checker
		opts    *CheckListOptions

		curState []ListItemState
		curInfos []ListItemInfo
	}

	// CheckListOptions options to create a new check list
	CheckListOptions struct {
		// Interval to refresh the checklist default is 500ms
		Interval time.Duration
		// StateIconMap the icon to display for each item state
		StateIconMap map[ListItemState]string
		// NoColor if true will not use colors at all
		NoColor bool
		// Fullscreen if true will clear all screen each refresh
		Fullscreen bool
		// ClearAfter if true will clear the check list from the
		// screen after all the checks are done
		ClearAfter bool
	}

	// CleanWriter a writer that can be cleaned, like clearing the terminal.
	CleanWriter interface {
		io.Writer
		// Clean cleans what was printed to the writer before
		// if the number of lines is -1 it means that the
		// entire screen needs to be cleared instead of just
		// the number of lines
		Clean(linesPrinted int) error
	}

	// ListItemState holds the state of a list item
	ListItemState string
	// ListItemInfo holds the values for all the table columns for a single list item
	ListItemInfo []string
	// Checker a function used to check the current state and info of a single list item
	Checker func(ctx context.Context) (ListItemState, ListItemInfo)
)

// Item states
const (
	// Waiting for the item to be ready
	Waiting ListItemState = "Waiting"
	// Ready item is ready
	Ready ListItemState = "Ready"
	// Error item checker returns some error
	Error ListItemState = "Error"
)

// NewCheckList creates a new checklist
func NewCheckList(w CleanWriter, headers ListItemInfo, items []Checker, opts *CheckListOptions) *CheckList {
	return &CheckList{
		w:       w,
		tb:      ansiterm.NewTabWriter(w, 0, 0, 2, ' ', 0),
		headers: headers,
		items:   items,
		opts:    getOptions(opts),
	}
}

func getOptions(o *CheckListOptions) *CheckListOptions {
	opts := CheckListOptions{
		Interval:     defaultInterval,
		StateIconMap: defaultStateIconMap,
	}

	if o == nil {
		return &opts
	}

	opts.NoColor = o.NoColor
	opts.Fullscreen = o.Fullscreen
	opts.ClearAfter = o.ClearAfter
	if os.Getenv("NO_COLOR") == "true" || os.Getenv("TERM") == "dumb" {
		opts.NoColor = true
	}

	if o.Interval.Milliseconds() != 0 {
		opts.Interval = o.Interval
	}

	if o.StateIconMap != nil {
		opts.StateIconMap = o.StateIconMap
	}

	return &opts
}

// Start starts the checklist
func (cl *CheckList) Start(ctx context.Context) error {
	t := time.NewTicker(cl.opts.Interval)
	linesPrinted := 0

	if cl.opts.ClearAfter {
		defer func() {
			_ = cl.w.Clean(linesPrinted)
		}()
	}

	for !allReady(cl.curState) {
		cl.refreshItems(ctx)

		if cl.opts.Fullscreen {
			linesPrinted = -1
		}

		if err := cl.w.Clean(linesPrinted); err != nil {
			return err
		}

		if err := cl.printHeader(); err != nil {
			return err
		}

		if err := cl.printItems(); err != nil {
			return err
		}

		linesPrinted = len(cl.items) + 1

		if err := cl.tb.Flush(); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}

	return nil
}

func (cl *CheckList) printHeader() error {
	return printToTabWriter(cl.tb, append([]string{""}, cl.headers...))
}

func (cl *CheckList) printItems() error {
	lines := make([]ListItemInfo, len(cl.curState))

	for i := range lines {
		lines[i] = make([]string, 0, len(cl.curInfos[i])+1)
		lines[i] = append(lines[i], cl.getStateIcon(cl.curState[i]))
		lines[i] = append(lines[i], cl.curInfos[i]...)
	}

	for _, line := range lines {
		if err := printToTabWriter(cl.tb, line); err != nil {
			return err
		}
	}

	return nil
}

func (cl *CheckList) refreshItems(ctx context.Context) {
	n := len(cl.items)
	newStates := make([]ListItemState, n)
	newInfos := make([]ListItemInfo, n)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i, checker := range cl.items {
		go func(i int, checker Checker) {
			state, info := checker(ctx)
			newStates[i] = state
			newInfos[i] = info
			wg.Done()
		}(i, checker)
	}

	wg.Wait()

	cl.curInfos = newInfos
	cl.curState = newStates
}

func (cl *CheckList) getStateIcon(s ListItemState) string {
	icon, ok := cl.opts.StateIconMap[s]
	if !ok {
		icon = defaultStateIconMap[s]
	} else if icon != defaultStateIconMap[s] {
		return icon
	}

	return cl.colorize(icon, defaultStateColorMap[s])
}

func (cl *CheckList) colorize(s string, codes ...int) string {
	if cl.opts.NoColor {
		return s
	}

	codeStrs := make([]string, len(codes))
	for i, code := range codes {
		codeStrs[i] = strconv.Itoa(code)
	}

	sequence := strings.Join(codeStrs, ";")

	return fmt.Sprintf("%s[%sm%s%s[%dm", escape, sequence, s, escape, noFormat)
}

func (s ListItemState) isFinal() bool {
	return s != Waiting
}
