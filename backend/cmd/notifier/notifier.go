package notifier

import "server/internal/transport/temporal"

func Run() error {
	err := temporal.StartNotifier()
	if err != nil {
		return err
	}
	return nil
}
