package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {
	var have BuildInfo
	want := BuildInfo{
		Version: DummyVersion,
		Date:    DummyDate,
		Branch:  DummyBranch,
		Commit:  DummyCommit,
	}

	Dummy(&have)
	assert.Exactly(t, want, have, "empty values should fallback to their dummy values")
}
