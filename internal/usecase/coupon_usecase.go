package usecase

import (
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
)

// CouponUsecase defines the business logic contract for coupons
type CouponUsecase interface {
	ClaimCoupon(userID string, couponName string) error
	GetCouponStatus(name string) (*domain.Coupon, error)
	CreateCoupon(name string, amount int) error
}
