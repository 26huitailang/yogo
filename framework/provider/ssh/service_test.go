package ssh

import (
	"testing"

	"github.com/26huitailang/yogo/framework/provider/config"
	"github.com/26huitailang/yogo/framework/provider/log"
	tests "github.com/26huitailang/yogo/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestYogoSSHService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.YogoConfigProvider{})
	container.Bind(&log.YogoLogServiceProvider{})

	Convey("test get client", t, func() {
		yogoRedis, err := NewYogoSSH(container)
		So(err, ShouldBeNil)
		service, ok := yogoRedis.(*YogoSSH)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("ssh.web-01"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		session, err := client.NewSession()
		So(err, ShouldBeNil)
		out, err := session.Output("pwd")
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
		session.Close()
	})
}
