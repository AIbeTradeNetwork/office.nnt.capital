package autofarm

import (
	"server/internal/transport/temporal"
)

func Run() error {
	err := temporal.StartAutofarm()
	if err != nil {
		return err
	}
	return nil
}
