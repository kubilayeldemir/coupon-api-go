package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-postgresql/internal/models"
	"go-postgresql/internal/services"
	"net/http"
)

type Handler struct {
	couponService *services.CouponService
}

func NewHandler(e *echo.Echo, couponService *services.CouponService) {
	h := &Handler{
		couponService: couponService,
	}
	g := e.Group("coupon/api/v1")
	g.GET("/:id", h.GetCouponById)
	g.GET("/all", h.GetAllCoupons)
	g.POST("/create", h.CreateNewCoupon)
	g.POST("/give/:couponId", h.GiveCouponToUser)
	g.POST("/give-transaction/:couponId", h.GiveCouponToUserWithTransaction)
	g.POST("/give-transaction-pgfunction/:couponId", h.GiveCouponToUserAndSaveEventWithPgFunction)
}

// GetCouponById example
//
//	@Summary		Get Coupon By Id
//	@Description	Get Coupon By Id
//	@ID				string
//	@Tags			Coupon Api V1
//	@Accept			json
//	@Success		200	{object}	models.Coupon
//	@Failure		400	{object}	models.Response
//	@Failure		500	{object}	models.Response
//	@Produce		json
//	@Param			id	path	string	true	"coupon Id"
//	@Router			/v1/{id} [get]
func (h Handler) GetCouponById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(400, nil)
	}
	coupon, err := h.couponService.GetCouponById(id)
	if err != nil {
		return c.JSON(500, nil)
	}
	if coupon.Id == "" {
		return c.JSON(404, nil)
	}
	return c.JSON(200, coupon)
}

// GetAllCoupons
//
//	@Summary		Get All Coupons
//	@Description	Get All Coupons
//	@Tags			Coupon Api V1
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.Coupon
//	@Failure		400	{object}	models.Response
//	@Failure		500	{object}	models.Response
//	@Router			/v1/all [get]
func (h Handler) GetAllCoupons(c echo.Context) error {
	coupons, _ := h.couponService.GetAllCupons()
	return c.JSON(200, coupons)
}

// CreateNewCoupon
//	@Summary		Create New Coupon
//	@Description	Create New Coupon
//	@Tags			Coupon Api V1
//	@Accept			json
//	@Produce		json
//	@Param			models.CreateNewCouponRequestModel	body		models.CreateNewCouponRequestModel	true	"Create new  coupon request "
//	@Success		200									{object}	models.Response
//	@Failure		400									{object}	models.Response
//	@Failure		500									{object}	models.Response
//	@Router			/v1/create [post]
func (h Handler) CreateNewCoupon(c echo.Context) error {
	createNewCouponRequestModel := new(models.CreateNewCouponRequestModel)

	if bindErr := c.Bind(createNewCouponRequestModel); bindErr != nil {
		panic(bindErr)
	}

	err := h.couponService.CreateCoupon(*createNewCouponRequestModel)
	if err != nil {
		return c.JSON(500, nil)
	}

	return c.String(http.StatusOK, "Hello, World!")
}

// GiveCouponToUser
//	@Summary		Give Coupon To User
//	@Description	Give Coupon To User (basic and wrong way)
//	@Tags			Coupon Api V1
//	@Param			couponId	path	string	true	"coupon Id"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Response
//	@Failure		400	{object}	models.Response
//	@Failure		500	{object}	models.Response
//	@Router			/v1/give/{couponId} [post]
func (h Handler) GiveCouponToUser(c echo.Context) error {
	couponId := c.Param("couponId")
	if couponId == "" {
		return c.JSON(400, nil)
	}
	userId := uuid.NewString()
	err := h.couponService.GiveCouponToUser(couponId, userId)
	if err != nil {
		return c.JSON(500, nil)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Coupon Given To User %s!", userId))
}

// GiveCouponToUserWithTransaction
//	@Summary		Give Coupon To User
//	@Description	Give Coupon To User (correct way with transaction)
//	@Tags			Coupon Api V1
//	@Param			couponId	path	string	true	"coupon Id"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Response
//	@Failure		400	{object}	models.Response
//	@Failure		500	{object}	models.Response
//	@Router			/v1/give-transaction/{couponId} [post]
func (h Handler) GiveCouponToUserWithTransaction(c echo.Context) error {
	couponId := c.Param("couponId")
	if couponId == "" {
		return c.JSON(400, nil)
	}
	userId := uuid.NewString()
	err := h.couponService.GiveCouponToUserWithTransaction(couponId, userId)
	if err != nil {
		return c.JSON(500, nil)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Coupon Given To User %s!", userId))
}

// GiveCouponToUserAndSaveEventWithPgFunction
//	@Summary		Give Coupon To User
//	@Description	Give Coupon To User (correct way with pg function)
//	@Tags			Coupon Api V1
//	@Param			couponId	path	string	true	"coupon Id"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Response
//	@Failure		400	{object}	models.Response
//	@Failure		500	{object}	models.Response
//	@Router			/v1/give-transaction-pgfunction/{couponId} [post]
func (h Handler) GiveCouponToUserAndSaveEventWithPgFunction(c echo.Context) error {
	couponId := c.Param("couponId")
	if couponId == "" {
		return c.JSON(400, nil)
	}
	userId := uuid.NewString()
	err := h.couponService.GiveCouponToUserAndSaveEventWithPgFunction(couponId, userId)
	if err != nil {
		return c.JSON(500, nil)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Coupon Given To User %s!", userId))
}
