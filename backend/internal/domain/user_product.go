package domain

import "time"

type UserProduct struct {
	UID             string
	UserUID         string
	BuyUID          string
	ProductCode     string
	ProductCategory string
	ComboUID        string
	StartAt         time.Time
	EndAt           time.Time
	Priority        int64
	Multiplier      int64
	FlowUID         string
}
