package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pambudisatria/scalable-coupon-system/internal/delivery/http"
	"github.com/pambudisatria/scalable-coupon-system/internal/domain"
	"github.com/pambudisatria/scalable-coupon-system/internal/pkg/config"
	repoPostgres "github.com/pambudisatria/scalable-coupon-system/internal/repository/postgres"
	"github.com/pambudisatria/scalable-coupon-system/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ... (load config, connect db, auto-migrate)
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
	couponUsecase := usecase.NewCouponUsecase(db, couponRepo, claimRepo)

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
