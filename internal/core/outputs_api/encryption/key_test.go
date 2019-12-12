package encryption

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, "keyid", *New("keyid", session.Must(session.NewSession())).ID)
}
