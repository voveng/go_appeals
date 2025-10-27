package services

import (
	"fmt"
	"go_appeals/models"
	"go_appeals/repository"
	"time"
)

type AppealService struct {
	repo *repository.AppealRepository
}

func NewAppealService(repo *repository.AppealRepository) *AppealService {
	return &AppealService{
		repo: repo,
	}
}

func (s *AppealService) CreateAppeal(req models.CreateAppealRequest) (*models.Appeal, error) {
	appeal := &models.Appeal{
		Theme:   req.Theme,
		Message: req.Message,
		Status:  models.StatusNew,
	}

	if appeal.Theme == "" || appeal.Message == "" {
		return nil, fmt.Errorf("theme and message are required")
	}

	savedAppeal, err := s.repo.Save(appeal)
	if err != nil {
		return nil, fmt.Errorf("failed to save appeal: %w", err)
	}
	return savedAppeal, nil
}

func (s *AppealService) GetStartedAppeals() ([]*models.Appeal, error) {
	allAppeals, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all appeals: %w", err)
	}

	startedAppeals := make([]*models.Appeal, 0)
	for _, appeal := range allAppeals {
		if appeal.IsNew() || appeal.IsInProgress() {
			startedAppeals = append(startedAppeals, appeal)
		}
	}

	return startedAppeals, nil
}

func (s *AppealService) GetAllAppeals() ([]*models.Appeal, error) {
	return s.repo.GetAll()
}

func (s *AppealService) GetAppealByID(id string) (*models.Appeal, error) {
	return s.repo.FindByID(id)
}

func (s *AppealService) StartProcessing(id string) (*models.Appeal, error) {
	appeal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !appeal.CanStartProcessing() {
		return nil, fmt.Errorf("cannot start processing appeal with status: %s", appeal.Status)
	}

	appeal.Status = models.StatusInProgress

	updatedAppeal, err := s.repo.Update(appeal)
	if err != nil {
		return nil, err
	}

	return updatedAppeal, nil
}

func (s *AppealService) CancelAppeal(id string) error {
	appeal, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !appeal.CanCancel() {
		return fmt.Errorf("cannot cancel appeal with status: %s", appeal.Status)
	}

	appeal.Status = models.StatusCancelled

	_, err = s.repo.Update(appeal)
	if err != nil {
		return err
	}

	return nil
}

func (s *AppealService) CancelAllInProgress() error {
	return s.repo.CancelInProgressAppeals()
}

func (s *AppealService) GetAppealsByDates(start, end time.Time) ([]*models.Appeal, error) {
	return s.repo.SelectAppealsByDates(start, end)
}

func (s *AppealService) CompleteAppeal(id string, req models.UpdateAppealSolutionRequest) (*models.Appeal, error) {
	appeal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !appeal.CanComplete() {
		return nil, fmt.Errorf("cannot complete appeal with status: %s", appeal.Status)
	}

	appeal.Status = models.StatusCompleted
	appeal.Solution = req.Solution
	updatedAppeal, err := s.repo.Update(appeal)
	if err != nil {
		return nil, err
	}

	return updatedAppeal, nil
}
