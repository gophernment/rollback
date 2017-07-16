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

	for i := 0; i < total; i++ {
		if err := <-chErr; err != nil {
			gerr = err
		}
	}

	if gerr != nil {
		close(chRollback)
		return gerr
	}

	close(chFinish)
	return nil
}
