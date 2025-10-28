package models

import (
	"testing"
	"time"
)

func TestAppealStatusMethods(t *testing.T) {
	t.Parallel()

	appeal := &Appeal{
		ID:        "test-id",
		Theme:     "Test Theme",
		Message:   "Test Message",
		Status:    StatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("IsNew", func(t *testing.T) {
		appeal.Status = StatusNew
		if !appeal.IsNew() {
			t.Errorf("Expected appeal to be new, got %s", appeal.Status)
		}

		appeal.Status = StatusInProgress
		if appeal.IsNew() {
			t.Errorf("Expected appeal not to be new, got %s", appeal.Status)
		}
	})

	t.Run("IsInProgress", func(t *testing.T) {
		appeal.Status = StatusInProgress
		if !appeal.IsInProgress() {
			t.Errorf("Expected appeal to be in progress, got %s", appeal.Status)
		}

		appeal.Status = StatusNew
		if appeal.IsInProgress() {
			t.Errorf("Expected appeal not to be in progress, got %s", appeal.Status)
		}
	})

	t.Run("IsCompleted", func(t *testing.T) {
		appeal.Status = StatusCompleted
		if !appeal.IsCompleted() {
			t.Errorf("Expected appeal to be completed, got %s", appeal.Status)
		}

		appeal.Status = StatusNew
		if appeal.IsCompleted() {
			t.Errorf("Expected appeal not to be completed, got %s", appeal.Status)
		}
	})

	t.Run("IsCancelled", func(t *testing.T) {
		appeal.Status = StatusCancelled
		if !appeal.IsCancelled() {
			t.Errorf("Expected appeal to be cancelled, got %s", appeal.Status)
		}

		appeal.Status = StatusNew
		if appeal.IsCancelled() {
			t.Errorf("Expected appeal not to be cancelled, got %s", appeal.Status)
		}
	})
}

func TestAppealCanMethods(t *testing.T) {
	t.Parallel()

	appeal := &Appeal{
		ID:        "test-id",
		Theme:     "Test Theme",
		Message:   "Test Message",
		Status:    StatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("CanStartProcessing", func(t *testing.T) {
		appeal.Status = StatusNew
		if !appeal.CanStartProcessing() {
			t.Error("Expected appeal with StatusNew to be able to start processing")
		}

		appeal.Status = StatusCancelled
		if !appeal.CanStartProcessing() {
			t.Error("Expected appeal with StatusCancelled to be able to start processing")
		}

		appeal.Status = StatusInProgress
		if appeal.CanStartProcessing() {
			t.Error("Expected appeal with StatusInProgress not to be able to start processing")
		}

		appeal.Status = StatusCompleted
		if appeal.CanStartProcessing() {
			t.Error("Expected appeal with StatusCompleted not to be able to start processing")
		}
	})

	t.Run("CanComplete", func(t *testing.T) {
		appeal.Status = StatusInProgress
		if !appeal.CanComplete() {
			t.Error("Expected appeal with StatusInProgress to be able to complete")
		}

		appeal.Status = StatusNew
		if appeal.CanComplete() {
			t.Error("Expected appeal with StatusNew not to be able to complete")
		}

		appeal.Status = StatusCompleted
		if appeal.CanComplete() {
			t.Error("Expected appeal with StatusCompleted not to be able to complete")
		}

		appeal.Status = StatusCancelled
		if appeal.CanComplete() {
			t.Error("Expected appeal with StatusCancelled not to be able to complete")
		}
	})

	t.Run("CanCancel", func(t *testing.T) {
		appeal.Status = StatusNew
		if !appeal.CanCancel() {
			t.Error("Expected appeal with StatusNew to be able to cancel")
		}

		appeal.Status = StatusInProgress
		if !appeal.CanCancel() {
			t.Error("Expected appeal with StatusInProgress to be able to cancel")
		}

		appeal.Status = StatusCompleted
		if appeal.CanCancel() {
			t.Error("Expected appeal with StatusCompleted not to be able to cancel")
		}

		appeal.Status = StatusCancelled
		if appeal.CanCancel() {
			t.Error("Expected appeal with StatusCancelled not to be able to cancel")
		}
	})
}

func TestCreateAppealRequestValidation(t *testing.T) {
	t.Parallel()

	t.Run("ValidRequests", func(t *testing.T) {
		req := CreateAppealRequest{
			Theme:   "Valid Theme",
			Message: "Valid Message",
		}

		if req.Theme != "Valid Theme" {
			t.Errorf("Expected Theme to be 'Valid Theme', got %s", req.Theme)
		}
		if req.Message != "Valid Message" {
			t.Errorf("Expected Message to be 'Valid Message', got %s", req.Message)
		}
	})

	t.Run("EmptyFieldsAllowedAtStructLevel", func(t *testing.T) {
		req := CreateAppealRequest{
			Theme:   "",
			Message: "",
		}

		if req.Theme != "" || req.Message != "" {
			t.Error("Expected empty values to be allowed at struct level")
		}
	})
}

func TestUpdateAppealSolutionRequest(t *testing.T) {
	t.Parallel()

	req := UpdateAppealSolutionRequest{
		Solution: "This is a solution",
	}

	if req.Solution != "This is a solution" {
		t.Errorf("Expected Solution to be 'This is a solution', got %s", req.Solution)
	}
}

func TestUpdateAppealCancelRequest(t *testing.T) {
	t.Parallel()

	req := UpdateAppealCancelRequest{
		Reason: "Cancellation reason",
	}

	if req.Reason != "Cancellation reason" {
		t.Errorf("Expected Reason to be 'Cancellation reason', got %s", req.Reason)
	}
}

func TestFilterDatesRequest(t *testing.T) {
	t.Parallel()

	req := FilterDatesRequest{
		Date:      "2023-01-01",
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
	}

	if req.Date != "2023-01-01" {
		t.Errorf("Expected Date to be '2023-01-01', got %s", req.Date)
	}
	if req.StartDate != "2023-01-01" {
		t.Errorf("Expected StartDate to be '2023-01-01', got %s", req.StartDate)
	}
	if req.EndDate != "2023-12-31" {
		t.Errorf("Expected EndDate to be '2023-12-31', got %s", req.EndDate)
	}
}

