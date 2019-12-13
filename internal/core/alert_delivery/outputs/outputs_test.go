package outputs

import "github.com/stretchr/testify/mock"

type mockHTTPWrapper struct {
	HTTPWrapper
	mock.Mock
}

func (m *mockHTTPWrapper) post(postInput *PostInput) *AlertDeliveryError {
	args := m.Called(postInput)
	return args.Get(0).(*AlertDeliveryError)
}
