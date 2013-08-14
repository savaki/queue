package queue

import (
	"strconv"
	"testing"
	"time"
)

func TestAssembleDeleteMessageBatch(t *testing.T) {
	// GIVEN
	del := make(chan string)
	count := 7
	go func() {
		for i := 0; i < count; i++ {
			del <- strconv.Itoa(i)
		}
	}()

	// WHEN
	batch := assembleDeleteMessageBatch(del, 100*time.Millisecond)
	if len(batch) != count {
		t.Fatalf("expected %d elements in our batch; actual was %d", count, len(batch))
	}

	for i := 0; i < count; i++ {
		if batch[i].Id != strconv.Itoa(i+1) {
			t.Fatal("expected id to be the index + 1")
		}

		if batch[i].ReceiptHandle != strconv.Itoa(i) {
			t.Fatal("expected the handle to be the index value (0 indexed)")
		}
	}
}
