package trafficlight_test

import (
	"context"
	trafficlight "golabs/traffic-light"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	t.Run("", func(t *testing.T) {
		trafficlight.Run(ctx)
	})
}
