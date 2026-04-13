package usecase

import (
	"errors"
	"fmt"

	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
)

var (
	ErrCouponAmountInvalid = errors.New("coupon amount must be greater than zero")
	ErrCouponDuplicate     = errors.New("coupon name already exists")
)

type couponUsecase struct {
	couponRepo domain.CouponRepository
	claimRepo  domain.ClaimRepository
}

// NewCouponUsecase creates a new instance of CouponUsecase
func NewCouponUsecase(couponRepo domain.CouponRepository, claimRepo domain.ClaimRepository) CouponUsecase {
	return &couponUsecase{
		couponRepo: couponRepo,
		claimRepo:  claimRepo,
	}
}

func (u *couponUsecase) CreateCoupon(name string, amount int) error {
	if amount <= 0 {
		return ErrCouponAmountInvalid
	}

	// Check if coupon already exists
	existing, _ := u.couponRepo.FindByName(name)
	if existing != nil && existing.Name != "" {
		return ErrCouponDuplicate
	}

	coupon := &domain.Coupon{
		Name:            name,
		Amount:          amount,
		RemainingAmount: amount,
	}

	err := u.couponRepo.Create(coupon)
	if err != nil {
		return fmt.Errorf("failed to create coupon: %w", err)
	}

	return nil
}

func (u *couponUsecase) ClaimCoupon(userID string, couponName string) error {
	return nil // Placeholder for now
}

func (u *couponUsecase) GetCouponStatus(name string) (*domain.Coupon, error) {
	return u.couponRepo.FindByName(name)
}
