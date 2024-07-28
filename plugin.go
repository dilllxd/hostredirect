package hostredirect

import (
	"github.com/dilllxd/hostredirect/internal/plugin"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "HostRedirect",
	Init: plugin.InitPlugin,
}
