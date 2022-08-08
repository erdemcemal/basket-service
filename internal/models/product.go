package models

import "github.com/shopspring/decimal"

type Product struct {
	Base
	Name      string
	UnitPrice decimal.Decimal
	VatRate   int32
	Quantity  int32
}
