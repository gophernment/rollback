package rollback

import (
	"testing"
)

func TestSyncParallelJobs(t *testing.T) {
	err := SyncParallel(&fakeCreateAcct{}, &fakeCreateWorker{})
	if err == nil {
		t.Error("it should error but not")
	}
}

func TestSyncParallelJobsRollback(t *testing.T) {
	fa := &fakeCreateAcct{}
	fw := &fakeCreateWorker{}
	fe := &fakeCreateEmp{}
	SyncParallel(fa, fw, fe)

	if !fa.rollback {
		t.Error("create account should rollback")
	}

	if !fe.rollback {
		t.Error("create employee should rollback")
	}

	if fw.rollback {
		t.Error("create worker should not rollback")
	}
}
