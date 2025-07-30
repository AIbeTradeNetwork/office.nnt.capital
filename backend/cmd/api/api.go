package api

import (
	"server/internal/transport/graphql/server"
)

func Run() error {
	err := server.Start()
	if err != nil {
		return err
	}
	return nil
}
