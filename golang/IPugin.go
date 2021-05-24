package iplugin

import (
	"context"
	"time"
)

type Metric struct {
	ID        uint64
	Name      string
	Time      time.Time
	Data      interface{}
	AddedDate time.Time
	Tags      map[string]string
}

type IPlugin interface {
	GetNewValues(ctx context.Context, start uint64, limit uint64) ([]Metric, error)
	Close() error
}
