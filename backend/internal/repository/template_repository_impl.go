package repository

import (
	"context"
	"database/sql"
	"github.com/45ai/backend/internal/model"
)

type templateRepositoryImpl struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) TemplateRepository {
	return &templateRepositoryImpl{db: db}
}

func (r *templateRepositoryImpl) GetAll(ctx context.Context) ([]model.Template, error) {
	query := "SELECT id, name, description, preview_image_url, credit_cost, is_active, created_at FROM templates WHERE is_active = true"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []model.Template
	for rows.Next() {
		var t model.Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.PreviewImageURL, &t.CreditCost, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *templateRepositoryImpl) GetByID(ctx context.Context, id int) (*model.Template, error) {
	query := "SELECT id, name, description, preview_image_url, credit_cost, is_active, created_at FROM templates WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	template := &model.Template{}
	err := row.Scan(&template.ID, &template.Name, &template.Description, &template.PreviewImageURL, &template.CreditCost, &template.IsActive, &template.CreatedAt)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *templateRepositoryImpl) GetByName(ctx context.Context, name string) (*model.Template, error) {
	query := "SELECT id, name, description, preview_image_url, credit_cost, is_active, created_at FROM templates WHERE name = ?"
	row := r.db.QueryRowContext(ctx, query, name)
	template := &model.Template{}
	err := row.Scan(&template.ID, &template.Name, &template.Description, &template.PreviewImageURL, &template.CreditCost, &template.IsActive, &template.CreatedAt)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *templateRepositoryImpl) Create(ctx context.Context, template *model.Template) error {
	query := "INSERT INTO templates (name, description, preview_image_url, credit_cost, is_active) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, template.Name, template.Description, template.PreviewImageURL, template.CreditCost, template.IsActive)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	template.ID = int(id)
	return nil
}

func (r *templateRepositoryImpl) Update(ctx context.Context, template *model.Template) error {
	// To be implemented in a future task
	return nil
}

func (r *templateRepositoryImpl) SetActive(ctx context.Context, id int, isActive bool) error {
	// To be implemented in a future task
	return nil
}

func (r *templateRepositoryImpl) Count(ctx context.Context) (int, error) {
	// To be implemented in a future task
	return 0, nil
} 