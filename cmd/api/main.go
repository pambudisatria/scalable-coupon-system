package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pambudisatria/scalable-coupon-system/internal/delivery/http"
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"github.com/pambudisatria/scalable-coupon-system/internal/pkg/config"
	repoPostgres "github.com/pambudisatria/scalable-coupon-system/internal/repository/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// couponUsecaseImpl is a dummy implementation of the usecase interface
type couponUsecaseImpl struct {
	couponRepo domain.CouponRepository
	claimRepo  domain.ClaimRepository
}

func (u *couponUsecaseImpl) ClaimCoupon(userID string, couponName string) error {
	return nil // Placeholder
}

func (u *couponUsecaseImpl) GetCouponStatus(name string) (*domain.Coupon, error) {
	return nil, nil // Placeholder
}

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Connect to Database (Postgres)
	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Auto-migration (GORM)
	err = db.AutoMigrate(&domain.Coupon{}, &domain.Claim{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 4. Initialize Repositories
	couponRepo := repoPostgres.NewCouponRepository(db)
	claimRepo := repoPostgres.NewClaimRepository(db)

	// 5. Initialize Usecases (Manual Wiring / DI)
	couponUsecase := &couponUsecaseImpl{
		couponRepo: couponRepo,
		claimRepo:  claimRepo,
	}

	// 6. Initialize Handlers & Routes
	app := fiber.New()
	couponHandler := http.NewCouponHandler(couponUsecase)
	couponHandler.InitRoutes(app)

	// 7. Start Server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
