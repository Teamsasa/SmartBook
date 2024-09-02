package usecase

import (
	"SmartBook/internal/database"
	"SmartBook/internal/model"
	"time"
)

type MemoUseCase struct {
	db database.Service
}

func NewMemoUseCase(db database.Service) *MemoUseCase {
	return &MemoUseCase{
		db: db,
	}
}

func (u *MemoUseCase) GetMemos(userID string) ([]model.Memo, error) {
	// DBからユーザーIDを元にメモを取得
	query := "SELECT ID, UserID, ArticleURL, Content, CreatedAt, UpdatedAt FROM memos WHERE UserID = $1"
	rows, err := u.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []model.Memo
	for rows.Next() {
		var memo model.Memo
		if err := rows.Scan(&memo.ID, &memo.UserID, &memo.ArticleURL, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt); err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}

func (u *MemoUseCase) CreateMemo(memo *model.MemoRequest) error {

	// もしメモがすでに存在していたらupdate, そうでなければinsert
	query := "SELECT ID FROM memos WHERE UserID = $1 AND ArticleURL = $2"
	rows, err := u.db.Query(query, memo.UserID, memo.ArticleURL)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// update
		query = "UPDATE memos SET Content = $1, UpdatedAt = $2 WHERE UserID = $3 AND ArticleURL = $4"
		_, err := u.db.Exec(query, memo.Content, time.Now(), memo.UserID, memo.ArticleURL)
		if err != nil {
			return err
		}
	} else {
		// insert
		query = "INSERT INTO memos (UserID, ArticleURL, Content, CreatedAt, UpdatedAt) VALUES ($1, $2, $3, $4, $5)"
		_, err := u.db.Exec(query, memo.UserID, memo.ArticleURL, memo.Content, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (u *MemoUseCase) GetMemo(userID, memoID string) (*model.Memo, error) {
	query := "SELECT ID, UserID, ArticleURL, Content, CreatedAt, UpdatedAt FROM memos WHERE UserID = $1 AND ID = $2"
	rows, err := u.db.Query(query, userID, memoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var memo model.Memo
	if err := rows.Scan(&memo.ID, &memo.UserID, &memo.ArticleURL, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt); err != nil {
		return nil, err
	}

	return &memo, nil
}

func (u *MemoUseCase) DeleteMemo(userID, memoID string) error {
	query := "DELETE FROM memos WHERE UserID = $1 AND ID = $2"
	_, err := u.db.Exec(query, userID, memoID)
	if err != nil {
		return err
	}

	return nil
}
