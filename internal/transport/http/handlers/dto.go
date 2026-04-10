package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  recurrenceDTO     `json:"recurrence"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  recurrenceDTO     `json:"recurrence"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type recurrenceDTO struct {
	Type       taskdomain.RecurrenceType `json:"type"`
	EveryNDays int                       `json:"every_n_days,omitempty"`
	DayOfMonth int                       `json:"day_of_month,omitempty"`
	Dates      []string                  `json:"dates,omitempty"`
	Parity     taskdomain.DayParity      `json:"parity,omitempty"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Recurrence:  newRecurrenceDTO(task.Recurrence),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func newRecurrenceDTO(recurrence taskdomain.Recurrence) recurrenceDTO {
	return recurrenceDTO{
		Type:       recurrence.Type,
		EveryNDays: recurrence.EveryNDays,
		DayOfMonth: recurrence.DayOfMonth,
		Dates:      recurrence.Dates,
		Parity:     recurrence.Parity,
	}
}

func (r recurrenceDTO) toDomain() taskdomain.Recurrence {
	return taskdomain.Recurrence{
		Type:       r.Type,
		EveryNDays: r.EveryNDays,
		DayOfMonth: r.DayOfMonth,
		Dates:      r.Dates,
		Parity:     r.Parity,
	}
}
