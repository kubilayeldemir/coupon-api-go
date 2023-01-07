package services

import (
	"errors"
	"fmt"
	"go-postgresql/internal/models"
	"go-postgresql/internal/repository"
	"log"
)

type CouponService struct {
	couponRepository *repository.CouponRepository
}

func NewCouponService(couponRepository *repository.CouponRepository) *CouponService {
	return &CouponService{couponRepository: couponRepository}
}

func (c CouponService) GetAllCupons() ([]models.Coupon, error) {
	coupons, err := c.couponRepository.GetAllCoupons()
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func (c CouponService) GetCouponById(id string) (models.Coupon, error) {
	coupon, err := c.couponRepository.GetCouponById(id)
	if err != nil {
		return models.Coupon{}, err
	}
	return coupon, nil
}

func (c CouponService) CreateCoupon(model models.CreateNewCouponRequestModel) error {
	err := c.couponRepository.SaveCoupon(models.Coupon{
		Id:       model.Id,
		Name:     model.Name,
		Type:     model.Type,
		Quantity: model.Amount,
	})
	return err
}

func (c CouponService) GiveCouponToUser(couponId, userId string) error {
	coupon, getCouponErr := c.GetCouponById(couponId)
	if getCouponErr != nil {
		return getCouponErr
	}
	if coupon.Quantity <= 0 {
		return errors.New("no coupons left")
	}
	newCouponQuantity := coupon.Quantity - 1
	updateCouponErr := c.couponRepository.UpdateCouponQuantity(couponId)
	if updateCouponErr != nil {
		return updateCouponErr
	}

	couponEventSaveErr := c.couponRepository.SaveCouponGivenEvent(couponId, userId, newCouponQuantity)
	if couponEventSaveErr != nil {
		log.Printf("Couldn't save coupon given event: couponId %s, userId %s, newCouponQuantity %d \n", couponId, userId, newCouponQuantity)
		return nil
	}

	return nil
}

func (c CouponService) GiveCouponToUserWithTransaction(couponId, userId string) error {
	coupon, getCouponErr := c.GetCouponById(couponId)
	if getCouponErr != nil {
		return getCouponErr
	}
	if coupon.Quantity <= 0 {
		return errors.New("no coupons left")
	}
	newCouponQuantity := coupon.Quantity - 1

	couponEventSaveErr := c.couponRepository.GiveCouponToUserAndSaveEventWithTransaction(couponId, userId, newCouponQuantity)
	if couponEventSaveErr != nil {
		log.Printf("Couldn't save coupon given event: couponId %s, userId %s, newCouponQuantity %d \n", couponId, userId, newCouponQuantity)
		return nil
	}

	return nil
}

func (c CouponService) GiveCouponToUserAndSaveEventWithPgFunction(couponId, userId string) error {
	coupon, getCouponErr := c.GetCouponById(couponId)
	if getCouponErr != nil {
		fmt.Println("couldnt get coupon")
		return getCouponErr
	}
	if coupon.Quantity <= 0 {
		fmt.Println("no coupons left")
		return errors.New("no coupons left")
	}

	couponEventSaveErr := c.couponRepository.GiveCouponAndSaveEventWithPgFunction(couponId, userId)
	if couponEventSaveErr != nil {
		log.Printf("GiveCouponAndSaveEventWithPgFunction Couldn't save coupon given event: couponId %s, userId %s, newCouponQuantity %d \n", couponId, userId)
		return couponEventSaveErr
	}
	return nil
}
