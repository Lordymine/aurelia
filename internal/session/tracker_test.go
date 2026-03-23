package session

import "testing"

func TestTracker_AddAndGet(t *testing.T) {
	tr := NewTracker()
	total := tr.Add(1, 1000, 500, 2, 0.05)
	if total != 1500 {
		t.Fatalf("expected 1500 total, got %d", total)
	}
	usage := tr.Get(1)
	if usage.InputTokens != 1000 || usage.OutputTokens != 500 || usage.NumTurns != 2 {
		t.Fatalf("unexpected usage: %+v", usage)
	}
	total = tr.Add(1, 500, 200, 1, 0.02)
	if total != 2200 {
		t.Fatalf("expected 2200 total, got %d", total)
	}
}

func TestTracker_RecordUsage(t *testing.T) {
	tr := NewTracker()
	shouldReset := tr.RecordUsage(1, 5, 0.10, 100000)
	if shouldReset {
		t.Fatal("should not reset below threshold")
	}
	tr.Add(1, 100000, 0, 0, 0)
	shouldReset = tr.RecordUsage(1, 1, 0.01, 100000)
	if !shouldReset {
		t.Fatal("should reset when over threshold")
	}
}

func TestTracker_Clear(t *testing.T) {
	tr := NewTracker()
	tr.Add(1, 1000, 0, 1, 0.01)
	tr.Clear(1)
	usage := tr.Get(1)
	if usage.TotalTokens() != 0 {
		t.Fatalf("expected zero after clear, got %d", usage.TotalTokens())
	}
}

func TestTracker_RecordUsageZeroThreshold(t *testing.T) {
	tr := NewTracker()
	shouldReset := tr.RecordUsage(1, 10, 0.50, 0)
	if shouldReset {
		t.Fatal("should not reset when maxTokens is 0 (disabled)")
	}
}
