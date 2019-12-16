package apihandlers

import (
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/internal/compliance/remediation_api/remediation"
)

type mockInvoker struct {
	remediation.InvokerAPI
	mock.Mock
}

func (m *mockInvoker) Remediate(input *models.RemediateResource) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockInvoker) GetRemediations() (*models.Remediations, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Remediations), args.Error(1)
}
