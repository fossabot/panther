package classification

import (
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/awslogs"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/osquerylogs"
)

// A slice containing all the available parsers
var availableParsers = []parsers.LogParser{
	&awslogs.CloudTrailParser{},
	&awslogs.S3ServerAccessParser{},
	&awslogs.VPCFlowParser{},
	&awslogs.ApplicationLoadBalancerParser{},
	&awslogs.AuroraMySQLAuditParser{},
	&osquerylogs.DifferentialParser{},
	&osquerylogs.BatchParser{},
	&osquerylogs.StatusParser{},
	&osquerylogs.SnapshotParser{},
}

// ParserPriorityQueue contains parsers in priority order
type ParserPriorityQueue struct {
	items []*ParserQueueItem
}

// initialize adds all parsers to the priority queue
// All parsers have the same priority
func (q *ParserPriorityQueue) initialize() {
	for _, parser := range availableParsers {
		q.items = append(q.items, &ParserQueueItem{
			parser:  parser,
			penalty: 1,
		})
	}
}

// ParserQueueItem contains all the information needed to initialize a schema.
type ParserQueueItem struct {
	parser parsers.LogParser
	// The smaller the number the higher the priority of the parser in the queue
	penalty int
}

// Len returns the length of the priority queue
func (q *ParserPriorityQueue) Len() int {
	return len(q.items)
}

// Less compares two items of the priority queue
func (q *ParserPriorityQueue) Less(i, j int) bool {
	return q.items[i].penalty < q.items[j].penalty
}

// Swap swaps two items in the priority queue
func (q *ParserPriorityQueue) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

// Push adds an element to the end of the SchemaQueue
func (q *ParserPriorityQueue) Push(x interface{}) {
	q.items = append(q.items, x.(*ParserQueueItem))
}

// Pop removes the last element of the queue
func (q *ParserPriorityQueue) Pop() interface{} {
	n := len(q.items)
	item := q.items[n-1]
	q.items[n-1] = nil // avoid memory leak
	q.items = q.items[0 : n-1]
	return item
}

// Peek returns the item with the higher priority without removing it
func (q *ParserPriorityQueue) Peek() *ParserQueueItem {
	return q.items[0]
}
