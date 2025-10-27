package models

import "time"

type AppealStatus string

const (
	StatusNew        AppealStatus = "New"
	StatusInProgress AppealStatus = "InProgress"
	StatusCompleted  AppealStatus = "Completed"
	StatusCancelled  AppealStatus = "Cancelled"
)

type Appeal struct {
	ID           string       `json:"id"`
	Theme        string       `json:"theme"`
	Message      string       `json:"message"`
	Status       AppealStatus `json:"status"`
	Solution     string       `json:"solution,omitempty"`
	CanselReason string       `json:"cansel_reason,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateAppealRequest struct {
	Theme   string `json:"theme" validate:"required,min=1"`
	Message string `json:"message" validate:"required,min=1"`
}

type UpdateAppealSolutionRequest struct {
	Solution string `json:"solution" validate:"required,min=1"`
}

type UpdateAppealCancelRequest struct {
	Reason string `json:"reason" validate:"required,min=1"`
}

type FilterDatesRequest struct {
	Date      string `json:"date,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
}

func (a *Appeal) IsNew() bool {
	return a.Status == StatusNew
}

func (a *Appeal) IsInProgress() bool {
	return a.Status == StatusInProgress
}

func (a *Appeal) IsCompleted() bool {
	return a.Status == StatusCompleted
}

func (a *Appeal) IsCancelled() bool {
	return a.Status == StatusCancelled
}

func (a *Appeal) CanStartProcessing() bool {
	return a.Status == StatusNew || a.Status == StatusCancelled
}

func (a *Appeal) CanComplete() bool {
	return a.Status == StatusInProgress
}

func (a *Appeal) CanCancel() bool {
	return a.Status == StatusNew || a.Status == StatusInProgress
}