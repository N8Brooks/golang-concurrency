// Package promise_time_limit contains the Promise Time Limit problem and shared
// timeout error contract.
package promise_time_limit

import "errors"

var ErrTimeLimitExceeded = errors.New("time limit exceeded")
