package task

import (
	"fmt"
	"slices"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func validateRecurrence(r taskdomain.Recurrence) (taskdomain.Recurrence, error) {
	if r.Type == "" {
		r.Type = taskdomain.RecurrenceNone
	}

	switch r.Type {
	case taskdomain.RecurrenceNone:
		return taskdomain.Recurrence{Type: taskdomain.RecurrenceNone}, nil
	case taskdomain.RecurrenceDaily:
		if r.EveryNDays <= 0 {
			return taskdomain.Recurrence{}, fmt.Errorf("%w: every_n_days must be greater than 0", ErrInvalidInput)
		}
		return taskdomain.Recurrence{
			Type:       taskdomain.RecurrenceDaily,
			EveryNDays: r.EveryNDays,
		}, nil
	case taskdomain.RecurrenceMonthly:
		if r.DayOfMonth < 1 || r.DayOfMonth > 30 {
			return taskdomain.Recurrence{}, fmt.Errorf("%w: day_of_month must be in range [1..30]", ErrInvalidInput)
		}
		return taskdomain.Recurrence{
			Type:       taskdomain.RecurrenceMonthly,
			DayOfMonth: r.DayOfMonth,
		}, nil
	case taskdomain.RecurrenceSpecificDates:
		if len(r.Dates) == 0 {
			return taskdomain.Recurrence{}, fmt.Errorf("%w: dates must be set for specific_dates", ErrInvalidInput)
		}
		cleanDates := make([]string, 0, len(r.Dates))
		seen := map[string]struct{}{}
		for i := range r.Dates {
			date, err := parseDate(r.Dates[i])
			if err != nil {
				return taskdomain.Recurrence{}, fmt.Errorf("%w: invalid date format in dates (expected YYYY-MM-DD)", ErrInvalidInput)
			}
			formatted := date.Format("2006-01-02")
			if _, exists := seen[formatted]; exists {
				continue
			}
			seen[formatted] = struct{}{}
			cleanDates = append(cleanDates, formatted)
		}
		slices.Sort(cleanDates)
		return taskdomain.Recurrence{
			Type:  taskdomain.RecurrenceSpecificDates,
			Dates: cleanDates,
		}, nil
	case taskdomain.RecurrenceEvenOdd:
		if r.Parity != taskdomain.DayParityEven && r.Parity != taskdomain.DayParityOdd {
			return taskdomain.Recurrence{}, fmt.Errorf("%w: parity must be even or odd", ErrInvalidInput)
		}
		return taskdomain.Recurrence{
			Type:   taskdomain.RecurrenceEvenOdd,
			Parity: r.Parity,
		}, nil
	default:
		return taskdomain.Recurrence{}, fmt.Errorf("%w: unknown recurrence type", ErrInvalidInput)
	}
}

func matchesDate(task taskdomain.Task, target time.Time) bool {
	switch task.Recurrence.Type {
	case "", taskdomain.RecurrenceNone:
		return sameDate(task.CreatedAt.UTC(), target)
	case taskdomain.RecurrenceDaily:
		if task.Recurrence.EveryNDays <= 0 {
			return false
		}
		start := dayStart(task.CreatedAt.UTC())
		day := dayStart(target)
		if day.Before(start) {
			return false
		}
		diffDays := int(day.Sub(start).Hours() / 24)
		return diffDays%task.Recurrence.EveryNDays == 0
	case taskdomain.RecurrenceMonthly:
		return target.Day() == task.Recurrence.DayOfMonth
	case taskdomain.RecurrenceSpecificDates:
		targetStr := target.Format("2006-01-02")
		for _, d := range task.Recurrence.Dates {
			if d == targetStr {
				return true
			}
		}
		return false
	case taskdomain.RecurrenceEvenOdd:
		if task.Recurrence.Parity == taskdomain.DayParityEven {
			return target.Day()%2 == 0
		}
		return target.Day()%2 != 0
	default:
		return false
	}
}

func parseDate(raw string) (time.Time, error) {
	return time.Parse("2006-01-02", raw)
}

func dayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func sameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
