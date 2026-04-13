package postgres

import (
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) domain.CouponRepository {
	return &couponRepository{db}
}

func (r *couponRepository) getDB(tx interface{}) *gorm.DB {
	if tx != nil {
		if gormDB, ok := tx.(*gorm.DB); ok {
			return gormDB
		}
	}
	return r.db
}

func (r *couponRepository) Create(coupon *domain.Coupon) error {
	return r.db.Create(coupon).Error
}

func (r *couponRepository) FindByName(name string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.db.First(&coupon, "name = ?", name).Error
	return &coupon, err
}

func (r *couponRepository) FindByNameWithLock(tx interface{}, name string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.getDB(tx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&coupon, "name = ?", name).Error
	return &coupon, err
}

func (r *couponRepository) UpdateStock(name string, amount int) error {
	return r.db.Model(&domain.Coupon{}).Where("name = ?", name).Update("remaining_amount", amount).Error
}

func (r *couponRepository) DecrementStock(tx interface{}, name string) error {
	return r.getDB(tx).Model(&domain.Coupon{}).
		Where("name = ? AND remaining_amount > 0", name).
		Update("remaining_amount", gorm.Expr("remaining_amount - 1")).Error
}
