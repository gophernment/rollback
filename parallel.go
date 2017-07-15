package rollback

func SyncParallel(workers ...Worker) error {
	var gerr error
	total := len(workers)
	chErr := make(chan error)
	chDone := make(chan struct{})
	chFinish := make(chan struct{})
	chRollback := make(chan struct{})

	defer func() {
		for i := 0; i < total; i++ {
			<-chDone
		}
	}()

	for i := range workers {
		go AsyncHandler(chErr, chDone, chFinish, chRollback, workers[i])
	}

	errCount := 0
	for i := 0; i < total; i++ {
		if err := <-chErr; err != nil {
			gerr = err
			errCount++
		}
	}

	if errCount != 0 {
		broadcast := total - errCount
		for i := 0; i < broadcast; i++ {
			chRollback <- struct{}{}
		}
		return gerr
	}

	chFinish <- struct{}{}
	chFinish <- struct{}{}
	chFinish <- struct{}{}

	return nil
}
