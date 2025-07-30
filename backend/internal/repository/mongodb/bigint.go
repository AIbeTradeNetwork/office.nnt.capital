package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"math/big"
)

type BigInt struct {
	i *big.Int
}

func NewBigInt(bigint *big.Int) *BigInt {
	return &BigInt{i: bigint}
}

func (bi *BigInt) Int() *big.Int {
	return bi.i
}

func (bi *BigInt) MarshalBSONValue() (bsontype.Type, []byte, error) {
	var txt []byte
	var err error
	if bi.i == nil || len(bi.i.Bits()) == 0 {
		txt = []byte("0")
	} else {
		txt, err = bi.i.MarshalText()
		if err != nil {
			return bson.TypeString, nil, err
		}
	}
	return bson.MarshalValue(string(txt))
}

func (bi *BigInt) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	bi.i = big.NewInt(0)
	var v interface{}
	err := bson.UnmarshalValue(bson.TypeString, data, &v)
	if err != nil {
		return fmt.Errorf("error unmarshal bson into value: %w", err)
	}
	err = bi.i.UnmarshalText([]byte(v.(string)))
	if err != nil {
		return fmt.Errorf("error unmarshal string into bigint: %w", err)
	}
	return nil
}
