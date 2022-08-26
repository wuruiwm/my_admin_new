package util

import (
	"app/global"
	"context"
	"go.uber.org/zap"
	"strings"
	"time"
)

func AdminConfig(key string) string {
	arr := strings.Split(key, ".")
	if len(arr) != 2 {
		return ""
	}
	cacheKey := "admin_config:" + arr[0] + ":" + arr[1]
	val, err := global.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		global.Db.Table("admin_config").
			Where("group", arr[0]).
			Where("key", arr[1]).
			Select("value").
			Scan(&val)
		err = global.Redis.Set(context.Background(), cacheKey, val, 86400*time.Second).Err()
		if err != nil {
			global.Logger.Error("admin_config", zap.String("cacheKey", cacheKey), zap.Any("error", err))
		}
	}
	return val
}
