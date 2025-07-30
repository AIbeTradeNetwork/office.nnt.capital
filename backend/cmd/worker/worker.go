package worker

import (
	"server/internal/transport/temporal"
)

func Run() error {
	err := temporal.Start()
	if err != nil {
		return err
	}
	return nil
}
