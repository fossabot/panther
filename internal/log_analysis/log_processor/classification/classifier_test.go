package classification

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/registry"
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

// admit to registry.Interface interface
type TestRegistry map[string]*registry.LogParserMetadata

func NewTestRegistry() TestRegistry {
	return make(map[string]*registry.LogParserMetadata)
}

func (r TestRegistry) Add(lpm *registry.LogParserMetadata) {
	r[lpm.Parser.LogType()] = lpm
}

func (r TestRegistry) Elements() map[string]*registry.LogParserMetadata {
	return r
}

func (r TestRegistry) LookupParser(logType string) (lpm *registry.LogParserMetadata) {
	return (registry.Registry)(r).LookupParser(logType) // call registry code
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

	availableParsers := []*registry.LogParserMetadata{
		{Parser: failingParser1},
		{Parser: succeedingParser},
		{Parser: failingParser2},
	}
	testRegistry := NewTestRegistry()
	parserRegistry = testRegistry // re-bind as interface
	for i := range availableParsers {
		testRegistry.Add(availableParsers[i]) // update registry
	}

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

	availableParsers := []*registry.LogParserMetadata{
		{Parser: failingParser},
	}
	testRegistry := NewTestRegistry()
	parserRegistry = testRegistry // re-bind as interface
	for i := range availableParsers {
		testRegistry.Add(availableParsers[i]) // update registry
	}

	classifier := NewClassifier()

	result := classifier.Classify("log")
	require.Equal(t, &ClassifierResult{}, result)
	failingParser.AssertNumberOfCalls(t, "Parse", 1)
}

func TestClassifyParserPanic(t *testing.T) {
	// uncomment to see the logs produced
	/*
		logger := zap.NewExample()
		defer logger.Sync()
		undo := zap.ReplaceGlobals(logger)
		defer undo()
	*/

	panicParser := &mockParser{}

	panicParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) { panic("test parser panic") })
	panicParser.On("LogType").Return("panic parser")

	availableParsers := []*registry.LogParserMetadata{
		{Parser: panicParser},
	}
	testRegistry := NewTestRegistry()
	parserRegistry = testRegistry // re-bind as interface
	for i := range availableParsers {
		testRegistry.Add(availableParsers[i]) // update registry
	}

	classifier := NewClassifier()

	result := classifier.Classify("log of death")
	require.Equal(t, &ClassifierResult{}, result)
	panicParser.AssertNumberOfCalls(t, "Parse", 1)
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
