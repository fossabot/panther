package delivery

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestGetSQSClient(t *testing.T) {
	assert.NotNil(t, getSQSClient())
}

// 95 ms / op
func BenchmarkSessionCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		session.Must(session.NewSession())
	}
}

// 2.7 ms / op
func BenchmarkClientCreation(b *testing.B) {
	sess := session.Must(session.NewSession())
	for i := 0; i < b.N; i++ {
		sqs.New(sess)
	}
}
