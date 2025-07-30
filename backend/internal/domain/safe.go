package domain

type SafeVariantType string

const (
	SafeVariantTypeUDEX     SafeVariantType = "UDEX"
	SafeVariantTypeAbt      SafeVariantType = "abt"
	SafeVariantTypeUsd      SafeVariantType = "usd"
	SafeVariantTypeVariants SafeVariantType = "variants"
)

type SafeVariant struct {
	Type       SafeVariantType `bson:"type" json:"type"`
	Variants   []SafeVariant   `bson:"variants" json:"variants"`
	Chance     int64           `bson:"chance" json:"chance"`
	Amount     int64           `bson:"amount" json:"amount"`
	Count      int64           `bson:"count" json:"count"`
	LimitCount int64           `bson:"limitCount" json:"limitCount"`
}

type Safe struct {
	UID      string        `bson:"uid" json:"uid"`
	Name     string        `bson:"name" json:"name"`
	Variants []SafeVariant `bson:"variants" json:"variants"`
}
