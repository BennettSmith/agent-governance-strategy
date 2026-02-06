package builder

import (
	"context"
	"testing"
)

func TestBuild_RequiresOutDir(t *testing.T) {
	_, err := Build(context.Background(), BuildOptions{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

