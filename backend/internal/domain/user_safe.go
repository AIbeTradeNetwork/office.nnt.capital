package domain

import "time"

type UserSafe struct {
	UID       string       `bson:"uid" json:"uid"`
	UserUID   string       `bson:"userUid" json:"userUid"`
	SafeUID   string       `bson:"safeUid" json:"safeUid"`
	Code      string       `bson:"code" json:"code"`
	Secret    string       `bson:"secret" json:"secret"`
	CreatedAt time.Time    `bson:"createdAt" json:"createdAt"`
	ClaimedAt time.Time    `bson:"claimedAt" json:"claimedAt"`
	Variant   *SafeVariant `bson:"variant" json:"variant"`
}
