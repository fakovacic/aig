package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/fakovacic/aig/internal/app/config"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/services/controllers"
	"github.com/fakovacic/aig/internal/services/stats"
)

// Ping return new pinger
func Ping(cfg config.Config) {
	pt := PingTask{
		cs:  controllers.New(cfg.Store()),
		ss:  stats.New(cfg.Store()),
		cfg: cfg,
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
				fmt.Println("send ping request to:", list[i].ID)
				model, err := pt.cs.ControllerCommand(ctx, list[i].ID, "stats", "read")
				if err != nil {
					log.Error("command send error", err)
					continue
				}

				_, err = pt.ss.Create(ctx, model.ToStats())
				if err != nil {
					log.Error("stats create error", err)
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
	}()
}

// PingTask task
type PingTask struct {
	cs  *controllers.Service
	ss  *stats.Service
	cfg config.Config
}
