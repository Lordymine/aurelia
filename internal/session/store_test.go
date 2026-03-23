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
