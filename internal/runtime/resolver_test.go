package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_UsesAURELIA_HOME(t *testing.T) {
	want := t.TempDir()
	t.Setenv("AURELIA_HOME", want)

	r, err := New()
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if r.Root() != want {
		t.Errorf("Root() = %q, want %q", r.Root(), want)
	}
}

func TestNew_DefaultsToUserHome(t *testing.T) {
	t.Setenv("AURELIA_HOME", "") // ensure env var is not set

	r, err := New()
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".aurelia")
	if r.Root() != want {
		t.Errorf("Root() = %q, want %q", r.Root(), want)
	}
}

func TestPathResolver_Root(t *testing.T) {
	base := t.TempDir()
	r := &PathResolver{root: base}

	if r.Root() != base {
		t.Errorf("Root() = %q, want %q", r.Root(), base)
	}
}

func TestPathResolver_Accessors(t *testing.T) {
	base := "/tmp/testinstance"
	r := &PathResolver{root: base}

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"Config", r.Config(), filepath.Join(base, "config")},
		{"Data", r.Data(), filepath.Join(base, "data")},
		{"Memory", r.Memory(), filepath.Join(base, "memory")},
		{"MemoryPersonas", r.MemoryPersonas(), filepath.Join(base, "memory", "personas")},
		{"Skills", r.Skills(), filepath.Join(base, "skills")},
		{"Logs", r.Logs(), filepath.Join(base, "logs")},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestProjectHelpers(t *testing.T) {
	base := t.TempDir()

	if got := ProjectRoot(base); got != filepath.Join(base, ".aurelia") {
		t.Fatalf("ProjectRoot() = %q", got)
	}
	if got := ProjectSkills(base); got != filepath.Join(base, ".aurelia", "skills") {
		t.Fatalf("ProjectSkills() = %q", got)
	}
}


