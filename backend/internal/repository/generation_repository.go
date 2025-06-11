package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/45ai/backend/internal/model"
)

// GenerationRepository defines the interface for generation data access
type GenerationRepository interface {
	Create(ctx context.Context, generation *model.Generation) error
	GetByJobID(ctx context.Context, jobID string) (*model.Generation, error)
	UpdateStatus(ctx context.Context, jobID string, status string, progress int) error
	UpdateWithResult(ctx context.Context, jobID string, status string, imageURL string, completedAt time.Time) error
	UpdateWithError(ctx context.Context, jobID string, error string, completedAt time.Time) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Generation, error)
}

type generationRepositoryImpl struct {
	db *sql.DB
}

// NewGenerationRepository creates a new generation repository
func NewGenerationRepository(db *sql.DB) GenerationRepository {
	return &generationRepositoryImpl{db: db}
}

func (r *generationRepositoryImpl) Create(ctx context.Context, generation *model.Generation) error {
	query := `
		INSERT INTO generations (job_id, user_id, template_id, status, progress, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	generation.CreatedAt = now
	generation.UpdatedAt = now
	
	result, err := r.db.ExecContext(ctx, query,
		generation.JobID,
		generation.UserID,
		generation.TemplateID,
		generation.Status,
		generation.Progress,
		generation.CreatedAt,
		generation.UpdatedAt,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	generation.ID = id
	return nil
}

func (r *generationRepositoryImpl) GetByJobID(ctx context.Context, jobID string) (*model.Generation, error) {
	query := `
		SELECT id, job_id, user_id, template_id, status, progress, image_url, error, 
			   started_at, completed_at, created_at, updated_at
		FROM generations 
		WHERE job_id = ?
	`
	
	generation := &model.Generation{}
	var imageURL, errorMsg sql.NullString
	var startedAt, completedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, jobID).Scan(
		&generation.ID,
		&generation.JobID,
		&generation.UserID,
		&generation.TemplateID,
		&generation.Status,
		&generation.Progress,
		&imageURL,
		&errorMsg,
		&startedAt,
		&completedAt,
		&generation.CreatedAt,
		&generation.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if imageURL.Valid {
		generation.ImageURL = imageURL.String
	}
	if errorMsg.Valid {
		generation.Error = errorMsg.String
	}
	if startedAt.Valid {
		generation.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		generation.CompletedAt = &completedAt.Time
	}
	
	return generation, nil
}

func (r *generationRepositoryImpl) UpdateStatus(ctx context.Context, jobID string, status string, progress int) error {
	query := `
		UPDATE generations 
		SET status = ?, progress = ?, updated_at = ?, started_at = CASE 
			WHEN status = 'pending' AND ? = 'processing' AND started_at IS NULL THEN ?
			ELSE started_at 
		END
		WHERE job_id = ?
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, progress, now, status, now, jobID)
	return err
}

func (r *generationRepositoryImpl) UpdateWithResult(ctx context.Context, jobID string, status string, imageURL string, completedAt time.Time) error {
	query := `
		UPDATE generations 
		SET status = ?, image_url = ?, completed_at = ?, updated_at = ?
		WHERE job_id = ?
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, imageURL, completedAt, now, jobID)
	return err
}

func (r *generationRepositoryImpl) UpdateWithError(ctx context.Context, jobID string, error string, completedAt time.Time) error {
	query := `
		UPDATE generations 
		SET status = 'failed', error = ?, completed_at = ?, updated_at = ?
		WHERE job_id = ?
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, error, completedAt, now, jobID)
	return err
}

func (r *generationRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Generation, error) {
	query := `
		SELECT id, job_id, user_id, template_id, status, progress, image_url, error,
			   started_at, completed_at, created_at, updated_at
		FROM generations 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var generations []*model.Generation
	for rows.Next() {
		generation := &model.Generation{}
		var imageURL, errorMsg sql.NullString
		var startedAt, completedAt sql.NullTime
		
		err := rows.Scan(
			&generation.ID,
			&generation.JobID,
			&generation.UserID,
			&generation.TemplateID,
			&generation.Status,
			&generation.Progress,
			&imageURL,
			&errorMsg,
			&startedAt,
			&completedAt,
			&generation.CreatedAt,
			&generation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if imageURL.Valid {
			generation.ImageURL = imageURL.String
		}
		if errorMsg.Valid {
			generation.Error = errorMsg.String
		}
		if startedAt.Valid {
			generation.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			generation.CompletedAt = &completedAt.Time
		}
		
		generations = append(generations, generation)
	}
	
	return generations, rows.Err()
} 