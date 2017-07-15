package async

import (
	"errors"
	"testing"
)

func fakeCreateAcctFunc(email, pass string) error {
	return nil
}

func fakeCreateWorkerErrorFunc(email string) error {
	return errors.New("worker error")
}

type fakeCreateAcct struct {
	email, pass string
	err         error
	chErr       chan error
}

func (h *fakeCreateAcct) Do() error {
	return fakeCreateAcctFunc(h.email, h.pass)
}

func (h *fakeCreateAcct) Rollback() {

}
func (h *fakeCreateAcct) Err() <-chan error {
	return h.chErr
}

type fakeCreateWorker struct {
	email string
	err   error
	chErr chan error
}

func (h *fakeCreateWorker) Do() error {
	return fakeCreateWorkerErrorFunc(h.email)
}

func (h *fakeCreateWorker) Rollback() {

}
func (h *fakeCreateWorker) Err() <-chan error {
	return h.chErr
}

func TestAsyncHandler(t *testing.T) {
	chErr := make(chan error)
	chFinish := make(chan struct{})
	chRollback := make(chan struct{})

	go AsyncHandler(chErr, chFinish, chRollback, &fakeCreateAcct{})

	err := <-chErr
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncHandlerError(t *testing.T) {
	chErr := make(chan error)
	chFinish := make(chan struct{})
	chRollback := make(chan struct{})

	go AsyncHandler(chErr, chFinish, chRollback, &fakeCreateWorker{})

	err := <-chErr
	if err == nil {
		t.Error("it should error but not")
	}
}
