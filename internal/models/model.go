package models

type CreateNewCouponRequestModel struct {
	Id     string
	Name   string
	Type   string
	Amount int
}

type Coupon struct {
	Id       string `json:"id"`
	Name     string
	Type     string
	Quantity int
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"errorDescription"`
}
