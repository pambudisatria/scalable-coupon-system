package domain

import (
	"time"
)

// Claim represents the record of a user claiming a coupon
type Claim struct {
	UserID     string    `gorm:"primaryKey;uniqueIndex:idx_user_coupon" json:"user_id"`
	CouponName string    `gorm:"primaryKey;uniqueIndex:idx_user_coupon" json:"coupon_name"`
	CreatedAt  time.Time `json:"created_at"`
}

// ClaimRepository defines the data access behavior for Claims
type ClaimRepository interface {
	Create(claim *Claim) error
	CreateWithTx(tx interface{}, claim *Claim) error
	Exists(userID string, couponName string) (bool, error)
	ExistsWithTx(tx interface{}, userID string, couponName string) (bool, error)
	GetClaimedUsersByCoupon(couponName string) ([]string, error)
}
