package async

import (
	"testing"
)

func TestSyncParallelJobs(t *testing.T) {
	err := SyncParallel(&fakeCreateAcct{}, &fakeCreateWorker{})
	if err == nil {
		t.Error("it should error but not")
	}
}
