package manager

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	period = 10 * time.Second
)

func (sm *shardingManager) startPeriodicBulkSyncer(ctx context.Context) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			logrus.Infoln("starting bulk sync")
			err := sm.bulkSync(ctx)
			if err != nil {
				logrus.Errorf("failed to bulk sync: %v", err)
			}
		case <-ctx.Done():
			logrus.Warnf("stopping periodic bulk syncer")
			ticker.Stop()
			return
		}
	}
}
