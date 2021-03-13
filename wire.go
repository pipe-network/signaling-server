//+build wireinject

package signaling_server

import (
	"github.com/google/wire"
	"github.com/pipe-network/signaling-server/application"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/infrastructure/providers"
	"github.com/pipe-network/signaling-server/interface/controllers"
)

var Providers = wire.NewSet(providers.ProvideUpgrader)

func InitializeMainApplication() application.MainApplication {
	wire.Build(
		Providers,
		services.NewSignalingMessageService,
		wire.Bind(new(services.ISignalingMessageService), new(*services.SignalingMessageService)),
		services.NewSaltyRTCService,
		controllers.NewSignalingController,
		application.NewMainApplication,
	)
	return application.MainApplication{}
}
