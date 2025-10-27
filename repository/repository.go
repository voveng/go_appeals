package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go_appeals/models"
)

type AppealRepository struct {
	db *sql.DB
}

func NewAppealRepository(dbPath string) (*AppealRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &AppealRepository{db: db}

	if err := repo.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return repo, nil
}

func (r *AppealRepository) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS appeals (
		id TEXT PRIMARY KEY,
		theme TEXT NOT NULL,
		message TEXT NOT NULL,
		status TEXT NOT NULL,
		solution TEXT,
		cansel_reason TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`
	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create appeals table: %w", err)
	}
	log.Println("Appeals table initialized or already exists.")
	return nil
}

func (r *AppealRepository) Save(appeal *models.Appeal) (*models.Appeal, error) {
	appeal.ID = uuid.New().String()
	appeal.CreatedAt = time.Now()
	appeal.UpdatedAt = time.Now()

	stmt, err := r.db.Prepare(
		"INSERT INTO appeals (id, theme, message, status, solution, cansel_reason, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare save statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		appeal.ID,
		appeal.Theme,
		appeal.Message,
		appeal.Status,
		appeal.Solution,
		appeal.CanselReason,
		appeal.CreatedAt,
		appeal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute save statement: %w", err)
	}

	return appeal, nil
}

func (r *AppealRepository) Update(appeal *models.Appeal) (*models.Appeal, error) {
	appeal.UpdatedAt = time.Now()

	stmt, err := r.db.Prepare(
		"UPDATE appeals SET theme=?, message=?, status=?, solution=?, cansel_reason=?, updated_at=? WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		appeal.Theme,
		appeal.Message,
		appeal.Status,
		appeal.Solution,
		appeal.CanselReason,
		appeal.UpdatedAt,
		appeal.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update statement: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("appeal with ID %s not found for update", appeal.ID)
	}

	return appeal, nil
}

func (r *AppealRepository) FindByID(id string) (*models.Appeal, error) {
	row := r.db.QueryRow(
		"SELECT id, theme, message, status, solution, cansel_reason, created_at, updated_at FROM appeals WHERE id = ?", id)

	appeal := &models.Appeal{}
	err := row.Scan(
		&appeal.ID,
		&appeal.Theme,
		&appeal.Message,
		&appeal.Status,
		&appeal.Solution,
		&appeal.CanselReason,
		&appeal.CreatedAt,
		&appeal.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("appeal with ID %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan appeal: %w", err)
	}

	return appeal, nil
}

func (r *AppealRepository) GetAll() ([]*models.Appeal, error) {
	rows, err := r.db.Query(
		"SELECT id, theme, message, status, solution, cansel_reason, created_at, updated_at FROM appeals")
	if err != nil {
		return nil, fmt.Errorf("failed to query appeals: %w", err)
	}
	defer rows.Close()

	appeals := make([]*models.Appeal, 0)
	for rows.Next() {
		appeal := &models.Appeal{}
		err := rows.Scan(
			&appeal.ID,
			&appeal.Theme,
			&appeal.Message,
			&appeal.Status,
			&appeal.Solution,
			&appeal.CanselReason,
			&appeal.CreatedAt,
			&appeal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appeal row: %w", err)
		}
		appeals = append(appeals, appeal)
	}

	return appeals, nil
}

func (r *AppealRepository) CancelInProgressAppeals() error {
	stmt, err := r.db.Prepare(
		"UPDATE appeals SET status = ?, updated_at = ? WHERE status IN (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare cancel statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(models.StatusCancelled, time.Now(), models.StatusNew, models.StatusInProgress)
	if err != nil {
		return fmt.Errorf("failed to execute cancel statement: %w", err)
	}

	return nil
}

func (r *AppealRepository) SelectAppealsByDates(start, end time.Time) ([]*models.Appeal, error) {
	rows, err := r.db.Query(
		"SELECT id, theme, message, status, solution, cansel_reason, created_at, updated_at FROM appeals WHERE created_at BETWEEN ? AND ?",
		start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query appeals: %w", err)
	}
	defer rows.Close()

	appeals := make([]*models.Appeal, 0)
	for rows.Next() {
		appeal := &models.Appeal{}
		err := rows.Scan(
			&appeal.ID,
			&appeal.Theme,
			&appeal.Message,
			&appeal.Status,
			&appeal.Solution,
			&appeal.CanselReason,
			&appeal.CreatedAt,
			&appeal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appeal row: %w", err)
		}
		appeals = append(appeals, appeal)
	}

	return appeals, nil
}

func (r *AppealRepository) Close() error {
	return r.db.Close()
}
