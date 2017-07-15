package rollback

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
		cherr <- err
		return
	}

	cherr <- nil

	select {
	case <-chFinish:
		return
	case <-chRollback:
		w.Rollback()
		return
	}
}
