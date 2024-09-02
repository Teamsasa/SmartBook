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

// もしメモがすでに存在していたらupdate, そうでなければinsert
func (u *MemoUseCase) UpsertMemo(req *model.MemoRequest) error {
	// 既存のメモを取得
	var memo model.MemoData
	result := u.db.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&memo)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// メモが存在しない場合は新規作成
	if result.Error == gorm.ErrRecordNotFound {
		memo = model.MemoData{
			UserID:    req.UserID,
			ArticleID: req.ArticleID,
			Content:   req.Content,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}
		result = u.db.Create(&memo)
		if result.Error != nil {
			return result.Error
		}
	} else {
		// メモが存在する場合は更新
		result = u.db.Model(&memo).Update("content", req.Content)
		if result.Error != nil {
			return result.Error
		}
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
