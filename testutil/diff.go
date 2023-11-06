package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertObject(t *testing.T, x any, y any) {
	t.Helper()

	if diff := cmp.Diff(x, y); len(diff) > 0 {
		t.Fatalf("diff: %s\n", diff)
	}
}
