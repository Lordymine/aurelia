package session

import "testing"

func TestStore_SetGetClear(t *testing.T) {
	s := NewStore()
	if id := s.Get(1); id != "" {
		t.Fatalf("expected empty, got %q", id)
	}
	s.Set(1, "sess-abc")
	if id := s.Get(1); id != "sess-abc" {
		t.Fatalf("expected sess-abc, got %q", id)
	}
	sid, active := s.GetWithState(1)
	if sid != "sess-abc" || !active {
		t.Fatalf("expected sess-abc/active, got %q/%v", sid, active)
	}
	s.Clear(1)
	if id := s.Get(1); id != "" {
		t.Fatalf("expected empty after clear, got %q", id)
	}
}

func TestStore_DeactivateAll(t *testing.T) {
	s := NewStore()
	s.Set(1, "sess-a")
	s.Set(2, "sess-b")
	s.Set(3, "sess-c")

	// All active before deactivation
	for _, chatID := range []int64{1, 2, 3} {
		if _, active := s.GetWithState(chatID); !active {
			t.Fatalf("chat %d should be active before DeactivateAll", chatID)
		}
	}

	s.DeactivateAll()

	// All inactive after deactivation, but IDs preserved
	for _, chatID := range []int64{1, 2, 3} {
		sid, active := s.GetWithState(chatID)
		if active {
			t.Fatalf("chat %d should be inactive after DeactivateAll", chatID)
		}
		if sid == "" {
			t.Fatalf("chat %d session ID should be preserved after DeactivateAll", chatID)
		}
	}

	// Get still returns the session ID
	if id := s.Get(1); id != "sess-a" {
		t.Fatalf("Get(1) = %q, want %q", id, "sess-a")
	}
}

func TestStore_DeactivateAll_Empty(t *testing.T) {
	s := NewStore()
	s.DeactivateAll() // should not panic
}

func TestStore_Cwd(t *testing.T) {
	s := NewStore()
	s.SetCwd(1, "/home/user")
	if cwd := s.GetCwd(1); cwd != "/home/user" {
		t.Fatalf("expected /home/user, got %q", cwd)
	}
	s.Clear(1)
	if cwd := s.GetCwd(1); cwd != "" {
		t.Fatalf("expected empty after clear, got %q", cwd)
	}
}
