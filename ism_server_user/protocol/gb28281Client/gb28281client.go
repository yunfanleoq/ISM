package gb28281client

import (
	"context"
	"os"

	"ISMServer/protocol/gb28281Client/internal/config"

	"ISMServer/protocol/gb28281Client/internal/useragent"
	"ISMServer/protocol/gb28281Client/internal/version"

	cli "github.com/jawher/mow.cli"
	"github.com/qiniu/x/xlog"
)

func Gb28281client() {
	//xlog := xlog.NewWith(context.Background())
	xlog.SetOutputLevel(0)
	xlog.SetFlags(xlog.Llevel | xlog.Llongfile | xlog.Ltime)
	xlog := xlog.NewWith(context.Background())
	app := cli.App("gb28181Simulator", "Runs the gb28181 simulator.")
	app.Spec = "[ -c=<configuration path> ] "
	confPath := app.StringOpt("c config", "sim.conf", "Specifies the configuration path (file) to use for the simulator.")
	app.Action = func() { run(xlog, app, confPath) }

	// Register sub-commands
	app.Command("version", "Prints the version of the executable.", version.Print)
	app.Run(os.Args)
}

func run(xlog *xlog.Logger, app *cli.Cli, conf *string) {
	xlog.Infof("gb28181 simulator is running...")
	cfg, err := config.ParseJsonConfig("conf/sim.conf")
	if err != nil {
		xlog.Errorf("load config file failed, err = %#v", err)
	}
	xlog.Infof("config file = %#v", cfg)
	srv, err := useragent.NewService(xlog, cfg)
	if err != nil {
		xlog.Infof("new service failed err = %#v", err)
		return
	}
	srv.HandleIncommingMsg()
}
