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

func (r *claimRepository) getDB(tx interface{}) *gorm.DB {
	if tx != nil {
		if gormDB, ok := tx.(*gorm.DB); ok {
			return gormDB
		}
	}
	return r.db
}

func (r *claimRepository) Create(claim *domain.Claim) error {
	return r.db.Create(claim).Error
}

func (r *claimRepository) CreateWithTx(tx interface{}, claim *domain.Claim) error {
	return r.getDB(tx).Create(claim).Error
}

func (r *claimRepository) Exists(userID string, couponName string) (bool, error) {
	return r.ExistsWithTx(nil, userID, couponName)
}

func (r *claimRepository) ExistsWithTx(tx interface{}, userID string, couponName string) (bool, error) {
	var count int64
	err := r.getDB(tx).Model(&domain.Claim{}).Where("user_id = ? AND coupon_name = ?", userID, couponName).Count(&count).Error
	return count > 0, err
}

func (r *claimRepository) GetClaimedUsersByCoupon(couponName string) ([]string, error) {
	var userIDs []string
	err := r.db.Model(&domain.Claim{}).Where("coupon_name = ?", couponName).Pluck("user_id", &userIDs).Error
	if userIDs == nil {
		userIDs = []string{}
	}
	return userIDs, err
}
