package processor

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
	"bufio"
	"sync"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/classification"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/destinations"
)

// ParsedEventBufferSize is the size of the buffer of the Go channel containing the parsed events.
// Since there are different goroutines writing and reading from that channel each with different I/O characteristics,
// we are specifying this buffer to avoid blocking the goroutines that write to the channel if the reader goroutine is
// temporarily busy. The writer goroutines will block writing but only when the buffer has been full - something we need
// to avoid using up lot of memory.
// see also: https://golang.org/doc/effective_go.html#channels
const ParsedEventBufferSize = 1000

var (
	parsedEventChannel      chan *common.ParsedEvent
	destinationErrorChannel chan error
)

// Handle orchestrates the tasks of parsing logs, classification, normalization
// and forwarding the logs to the appropriate destination
func Handle(dataStreams []*common.DataStream) error {
	zap.L().Info("handling data streams", zap.Int("numDataStreams", len(dataStreams)))
	parsedEventChannel = make(chan *common.ParsedEvent, ParsedEventBufferSize)
	destinationErrorChannel = make(chan error)

	go func() {
		destination := destinations.CreateDestination()
		err := destination.SendEvents(parsedEventChannel)
		if err != nil {
			destinationErrorChannel <- err
		}
		close(destinationErrorChannel)
	}()

	var streamProcessingWg sync.WaitGroup
	for _, dataStream := range dataStreams {
		streamProcessingWg.Add(1)
		go func(input *common.DataStream) {
			processStream(input)
			streamProcessingWg.Done()
		}(dataStream)
	}

	go func() {
		zap.L().Info("waiting for goroutines to stop reading data", zap.Int("numDataStreams", len(dataStreams)))
		// Close the channel after all goroutines have finished writing to it.
		// The Destination that is reading the channel will terminate
		// after consuming all the buffered messages
		streamProcessingWg.Wait()
		zap.L().Info("data processing goroutines finished")
		close(parsedEventChannel)
	}()

	// Blocking until the destination has finished.
	// If the destination finished successfully this will return nil
	// otherwise it will return an error and will cause Lambda invocation to fail
	err := <-destinationErrorChannel
	return err
}

//processStream loads the data from S3, parses it and writes them to the output channel
func processStream(input *common.DataStream) {
	zap.L().Info("starting to process data stream")
	classifier := classification.NewClassifier()
	logLines := 0
	successfullyClassified := 0
	classificationFailures := 0
	scanner := bufio.NewScanner(*input.Reader)
	for scanner.Scan() {
		logLines++

		classificationResult := classifier.Classify(scanner.Text())
		if classificationResult.LogType == nil {
			zap.L().Warn("failed to classify log line", zap.Int("lineNum", logLines))
			classificationFailures++
			continue
		}
		successfullyClassified++

		for _, parsedEvent := range classificationResult.Events {
			message := &common.ParsedEvent{
				Event:   parsedEvent,
				LogType: *classificationResult.LogType,
			}
			parsedEventChannel <- message
		}
	}
	zap.L().Info("finished processing stream",
		zap.Int("logLines", logLines),
		zap.Int("successfullyClassified", successfullyClassified),
		zap.Int("classificationFailures", classificationFailures))
}
