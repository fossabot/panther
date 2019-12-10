package processor

import (
	"bufio"
	"sync"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/classification"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/destinations"
)

var (
	streamProcessingWg sync.WaitGroup
	destinationsWg     sync.WaitGroup
	parsedEventChannel chan *common.ParsedEvent
	errorChannel       chan error
)

// Handle orchestrates the tasks of parsing logs, classification, normalization
// and forwarding the logs to the appropriate destination
func Handle(dataStreams []*common.DataStream) error {
	zap.L().Info("handling data streams", zap.Int("numDataStreams", len(dataStreams)))
	// Creating a buffered output channel
	parsedEventChannel = make(chan *common.ParsedEvent, 1000)
	errorChannel = make(chan error)
	streamProcessingWg = sync.WaitGroup{}
	destinationsWg = sync.WaitGroup{}

	startDestinations()

	for _, dataStream := range dataStreams {
		streamProcessingWg.Add(1)
		go func(input *common.DataStream) {
			defer streamProcessingWg.Done()
			processStream(input)
		}(dataStream)
	}
	zap.L().Info("waiting for stream processing to finish", zap.Int("numDataStreams", len(dataStreams)))
	streamProcessingWg.Wait() //Wait for all go routines that process the data streams to stop
	// Once all have stopped, close the channel to signal to destinations that there are no more
	// parsed events left to process.
	close(parsedEventChannel)

	zap.L().Info("waiting for destination to finish")
	destinationsWg.Wait()
	close(errorChannel)

	// If any error has occurred, the error will be returned by the Lambda
	// If no error has occurred, this will be nil
	return <-errorChannel
}

func startDestinations() {
	destination := destinations.CreateDestination()
	destinationsWg.Add(1)
	go func() {
		defer destinationsWg.Done()
		destination.SendEvents(parsedEventChannel, errorChannel)
	}()
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
				LogType: classificationResult.LogType,
			}
			parsedEventChannel <- message
		}
	}
	zap.L().Info("finished processing stream",
		zap.Int("logLines", logLines),
		zap.Int("successfullyClassified", successfullyClassified),
		zap.Int("classificationFailures", classificationFailures))
}
