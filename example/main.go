package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rkrmr33/checklist"
)

func main() {
	w := checklist.NewTerminalWriter(os.Stdout)

	items := []fakeItem{
		{
			name:         "item1",
			info:         "some fake item",
			readyAfter:   4 * time.Second,
			state:        checklist.Waiting,
			finalState:   checklist.Ready,
			finalMessage: "task completed",
		},
		{
			name:         "item2",
			info:         "another fake item",
			readyAfter:   3 * time.Second,
			state:        checklist.Waiting,
			finalState:   checklist.Ready,
			finalMessage: "task completed",
		},
		{
			name:         "item3",
			info:         "last fake item",
			readyAfter:   7 * time.Second,
			state:        checklist.Waiting,
			finalState:   checklist.Error,
			finalMessage: "failed to complete task",
		},
	}

	cl := checklist.NewCheckList(
		w,
		[]string{"NAME", "INFO", "MESSAGE"},
		[]checklist.Checker{
			checkerForItem(items[0]),
			checkerForItem(items[1]),
			checkerForItem(items[2]),
		},
		&checklist.CheckListOptions{ClearAfter: true},
	)

	fmt.Println("Starting:")

	cl.Start(context.Background())

	fmt.Println("Finished!")
}

type (
	fakeItem struct {
		name         string
		info         string
		readyAfter   time.Duration
		state        checklist.ListItemState
		message      string
		finalState   checklist.ListItemState
		finalMessage string
	}
)

func checkerForItem(i fakeItem) checklist.Checker {
	go func() {
		<-time.After(i.readyAfter)
		i.state = i.finalState
		i.message = i.finalMessage
	}()

	return func(ctx context.Context) (checklist.ListItemState, checklist.ListItemInfo) {
		// simulate time to calculate state
		<-time.After(100 * time.Millisecond)

		return i.state, []string{i.name, i.info, i.message}
	}
}
