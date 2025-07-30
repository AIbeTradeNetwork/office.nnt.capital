package ton

import (
	"server/internal/transport/ton"
)

func RunListener() error {
	err := ton.StartListener()
	if err != nil {
		return err
	}
	return nil
}
