package usecase


// CouponStatusResponse represents the unified response containing coupon info and claims
type CouponStatusResponse struct {
	Name            string   `json:"name"`
	Amount          int      `json:"amount"`
	RemainingAmount int      `json:"remaining_amount"`
	ClaimedBy       []string `json:"claimed_by"`
}

// CouponUsecase defines the business logic contract for coupons
type CouponUsecase interface {
	ClaimCoupon(userID string, couponName string) error
	GetCouponStatus(name string) (*CouponStatusResponse, error)
	CreateCoupon(name string, amount int) error
}
