package idpool

import (
	"context"
	"errors"
	"fmt"
	"github.com/lynx-go/x/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

type idPools struct {
	rdb *redis.Client
	o   Option
}

func (pool *idPools) Put(ctx context.Context, id int64) error {
	err := pool.rdb.SRem(ctx, pool.key(), cast.ToString(id)).Err()
	if err != nil {
		return err
	}
	if pool.o.OnPut != nil {
		pool.o.OnPut(ctx, id)
	}
	return nil
}

func (pool *idPools) key() string {
	return fmt.Sprintf("id_pool:%s", pool.o.IDType)
}

type GenerateFunc func(ctx context.Context) (int64, error)
type OnPutFunc func(ctx context.Context, id int64)

func (pool *idPools) generate(ctx context.Context) (int64, error) {
	return pool.o.Generate(ctx)
}

func (pool *idPools) get(ctx context.Context) (int64, error) {
	key := pool.key()
	nextId, err := pool.generate(ctx)
	if err != nil {
		return 0, err
	}
	ct := 0
	for {
		if ct > pool.o.MaxAttempts {
			return 0, errors.New("生成 ID 失败")
		}
		exists, err := pool.rdb.SIsMember(ctx, key, nextId).Result()
		if err != nil {
			log.ErrorContext(ctx, "查询 ID 是否重复失败", err)
			return 0, err
		}
		if !exists {
			return nextId, nil
		}
		ct += 1
	}
}

// Get 生成编号
func (pool *idPools) Get(ctx context.Context) (int64, error) {
	return pool.get(ctx)
}

type Option struct {
	MaxAttempts int          `json:"maxAttempts"`
	IDType      IDType       `json:"idType"`
	Generate    GenerateFunc `json:"-"`
	OnPut       OnPutFunc    `json:"-"`
}

func (o *Option) Validate() error {
	if o.IDType == "" {
		return errors.New("IDType 不能为空")
	}

	if o.Generate == nil {
		return errors.New("Generate 不能为空")
	}
	return nil
}

func New(rdb *redis.Client, o Option) (IDPool, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return &idPools{rdb: rdb, o: o}, nil
}

func RandomIDGenerateFunc(minId, maxId int64) GenerateFunc {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(ctx context.Context) (int64, error) {
		return r.Int63n(maxId-minId) + minId, nil
	}
}
