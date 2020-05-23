package loaders

import (
	"kitty_mon/config"
	"kitty_mon/reading"
	"kitty_mon/util"
)

func ConfigLoader() {
	config.Opts = config.NewOpts()

	if config.Opts.V {
		util.Fpl(config.App_name, config.Version)
		util.Fpl(reading.CatTemp())
		util.Fpl("Local IPs:", util.IPs(false))
		return
	}
}
