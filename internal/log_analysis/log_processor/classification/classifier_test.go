package classification

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

type mockParser struct {
	parsers.LogParser
	mock.Mock
}

func (m *mockParser) Parse(log string) []interface{} {
	args := m.Called(log)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.([]interface{})
}

func (m *mockParser) LogType() string {
	args := m.Called()
	return args.String(0)
}

func TestClassifyRespectsPriorityOfParsers(t *testing.T) {
	succeedingParser := &mockParser{}
	failingParser1 := &mockParser{}
	failingParser2 := &mockParser{}

	succeedingParser.On("Parse", mock.Anything).Return([]interface{}{"event"})
	succeedingParser.On("LogType").Return("success")
	failingParser1.On("Parse", mock.Anything).Return(nil)
	failingParser1.On("LogType").Return("failure1")
	failingParser2.On("Parse", mock.Anything).Return(nil)
	failingParser2.On("LogType").Return("failure2")

	availableParsers = []parsers.LogParser{failingParser1, succeedingParser, failingParser2}

	classifier := NewClassifier()

	expectedResult := &ClassifierResult{
		Events:  []interface{}{"event"},
		LogType: aws.String("success"),
	}

	repetitions := 1000
	for i := 0; i < repetitions; i++ {
		result := classifier.Classify("log")
		require.Equal(t, expectedResult, result)
	}

	succeedingParser.AssertNumberOfCalls(t, "Parse", repetitions)
	requireLessOrEqualNumberOfCalls(t, failingParser1, "Parse", 1)
	requireLessOrEqualNumberOfCalls(t, failingParser2, "Parse", 1)
}

func TestClassifyNoMatch(t *testing.T) {
	failingParser := &mockParser{}

	failingParser.On("Parse", mock.Anything).Return(nil)
	failingParser.On("LogType").Return("failure")

	availableParsers = []parsers.LogParser{failingParser}

	classifier := NewClassifier()

	result := classifier.Classify("log")
	require.Equal(t, &ClassifierResult{}, result)
	failingParser.AssertNumberOfCalls(t, "Parse", 1)
}

func requireLessOrEqualNumberOfCalls(t *testing.T, underTest *mockParser, method string, number int) {
	timesCalled := 0
	for _, call := range underTest.Calls {
		if call.Method == method {
			timesCalled++
		}
	}
	require.LessOrEqual(t, timesCalled, number)
}
