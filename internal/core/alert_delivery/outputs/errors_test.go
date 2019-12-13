package outputs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertDeliveryError(t *testing.T) {
	err := &AlertDeliveryError{Message: "http error"}
	assert.Equal(t, err.Message, err.Error())
}
