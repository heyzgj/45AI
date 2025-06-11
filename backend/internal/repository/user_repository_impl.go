package repository

import (
	"context"
	"database/sql"

	"github.com/45ai/backend/internal/model"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (wechat_openid, nickname, avatar_url, credits) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.WechatOpenID, user.Nickname, user.AvatarURL, user.Credits)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, wechat_openid, nickname, avatar_url, credits, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	user := &model.User{}
	var nickname, avatarURL sql.NullString
	err := row.Scan(&user.ID, &user.WechatOpenID, &nickname, &avatarURL, &user.Credits, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Nickname = nickname.String
	user.AvatarURL = avatarURL.String
	return user, nil
}

func (r *userRepositoryImpl) GetByWechatOpenID(ctx context.Context, openID string) (*model.User, error) {
	query := "SELECT id, wechat_openid, nickname, avatar_url, credits, created_at, updated_at FROM users WHERE wechat_openid = ?"
	row := r.db.QueryRowContext(ctx, query, openID)
	user := &model.User{}
	var nickname, avatarURL sql.NullString
	err := row.Scan(&user.ID, &user.WechatOpenID, &nickname, &avatarURL, &user.Credits, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Nickname = nickname.String
	user.AvatarURL = avatarURL.String
	return user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET nickname = ?, avatar_url = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Nickname, user.AvatarURL, user.ID)
	return err
}

func (r *userRepositoryImpl) UpdateCredits(ctx context.Context, userID int64, amount int) error {
	query := "UPDATE users SET credits = credits + ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, amount, userID)
	return err
}

func (r *userRepositoryImpl) Exists(ctx context.Context, openID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE wechat_openid = ?)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, openID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
} 