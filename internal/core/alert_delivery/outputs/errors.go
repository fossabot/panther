package outputs

// AlertDeliveryError indicates whether a failed alert should be retried.
type AlertDeliveryError struct {
	// Message is the description of the problem: what went wrong.
	Message string

	// Permanent indicates whether the alert output should be retried.
	// For example, outputs which don't exist or errors creating the request are permanent failures.
	// But any error talking to the output itself can be retried by the Lambda function later.
	Permanent bool
}

func (e *AlertDeliveryError) Error() string { return e.Message }
