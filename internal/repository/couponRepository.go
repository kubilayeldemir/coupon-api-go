package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go-postgresql/internal/models"
	"log"
)

type CouponRepository struct {
	DB *sql.DB
}

//            "Username=postgres;Password=1997;Server=localhost;Port=5432;Database=exchange;Trust Server Certificate=true;";
const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "1997"
	DBNAME   = "couponDb"
)

func NewCouponRepository() *CouponRepository {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)
	DB, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	return &CouponRepository{
		DB: DB,
	}
}

func (c *CouponRepository) GetAllCoupons() ([]models.Coupon, error) {
	var coupons []models.Coupon

	rows, err := c.DB.Query("SELECT * FROM COUPON")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var coupon models.Coupon
		if err := rows.Scan(&coupon.Id, &coupon.Name, &coupon.Type, &coupon.Quantity); err != nil {
			log.Fatal(err)
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}

func (c *CouponRepository) GetCouponById(id string) (models.Coupon, error) {
	var coupon models.Coupon

	row := c.DB.QueryRow("SELECT * FROM COUPON WHERE id = $1", id)
	err := row.Scan(&coupon.Id, &coupon.Name, &coupon.Type, &coupon.Quantity)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("no coupon with id %d\n", id)
		return coupon, nil
	case err != nil:
		return coupon, err
		log.Printf("query error: %v\n", err)
	}
	return coupon, nil
}

func (c *CouponRepository) SaveCoupon(coupon models.Coupon) error {
	res, err := c.DB.Exec("INSERT INTO COUPON(Id, Name, Type, Quantity) VALUES($1,$2,$3,$4)", coupon.Id, coupon.Name, coupon.Type, coupon.Quantity)
	if err != nil {
		fmt.Printf("Save coupon error: %v", err)
		return err
	}
	println(res)
	return nil
}

func (c *CouponRepository) UpdateCouponQuantity(couponId string) error {
	res, err := c.DB.Exec(`
		UPDATE coupon
		SET quantity = quantity - 1 
		WHERE id = $1
`, couponId)
	if err != nil {
		fmt.Printf("Save coupon error: %v", err)
		return err
	}
	println(res)
	return nil
}

func (c CouponRepository) SaveCouponGivenEvent(couponId, userId string, newQuantity int) error {
	_, err := c.DB.Exec("INSERT INTO coupon_given_events(CouponId, UserId, NewQuantity) VALUES($1,$2,$3)", couponId, userId, newQuantity)
	if err != nil {
		fmt.Printf("SaveCouponGivenEvent: %v", err)
		return err
	}
	return nil
}

func (c CouponRepository) GiveCouponToUserAndSaveEventWithTransaction(couponId, userId string, newQuantity int) error {
	tx, err := c.DB.Begin()

	if err != nil {
		fmt.Printf("SaveCouponGivenEvent Transaction error: %v \n", err)
		return err
	}

	couponRow := tx.QueryRow(`
		SELECT * FROM coupon WHERE id = $1 FOR UPDATE;
		`, couponId)
	var coupon models.Coupon

	couponRowErr := couponRow.Scan(&coupon.Id, &coupon.Name, &coupon.Type, &coupon.Quantity)

	switch {
	case couponRowErr == sql.ErrNoRows:
		_ = tx.Rollback()
		log.Printf("no coupon with id %d\n", couponId)
		return couponRowErr
	case couponRowErr != nil:
		_ = tx.Rollback()
		log.Printf("query error: %v\n", couponRowErr)
		return couponRowErr
	case coupon.Quantity < 1:
		_ = tx.Rollback()
		return errors.New("Quantity problems...")
	}

	_, updateErr := tx.Exec(`
		UPDATE coupon
		SET quantity = quantity - 1 
		WHERE id = $1
		`, couponId)

	if updateErr != nil {
		_ = tx.Rollback()
		fmt.Printf("Update coupon Transaction error: %v \n", updateErr)
		return err
	}

	_, eventErr := tx.Exec("INSERT INTO coupon_given_events(CouponId, UserId, NewQuantity) VALUES($1,$2,$3)", couponId, userId, coupon.Quantity-1)
	if eventErr != nil {
		fmt.Printf("SaveCouponGivenEvent: %v \n", err)
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c CouponRepository) GiveCouponAndSaveEventWithPgFunction(couponId, userId string) error {
	_, err := c.DB.Exec(`SELECT * FROM give_coupon_to_user($1, $2)`, couponId, userId)
	if err != nil {
		fmt.Printf(" ERROR GiveCouponAndSaveEventWithPgFunction: %v \n", err)
		return err
	}
	return nil
}
