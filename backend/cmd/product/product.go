package product

import (
	"server/internal/transport/temporal"
)

func Run() error {
	err := temporal.StartProduct()
	if err != nil {
		return err
	}
	return nil
}
