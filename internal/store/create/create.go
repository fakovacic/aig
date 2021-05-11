package create

import (
	"context"

	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/client"
)

// Tables initialize database tables
func Tables(s *store.Store, sList ...store.Model) {
	ctx := context.TODO()
	for i := range sList {
		uStore := client.New(s, sList[i])
		uStore.Init(ctx)
	}
}
