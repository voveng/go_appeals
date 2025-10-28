package repository

import (
	"database/sql"
	"go_appeals/internal/models"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestRepository(t *testing.T) (*AppealRepository, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	repo := &AppealRepository{db: db}
	if err := repo.InitSchema(); err != nil {
		t.Fatalf("Failed to initialize database schema: %v", err)
	}

	return repo, func() {
		if err := repo.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}
}

func TestInitSchema(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	_ = repo
}

func TestClose(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	repo := &AppealRepository{db: db}
	if err := repo.Close(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func TestSelectAppealsByDates(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	startDate := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)  // Начало дня
	endDate := time.Date(2025, 10, 27, 23, 59, 59, 0, time.UTC) // Конец дня

	appeals := []*models.Appeal{
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       "open",
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Date(2025, 10, 27, 12, 0, 0, 0, time.UTC), // Середина дня
			UpdatedAt:    time.Date(2025, 10, 27, 12, 0, 0, 0, time.UTC),
		},
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       "open",
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Date(2025, 10, 27, 18, 0, 0, 0, time.UTC), // Вечер дня
			UpdatedAt:    time.Date(2025, 10, 27, 18, 0, 0, 0, time.UTC),
		},
	}

	for _, appeal := range appeals {
		_, err := repo.Save(appeal)
		if err != nil {
			t.Errorf("Failed to save appeal: %v", err)
		}
	}

	selectedAppeals, err := repo.SelectAppealsByDates(startDate, endDate)
	if err != nil {
		t.Errorf("Failed to select appeals by dates: %v", err)
	}

	if len(selectedAppeals) != len(appeals) {
		t.Errorf("Expected %d appeals, got %d", len(appeals), len(selectedAppeals))
	}

	if len(selectedAppeals) == len(appeals) {
		for i, appeal := range appeals {
			if selectedAppeals[i].ID != appeal.ID {
				t.Errorf("Expected appeal ID %s, got %s", appeal.ID, selectedAppeals[i].ID)
			}
		}
	}
}

