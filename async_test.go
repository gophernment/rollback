package rollback

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

func fakeCreateEmpFunc(email string) error {
	return nil
}

type fakeCreateAcct struct {
	email, pass string
	err         error
	chErr       chan error
	rollback    bool
}

func (h *fakeCreateAcct) Do() error {
	return fakeCreateAcctFunc(h.email, h.pass)
}

func (h *fakeCreateAcct) Rollback() {
	h.rollback = true
}

type fakeCreateWorker struct {
	email    string
	err      error
	chErr    chan error
	rollback bool
}

func (h *fakeCreateWorker) Do() error {
	return fakeCreateWorkerErrorFunc(h.email)
}

func (h *fakeCreateWorker) Rollback() {
	h.rollback = true
}

type fakeCreateEmp struct {
	email    string
	err      error
	chErr    chan error
	rollback bool
}

func (h *fakeCreateEmp) Do() error {
	return fakeCreateEmpFunc(h.email)
}

func (h *fakeCreateEmp) Rollback() {
	h.rollback = true
}

func TestAsyncHandler(t *testing.T) {
	chErr := make(chan error)
	chDone := make(chan struct{})
	chFinish := make(chan struct{})
	chRollback := make(chan struct{})

	go AsyncHandler(chErr, chDone, chFinish, chRollback, &fakeCreateAcct{})

	err := <-chErr
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncHandlerError(t *testing.T) {
	chErr := make(chan error)
	chDone := make(chan struct{})
	chFinish := make(chan struct{})
	chRollback := make(chan struct{})

	go AsyncHandler(chErr, chDone, chFinish, chRollback, &fakeCreateWorker{})

	err := <-chErr
	if err == nil {
		t.Error("it should error but not")
	}
}
