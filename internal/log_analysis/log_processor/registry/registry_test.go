package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanic(t *testing.T) {
	assert.Panics(t, func() { AvailableParsers().LookupParser("doesnotexist") }, "Failed to panic, this is very dangerous!")
}