func TestCancelInProgressAppeals(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	appeals := []*models.Appeal{
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       models.StatusNew,
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       models.StatusInProgress,
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       models.StatusCompleted,
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Theme:        "Test theme",
			Message:      "Test message",
			Status:       models.StatusCancelled,
			Solution:     "Test solution",
			CanselReason: "Test reason",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	savedAppeals := make([]*models.Appeal, len(appeals))
	for i, appeal := range appeals {
		saved, err := repo.Save(appeal)
		if err != nil {
			t.Errorf("Failed to save appeal: %v", err)
		} else {
			savedAppeals[i] = saved
		}
	}

	allAppeals, err := repo.GetAll()
	if err != nil {
		t.Errorf("Failed to get all appeals: %v", err)
	}

	initialInProgressCount := 0
	for _, appeal := range allAppeals {
		if appeal.Status == models.StatusNew || appeal.Status == models.StatusInProgress {
			initialInProgressCount++
		}
	}

	if initialInProgressCount != 2 {
		t.Errorf("Expected 2 appeals with StatusNew or StatusInProgress, got %d", initialInProgressCount)
	}

	if err := repo.CancelInProgressAppeals(); err != nil {
		t.Errorf("Failed to cancel in-progress appeals: %v", err)
	}

	allAppealsAfter, err := repo.GetAll()
	if err != nil {
		t.Errorf("Failed to get all appeals: %v", err)
	}

	cancelledCount := 0
	for _, appeal := range allAppealsAfter {
		if appeal.Status == models.StatusCancelled {
			cancelledCount++
		}
	}

	if cancelledCount != 3 {
		t.Errorf("Expected 3 cancelled appeals (2 converted + 1 already cancelled), got %d", cancelledCount)
	}

	newAppeal := savedAppeals[0]
	inProgressAppeal := savedAppeals[1]

	updatedNewAppeal, err := repo.FindByID(newAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal with original ID %s: %v", newAppeal.ID, err)
	} else if updatedNewAppeal.Status != models.StatusCancelled {
		t.Errorf("Expected appeal %s to be cancelled, got %s", newAppeal.ID, updatedNewAppeal.Status)
	}

	updatedInProgressAppeal, err := repo.FindByID(inProgressAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal with original ID %s: %v", inProgressAppeal.ID, err)
	} else if updatedInProgressAppeal.Status != models.StatusCancelled {
		t.Errorf("Expected appeal %s to be cancelled, got %s", inProgressAppeal.ID, updatedInProgressAppeal.Status)
	}

	completedAppeal := savedAppeals[2]
	alreadyCancelledAppeal := savedAppeals[3]

	updatedCompletedAppeal, err := repo.FindByID(completedAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal with original ID %s: %v", completedAppeal.ID, err)
	} else if updatedCompletedAppeal.Status != models.StatusCompleted {
		t.Errorf("Expected appeal %s to remain completed, got %s", completedAppeal.ID, updatedCompletedAppeal.Status)
	}

	updatedAlreadyCancelledAppeal, err := repo.FindByID(alreadyCancelledAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal with original ID %s: %v", alreadyCancelledAppeal.ID, err)
	} else if updatedAlreadyCancelledAppeal.Status != models.StatusCancelled {
		t.Errorf("Expected appeal %s to remain cancelled, got %s", alreadyCancelledAppeal.ID, updatedAlreadyCancelledAppeal.Status)
	}
}

func TestSave(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	appeal := &models.Appeal{
		Theme:        "Test theme",
		Message:      "Test message",
		Status:       models.StatusNew,
		Solution:     "Test solution",
		CanselReason: "Test reason",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	savedAppeal, err := repo.Save(appeal)
	if err != nil {
		t.Errorf("Failed to save appeal: %v", err)
	}

	if savedAppeal.ID == "" {
		t.Errorf("Expected saved appeal to have a non-empty ID, got %s", savedAppeal.ID)
	}
}

func TestGetAll(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	appeals := []*models.Appeal{
		{
			Theme:        "Test theme 1",
			Message:      "Test message 1",
			Status:       models.StatusNew,
			Solution:     "Test solution 1",
			CanselReason: "Test reason 1",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Theme:        "Test theme 2",
			Message:      "Test message 2",
			Status:       models.StatusNew,
			Solution:     "Test solution 2",
			CanselReason: "Test reason 2",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, appeal := range appeals {
		_, err := repo.Save(appeal)
		if err != nil {
			t.Errorf("Failed to save appeal: %v", err)
		}
	}

	savedAppeals, err := repo.GetAll()
	if err != nil {
		t.Errorf("Failed to get all appeals: %v", err)
	}

	if len(savedAppeals) != len(appeals) {
		t.Errorf("Expected %d appeals, got %d", len(appeals), len(savedAppeals))
	}
}

func TestFindByID(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	appeal := &models.Appeal{
		Theme:        "Test theme",
		Message:      "Test message",
		Status:       models.StatusNew,
		Solution:     "Test solution",
		CanselReason: "Test reason",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	savedAppeal, err := repo.Save(appeal)
	if err != nil {
		t.Errorf("Failed to save appeal: %v", err)
	}

	foundAppeal, err := repo.FindByID(savedAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal by ID: %v", err)
	}

	if foundAppeal.ID != savedAppeal.ID {
		t.Errorf("Expected appeal with ID %s, got %s", savedAppeal.ID, foundAppeal.ID)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	repo, cleanup := newTestRepository(t)
	defer cleanup()

	originalAppeal := &models.Appeal{
		Theme:        "Original theme",
		Message:      "Original message",
		Status:       models.StatusNew,
		Solution:     "Original solution",
		CanselReason: "Original reason",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	savedAppeal, err := repo.Save(originalAppeal)
	if err != nil {
		t.Fatalf("Failed to save appeal: %v", err)
	}

	updatedAppeal := &models.Appeal{
		ID:           savedAppeal.ID,
		Theme:        "Updated theme",
		Message:      "Updated message",
		Status:       models.StatusInProgress,
		Solution:     "Updated solution",
		CanselReason: "Updated reason",
		CreatedAt:    savedAppeal.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	returnedAppeal, err := repo.Update(updatedAppeal)
	if err != nil {
		t.Errorf("Failed to update appeal: %v", err)
	}

	if returnedAppeal.Theme != "Updated theme" {
		t.Errorf("Expected theme %s, got %s", "Updated theme", returnedAppeal.Theme)
	}
	if returnedAppeal.Status != models.StatusInProgress {
		t.Errorf("Expected status %s, got %s", models.StatusInProgress, returnedAppeal.Status)
	}

	foundAppeal, err := repo.FindByID(savedAppeal.ID)
	if err != nil {
		t.Errorf("Failed to find appeal by ID: %v", err)
	}

	if foundAppeal.Theme != "Updated theme" {
		t.Errorf("Expected theme %s in database, got %s", "Updated theme", foundAppeal.Theme)
	}
	if foundAppeal.Status != models.StatusInProgress {
		t.Errorf("Expected status %s in database, got %s", models.StatusInProgress, foundAppeal.Status)
	}

	if foundAppeal.UpdatedAt.Equal(savedAppeal.UpdatedAt) {
		t.Errorf("Expected UpdatedAt to be different after update")
	}

	nonExistentAppeal := &models.Appeal{
		ID:           "non-existent",
		Theme:        "Test theme",
		Message:      "Test message",
		Status:       models.StatusNew,
		Solution:     "Test solution",
		CanselReason: "Test reason",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = repo.Update(nonExistentAppeal)
	if err == nil {
		t.Errorf("Expected error when updating non-existent appeal")
	}
}
