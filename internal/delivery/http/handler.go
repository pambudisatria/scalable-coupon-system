package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pambudisatria/scalable-coupon-system/internal/usecase"
)

// CouponHandler handles HTTP requests for coupons
type CouponHandler struct {
	usecase usecase.CouponUsecase
}

// NewCouponHandler initializes a new handler
func NewCouponHandler(u usecase.CouponUsecase) *CouponHandler {
	return &CouponHandler{usecase: u}
}

// InitRoutes sets up the routes for the coupon handler
func (h *CouponHandler) InitRoutes(app *fiber.App) {
	app.Post("/coupons/claim", h.Claim)
	app.Get("/coupons/:name", h.GetStatus)
}

// Claim is a placeholder for the claim endpoint
func (h *CouponHandler) Claim(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Claim logic not implemented",
	})
}

// GetStatus is a placeholder for the status endpoint
func (h *CouponHandler) GetStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "GetStatus logic not implemented",
	})
}
