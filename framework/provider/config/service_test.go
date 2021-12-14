package config

import (
	"path/filepath"
	"testing"

	"github.com/26huitailang/yogo/framework/contract"
	tests "github.com/26huitailang/yogo/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestYogoConfig_GetInt(t *testing.T) {
	Convey("test yogo env normal case", t, func() {
		basePath := tests.BasePath
		folder := filepath.Join(basePath, "config")
		serv, err := NewYogoConfig(folder, map[string]string{}, contract.EnvDevelopment)
		So(err, ShouldBeNil)
		conf := serv.(*YogoConfig)
		timeout := conf.GetInt("database.mysql.timeout")
		So(timeout, ShouldEqual, 1)
	})
}
