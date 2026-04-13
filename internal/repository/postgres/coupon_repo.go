package postgres

import (
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"gorm.io/gorm"
)

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) domain.CouponRepository {
	return &couponRepository{db}
}

func (r *couponRepository) Create(coupon *domain.Coupon) error {
	return r.db.Create(coupon).Error
}

func (r *couponRepository) FindByName(name string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.db.First(&coupon, "name = ?", name).Error
	return &coupon, err
}

func (r *couponRepository) UpdateStock(name string, amount int) error {
	return r.db.Model(&domain.Coupon{}).Where("name = ?", name).Update("remaining_amount", amount).Error
}
