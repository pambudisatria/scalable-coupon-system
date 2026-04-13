package postgres

import (
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"gorm.io/gorm"
)

type claimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) domain.ClaimRepository {
	return &claimRepository{db}
}

func (r *claimRepository) Create(claim *domain.Claim) error {
	return r.db.Create(claim).Error
}

func (r *claimRepository) Exists(userID string, couponName string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Claim{}).Where("user_id = ? AND coupon_name = ?", userID, couponName).Count(&count).Error
	return count > 0, err
}
