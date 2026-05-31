package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrCouponAmountInvalid = errors.New("coupon amount must be greater than zero")
	ErrCouponDuplicate     = errors.New("coupon name already exists")
	ErrCouponNotFound      = errors.New("coupon not found")
	ErrCouponNoStock       = errors.New("coupon out of stock")
	ErrUserAlreadyClaimed  = errors.New("user has already claimed this coupon")
)

type couponUsecase struct {
	db         *gorm.DB
	couponRepo domain.CouponRepository
	claimRepo  domain.ClaimRepository
}

// NewCouponUsecase creates a new instance of CouponUsecase
func NewCouponUsecase(db *gorm.DB, couponRepo domain.CouponRepository, claimRepo domain.ClaimRepository) CouponUsecase {
	return &couponUsecase{
		db:         db,
		couponRepo: couponRepo,
		claimRepo:  claimRepo,
	}
}

// ... CreateCoupon and other methods ...

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
		// Detect duplicate key violation (unique constraint on primary key Name)
		// This is the last-defense guard for concurrent requests that bypass FindByName
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCouponDuplicate
		}
		return fmt.Errorf("failed to create coupon: %w", err)
	}

	return nil
}

func (u *couponUsecase) ClaimCoupon(userID string, couponName string) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check if user already claimed
		exists, err := u.claimRepo.ExistsWithTx(tx, userID, couponName)
		if err != nil {
			return fmt.Errorf("failed to check claim existence: %w", err)
		}
		if exists {
			return ErrUserAlreadyClaimed
		}

		// 2. Lock coupon row and check stock
		coupon, err := u.couponRepo.FindByNameWithLock(tx, couponName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCouponNotFound
			}
			return fmt.Errorf("failed to find coupon with lock: %w", err)
		}

		if coupon.RemainingAmount <= 0 {
			return ErrCouponNoStock
		}

		// 3. Insert claim
		claim := &domain.Claim{
			UserID:     userID,
			CouponName: couponName,
		}
		err = u.claimRepo.CreateWithTx(tx, claim)
		if err != nil {
			// Handle unique constraint violation (concurrency safety last defense)
			// In GORM/Postgres, we check for duplicate key error
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrUserAlreadyClaimed
			}
			return fmt.Errorf("failed to create claim: %w", err)
		}

		// 4. Decrement stock
		err = u.couponRepo.DecrementStock(tx, couponName)
		if err != nil {
			return fmt.Errorf("failed to decrement stock: %w", err)
		}

		log.Printf("[CONCURRENCY] User %s successfully claimed coupon %s", userID, couponName)
		return nil
	})
}

func (u *couponUsecase) GetCouponStatus(name string) (*CouponStatusResponse, error) {
	coupon, err := u.couponRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	claims, err := u.claimRepo.GetClaimedUsersByCoupon(name)
	if err != nil {
		return nil, err
	}

	return &CouponStatusResponse{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: coupon.RemainingAmount,
		ClaimedBy:       claims,
	}, nil
}
