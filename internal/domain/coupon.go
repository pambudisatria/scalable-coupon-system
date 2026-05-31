package domain

import "errors"

// Sentinel error when an update did not affect any rows
var ErrNoRowsAffected = errors.New("update affected no rows")

// Coupon represents the core coupon entity
type Coupon struct {
	Name            string `gorm:"primaryKey" json:"name"`
	Amount          int    `json:"amount"`
	RemainingAmount int    `json:"remaining_amount"`
}

// CouponRepository defines the data access behavior for Coupons
type CouponRepository interface {
	Create(coupon *Coupon) error
	FindByName(name string) (*Coupon, error)
	FindByNameWithLock(tx interface{}, name string) (*Coupon, error)
	UpdateStock(name string, amount int) error
	DecrementStock(tx interface{}, name string) error
}
