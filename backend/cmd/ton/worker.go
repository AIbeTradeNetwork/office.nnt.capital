package ton

import (
	"server/internal/transport/ton"
)

func RunWorker() error {
	err := ton.StartWorker()
	if err != nil {
		return err
	}
	return nil
}
