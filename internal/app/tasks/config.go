package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/fakovacic/aig/internal/app/config"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/services/controllers"
)

// Config return new pinger
func Config(cfg config.Config) {
	pt := ConfigTask{
		cs: controllers.New(cfg.Store()),
	}

	ctx := context.Background()

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			list, err := pt.cs.List(ctx, meta.New())
			if err != nil {
				log.Error("list of controllers not found", err)
				continue
			}

			if len(list) == 0 {
				log.Error("controllers not found")
				continue
			}

			for i := range list {
				commandList := list[i].ValidCommands()

				if len(commandList) == 0 {
					continue
				}

				for control, power := range commandList {
					fmt.Printf("send command %s %s to %s \n", control, power, list[i].ID)
					model, err := pt.cs.ControllerCommand(ctx, list[i].ID, control, power)
					if err != nil {
						log.Error("command send error", err)
						continue
					}

					ctx = context.WithValue(ctx, "external_id", model.ID)
					_, err = pt.cs.Update(ctx, model)
					if err != nil {
						log.Error("update controller error", err)
						continue
					}
				}
			}
		}
	}()
}

// ConfigTask task
type ConfigTask struct {
	cs *controllers.Service
}
