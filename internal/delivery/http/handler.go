package http

import (
	"errors"

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
func (h *CouponHandler) InitRoutes(app fiber.Router) {
	app.Post("/coupons", h.CreateCoupon)
	app.Post("/coupons/claim", h.Claim)
	app.Get("/coupons/:name", h.GetStatus)
}

// CreateCoupon handles the creation of a new coupon
func (h *CouponHandler) CreateCoupon(c *fiber.Ctx) error {
	var req struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	err := h.usecase.CreateCoupon(req.Name, req.Amount)
	if err != nil {
		if errors.Is(err, usecase.ErrCouponDuplicate) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, usecase.ErrCouponAmountInvalid) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create coupon"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "coupon created successfully",
	})
}

// Claim handles the coupon claim request
func (h *CouponHandler) Claim(c *fiber.Ctx) error {
	var req struct {
		UserID     string `json:"user_id"`
		CouponName string `json:"coupon_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.UserID == "" || req.CouponName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id and coupon_name are required"})
	}

	err := h.usecase.ClaimCoupon(req.UserID, req.CouponName)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyClaimed) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, usecase.ErrCouponNoStock) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, usecase.ErrCouponNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to claim coupon"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "coupon claimed successfully",
	})
}

// GetStatus is a placeholder for the status endpoint
func (h *CouponHandler) GetStatus(c *fiber.Ctx) error {
	name := c.Params("name")
	coupon, err := h.usecase.GetCouponStatus(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "coupon not found"})
	}
	return c.JSON(coupon)
}
