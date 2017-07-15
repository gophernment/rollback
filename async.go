package rollback

import (
	"fmt"

	"github.com/pkg/errors"
)

type Worker interface {
	Do() error
	Rollback()
}

func AsyncHandler(cherr chan error, chDone, chFinish, chRollback chan struct{}, w Worker) {
	defer func() {
		chDone <- struct{}{}
	}()

	err := w.Do()
	if err != nil {
		cherr <- errors.Wrap(err, "common create user")
		return
	}

	cherr <- nil

	select {
	case <-chFinish:
		fmt.Println("create Account finish")
		return
	case <-chRollback:
		fmt.Println("create Account rollback")
		w.Rollback()
		return
	}
}
