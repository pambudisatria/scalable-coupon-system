package usecase

import (
	"errors"
	"sync"
	"testing"

	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"github.com/pambudisatria/scalable-coupon-system/internal/repository/postgres"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateCoupon_ConcurrentDuplicate(t *testing.T) {
	// 1. Setup DB Connection (skip if Postgres is unavailable)
	dsn := "host=localhost user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable"
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Postgres: %v", err)
	}

	// 2. Clean up and migrate
	db.Exec("DELETE FROM claims")
	db.Exec("DELETE FROM coupons")
	db.AutoMigrate(&domain.Coupon{}, &domain.Claim{})

	// 3. Setup repos and usecase
	couponRepo := postgres.NewCouponRepository(db)
	claimRepo := postgres.NewClaimRepository(db)
	uc := NewCouponUsecase(db, couponRepo, claimRepo)

	// 4. Launch 10 concurrent goroutines all creating the same coupon name
	couponName := "RACE_COUPON_CREATE"
	numRequests := 10
	var wg sync.WaitGroup
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- uc.CreateCoupon(couponName, 5)
		}()
	}

	wg.Wait()
	close(results)

	// 5. Analyze results
	successCount := 0
	duplicateCount := 0
	otherErrorCount := 0

	for err := range results {
		if err == nil {
			successCount++
		} else if errors.Is(err, ErrCouponDuplicate) {
			duplicateCount++
		} else {
			otherErrorCount++
			t.Logf("Unexpected error (must NOT be 500-class): %v", err)
		}
	}

	// 6. Assert: exactly 1 success, rest must be ErrCouponDuplicate, zero unexpected errors
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful creation, got %d", successCount)
	}

	expectedDuplicates := numRequests - 1
	if duplicateCount != expectedDuplicates {
		t.Errorf("Expected %d ErrCouponDuplicate errors, got %d", expectedDuplicates, duplicateCount)
	}

	if otherErrorCount != 0 {
		t.Errorf("Expected 0 unexpected errors (no 500-class), got %d", otherErrorCount)
	}

	// 7. DB verification: exactly 1 row must exist
	var count int64
	db.Model(&domain.Coupon{}).Where("name = ?", couponName).Count(&count)
	if count != 1 {
		t.Errorf("Expected exactly 1 coupon row in DB, got %d", count)
	}

	t.Logf("CreateCoupon concurrency test passed: success=%d, duplicate=%d, unexpected=%d",
		successCount, duplicateCount, otherErrorCount)
}

func TestCreateCoupon_SequentialDuplicate(t *testing.T) {
	// 1. Setup DB Connection (skip if Postgres is unavailable)
	dsn := "host=localhost user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable"
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Postgres: %v", err)
	}

	// 2. Clean up and migrate
	db.Exec("DELETE FROM claims")
	db.Exec("DELETE FROM coupons")
	db.AutoMigrate(&domain.Coupon{}, &domain.Claim{})

	// 3. Setup repos and usecase
	couponRepo := postgres.NewCouponRepository(db)
	claimRepo := postgres.NewClaimRepository(db)
	uc := NewCouponUsecase(db, couponRepo, claimRepo)

	couponName := "SEQUENTIAL_DUPLICATE"

	// 4. First creation must succeed
	err = uc.CreateCoupon(couponName, 10)
	if err != nil {
		t.Fatalf("First CreateCoupon should succeed, got: %v", err)
	}

	// 5. Second creation with same name must return ErrCouponDuplicate
	err = uc.CreateCoupon(couponName, 10)
	if !errors.Is(err, ErrCouponDuplicate) {
		t.Errorf("Expected ErrCouponDuplicate on second call, got: %v", err)
	}

	t.Logf("Sequential duplicate test passed: second call correctly returned ErrCouponDuplicate")
}
