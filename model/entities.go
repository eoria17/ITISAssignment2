package model

import "time"

type Menu struct {
	ID       int `gorm:"primary_key"`
	Name     string
	Price    float64
	ImageURL string
}

func (Menu) TableName() string {
	return "menu"
}

type Order struct {
	ID          int `gorm:"primary_key"`
	OrderNumber string
	Date        time.Time `gorm:"type:date"`
	Total       float64
}

func (Order) TableName() string {
	return "order"
}

type OrderLine struct {
	ID              int `gorm:"primary_key"`
	OrderLineNumber string
	MenuID          int
	Menu            *Menu
	Amount          int
	Subtotal        float64
}

func (OrderLine) TableName() string {
	return "order_line"
}
