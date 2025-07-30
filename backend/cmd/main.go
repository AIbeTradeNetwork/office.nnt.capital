package main

import (
	"log"
	"server/cmd/api"
	"server/cmd/autofarm"
	"server/cmd/notifier"
	"server/cmd/product"
	"server/cmd/ton"
	"server/cmd/worker"
	"server/internal/config"
	"sync"
)

func main() {
	cfg := config.Get()
	wg := sync.WaitGroup{}
	var err error

	if cfg.ApiRole == "enabled" {
		wg.Add(1)
		go func() {
			err = api.Run()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.WorkerRole == "enabled" {
		wg.Add(1)
		go func() {
			err = worker.Run()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.NotifierRole == "enabled" {
		wg.Add(1)
		go func() {
			err = notifier.Run()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.TonListenerRole == "enabled" {
		wg.Add(1)
		go func() {
			err = ton.RunListener()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.TonWorkerRole == "enabled" {
		wg.Add(1)
		go func() {
			err = ton.RunWorker()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.ProductRole == "enabled" {
		wg.Add(1)
		go func() {
			err = product.Run()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	if cfg.AutofarmRole == "enabled" {
		wg.Add(1)
		go func() {
			err = autofarm.Run()
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
