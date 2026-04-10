package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	Recurrence  Recurrence `json:"recurrence"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type RecurrenceType string

const (
	RecurrenceNone          RecurrenceType = "none"
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceEvenOdd       RecurrenceType = "even_odd"
)

type DayParity string

const (
	DayParityEven DayParity = "even"
	DayParityOdd  DayParity = "odd"
)

type Recurrence struct {
	Type       RecurrenceType `json:"type"`
	EveryNDays int            `json:"every_n_days,omitempty"`
	DayOfMonth int            `json:"day_of_month,omitempty"`
	Dates      []string       `json:"dates,omitempty"`
	Parity     DayParity      `json:"parity,omitempty"`
}
