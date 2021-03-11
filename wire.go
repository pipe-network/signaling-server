//+build wireinject

package signaling_server

import (
	"github.com/google/wire"
	"github.com/pipe-network/signaling-server/application"
	"github.com/pipe-network/signaling-server/infrastructure/providers"
	"github.com/pipe-network/signaling-server/interface/controllers"
)

var Providers = wire.NewSet(providers.ProvideUpgrader)

func InitializeMainApplication() application.MainApplication {
	wire.Build(
		Providers,
		controllers.NewSignalingController,
		application.NewMainApplication,
	)
	return application.MainApplication{}
}
