package classification

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"container/heap"
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// ClassifierAPI is the interface for a classifier
type ClassifierAPI interface {
	// Classify attempts to classify the provided log line
	Classify(log string) *ClassifierResult
}

// ClassifierResult is the result of the ClassifierAPI#Classify method
type ClassifierResult struct {
	// Events contains the parsed events
	// If the classification process was not successful and the log is from an
	// unsupported type, this will be nil
	Events []interface{}
	// LogType is the identified type of the log
	LogType *string
}

// NewClassifier returns a new instance of a ClassifierAPI implementation
func NewClassifier() ClassifierAPI {
	parserQueue := &ParserPriorityQueue{}
	parserQueue.initialize()
	return &Classifier{parsers: parserQueue}
}

// Classifier is the struct responsible for classifying logs
type Classifier struct {
	parsers *ParserPriorityQueue
}

// catch panics from parsers, log and continue
func safeLogParse(parser parsers.LogParser, log string) (parsedEvents []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("parser panic",
				zap.String("parser", parser.LogType()),
				zap.Error(fmt.Errorf("%v", r)),
				zap.String("stacktrace", string(debug.Stack())),
				zap.String("log", log))
			parsedEvents = nil // return indicator that parse failed
		}
	}()
	parsedEvents = parser.Parse(log)
	return parsedEvents
}

// Classify attempts to classify the provided log line
func (c *Classifier) Classify(log string) *ClassifierResult {
	// Slice containing the popped queue items
	var popped []interface{}
	result := &ClassifierResult{}

	for c.parsers.Len() > 0 {
		currentItem := c.parsers.Peek()
		parsedEvents := safeLogParse(currentItem.parser, log)

		// Parser failed to parse event
		if parsedEvents == nil {
			zap.L().Debug("failed to parse event", zap.String("expectedLogType", currentItem.parser.LogType()))
			// Removing parser from queue
			popped = append(popped, heap.Pop(c.parsers))
			// Increasing penalty of the parser
			// Due to increased penalty the parser will be lower priority in the queue
			currentItem.penalty++
			continue
		}

		// Since the parsing was successful, remove all penalty from the parser
		// The parser will be higher priority in the queue
		currentItem.penalty = 0
		result.LogType = aws.String(currentItem.parser.LogType())
		result.Events = parsedEvents
		break
	}

	// Put back the popped items to the ParserPriorityQueue.
	for _, item := range popped {
		heap.Push(c.parsers, item)
	}
	return result
}
