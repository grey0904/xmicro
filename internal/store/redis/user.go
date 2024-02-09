package store_redis

import (
	"context"
	"fmt"
	"time"
	"xmicro/internal/app"
	"xmicro/internal/constant"
)

func SetRedisUserToken(ctx context.Context, userId string, token string) error {
	key := fmt.Sprintf(constant.UserTokenKey, userId)
	return app.Redis.Set(ctx, key, token, time.Hour*24).Err()
}
