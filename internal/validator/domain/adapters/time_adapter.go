package adapters

import "time"

type TimeAdapter interface {
	Now() time.Time
}
