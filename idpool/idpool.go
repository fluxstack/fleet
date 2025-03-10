package idpool

import "context"

type IDPool interface {
	Get(ctx context.Context) (int64, error)
	Put(ctx context.Context, id int64) error
}

// IDType 序列类型
type IDType string
