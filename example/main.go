package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rkrmr33/checklist"
)

func main() {
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
			finalState:   checklist.Error,
			errMessage:   "some error aef aejfj aefj aejf jae fja efja ejf aejf aejf aje fja efja efj aejfaje fjae fja efj aejf ajef jae fjae fja efja jf aejf aejf aje fja efja efj aejfaje fjae fja efj aejf ajef jae fjae fja efja jf aejf aejf aje fja efja efj aejfaje fjae fja efj aejf ajef jae fjae fja efja ejf aejf ajef aje fjae fjae fj aejf aejf aje fjae f",
			finalMessage: "failed to complete task",
		},
		{
			name:         "item3",
			info:         "last fake item",
			readyAfter:   7 * time.Second,
			state:        checklist.Waiting,
			finalState:   checklist.Ready,
			finalMessage: "task completed",
		},
	}

	fmt.Println("Starting:")

	for i := 3; i > 0; i, _ = i-1, <-time.After(time.Millisecond*500) {
		fmt.Println(i)
	}

	cl := checklist.NewCheckList(
		os.Stdout,
		[]string{"NAME", "INFO", "MESSAGE", "ERROR"},
		[]checklist.Checker{
			checkerForItem(items[0]),
			checkerForItem(items[1]),
			checkerForItem(items[2]),
		},
		&checklist.CheckListOptions{
			ClearAfter:   false,
			Fullscreen:   false,
			WaitAllReady: false,
		},
	)

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
		err          string
		errMessage   string
		finalState   checklist.ListItemState
		finalMessage string
	}
)

func checkerForItem(i fakeItem) checklist.Checker {
	go func() {
		<-time.After(i.readyAfter)
		i.state = i.finalState
		i.message = i.finalMessage
		i.err = i.errMessage
	}()

	return func(ctx context.Context) (checklist.ListItemState, checklist.ListItemInfo) {
		// simulate time to calculate state
		<-time.After(100 * time.Millisecond)

		return i.state, []string{i.name, i.info, i.message, i.err}
	}
}
