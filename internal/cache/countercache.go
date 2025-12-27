package cache

import "context"

type CounterCache struct {
	*Redis
}

const (
	COUNTER_KEY = "short_url_counter"
)

func NewCounterCache(r *Redis) *CounterCache {
	return &CounterCache{r}
}

func (c *CounterCache) GetCounter(ctx context.Context) (int64, error) {
	return c.Client.Incr(ctx, COUNTER_KEY).Result()
}
