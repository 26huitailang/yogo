package redis

import (
	"context"
	"github.com/26huitailang/yogo/framework/provider/config"
	"github.com/26huitailang/yogo/framework/provider/log"
	tests "github.com/26huitailang/yogo/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestYogoService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.YogoConfigProvider{})
	container.Bind(&log.YogoLogServiceProvider{})

	Convey("test get client", t, func() {
		yogoRedis, err := NewYogoRedis(container)
		So(err, ShouldBeNil)
		service, ok := yogoRedis.(*YogoRedis)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("redis.write"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		ctx := context.Background()
		err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()
		So(err, ShouldBeNil)
		val, err := client.Get(ctx, "foo").Result()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "bar")
		err = client.Del(ctx, "foo").Err()
		So(err, ShouldBeNil)
	})
}
