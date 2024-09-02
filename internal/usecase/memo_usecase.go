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

// もしメモがすでに存在していたらupdate, そうでなければinsert
func (u *MemoUseCase) UpsertMemo(memo *model.MemoRequest) error {

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

func (u *MemoUseCase) GetMemo(memo *model.MemoRequest) (*model.Memo, error) {
	query := "SELECT ID, UserID, ArticleURL, Content, CreatedAt, UpdatedAt FROM memos WHERE UserID = $1 AND ArticleURL = $2"
	rows, err := u.db.Query(query, memo.UserID, memo.ArticleURL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var m model.Memo
	if rows.Next() {
		if err := rows.Scan(&m.ID, &m.UserID, &m.ArticleURL, &m.Content, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (u *MemoUseCase) DeleteMemo(memo *model.MemoRequest) error {
	query := "DELETE FROM memos WHERE UserID = $1 AND ArticleURL = $2"
	_, err := u.db.Exec(query, memo.UserID, memo.ArticleURL)
	if err != nil {
		return err
	}

	return nil
}
