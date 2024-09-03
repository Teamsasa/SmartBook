package usecase

import (
	"SmartBook/internal/model"
	"time"

	"gorm.io/gorm"
)

type MemoUseCase struct {
	db *gorm.DB
}

func NewMemoUseCase(db *gorm.DB) *MemoUseCase {
	return &MemoUseCase{
		db: db,
	}
}

func (u *MemoUseCase) GetMemos(userID string) ([]model.MemoData, error) {
	var memos []model.MemoData
	result := u.db.Where("user_id = ?", userID).Find(&memos)
	if result.Error != nil {
		return nil, result.Error
	}

	return memos, nil
}

// articleとmemoを作成する。どちらが失敗したらロールバックする。
func (u *MemoUseCase) CreateMemo(memoCreateReq *model.MemoData, articleCreateReq *model.ArticleData) error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Create(articleCreateReq)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(memoCreateReq)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u *MemoUseCase) UpdateMemo(req *model.MemoRequest) error {
	memo := model.MemoData{
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	result := u.db.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).Updates(&memo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *MemoUseCase) GetMemo(req *model.MemoRequest) (*model.MemoData, error) {
	var memo model.MemoData
	result := u.db.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&memo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &memo, nil
}

func (u *MemoUseCase) DeleteMemo(req *model.MemoRequest) error {
	result := u.db.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).Delete(&model.MemoData{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
