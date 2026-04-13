package usecase

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"github.com/pambudisatria/scalable-coupon-system/internal/repository/postgres"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestClaimCoupon_Concurrency(t *testing.T) {
	// 1. Setup DB Connection (Ensure Postgres is running)
	dsn := "host=localhost user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable"
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Postgres: %v", err)
	}

	// 2. Clean up and Migrate
	db.Exec("DELETE FROM claims")
	db.Exec("DELETE FROM coupons")
	db.AutoMigrate(&domain.Coupon{}, &domain.Claim{})

	// 3. Setup Repos and Usecase
	couponRepo := postgres.NewCouponRepository(db)
	claimRepo := postgres.NewClaimRepository(db)
	uc := NewCouponUsecase(db, couponRepo, claimRepo)

	// 4. Create a coupon with limited stock
	couponName := "FLASH_SALE_10"
	stock := 10
	err = uc.CreateCoupon(couponName, stock)
	if err != nil {
		t.Fatalf("Failed to create test coupon: %v", err)
	}

	// 5. Simulate 50 concurrent users
	numUsers := 50
	var wg sync.WaitGroup
	results := make(chan error, numUsers)

	for i := 1; i <= numUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", id)
			err := uc.ClaimCoupon(userID, couponName)
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	// 6. Analyze results
	successCount := 0
	alreadyClaimedCount := 0
	noStockCount := 0
	otherErrorCount := 0

	for err := range results {
		if err == nil {
			successCount++
		} else if err == ErrUserAlreadyClaimed {
			alreadyClaimedCount++
		} else if err == ErrCouponNoStock {
			noStockCount++
		} else {
			otherErrorCount++
			t.Logf("Unexpected error: %v", err)
		}
	}

	// 7. Verify correctness
	if successCount != stock {
		t.Errorf("Expected exactly %d successes, got %d", stock, successCount)
	}
	
	expectedNoStock := numUsers - stock
	if noStockCount != expectedNoStock {
		t.Errorf("Expected %d stock-out errors, got %d", expectedNoStock, noStockCount)
	}

	// 8. Final DB verification
	var finalCoupon domain.Coupon
	db.First(&finalCoupon, "name = ?", couponName)
	if finalCoupon.RemainingAmount != 0 {
		t.Errorf("Expected remaining stock 0, got %d", finalCoupon.RemainingAmount)
	}

	var totalClaims int64
	db.Model(&domain.Claim{}).Where("coupon_name = ?", couponName).Count(&totalClaims)
	if totalClaims != int64(stock) {
		t.Errorf("Expected %d claims in DB, got %d", stock, totalClaims)
	}

	t.Logf("Concurrency test passed: Success: %d, NoStock: %d", successCount, noStockCount)
}
