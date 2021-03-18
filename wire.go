//+build wireinject

package signaling_server

import (
	"github.com/google/wire"
	"github.com/pipe-network/signaling-server/application"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/infrastructure/providers"
	"github.com/pipe-network/signaling-server/infrastructure/storages"
	"github.com/pipe-network/signaling-server/interface/controllers"
)

var Providers = wire.NewSet(
	providers.ProvideUpgrader,
)

var FlagProviders = wire.NewSet(
	providers.ProvideServerAddress,
	providers.ProvidePublicKeyPath,
	providers.ProvidePrivateKeyPath,
	providers.ProvideTLSCertFilePath,
	providers.ProvideTLSKeyFilePath,
)

func InitializeMainApplication() (application.MainApplication, error) {
	panic(
		wire.Build(
			Providers,
			FlagProviders,
			storages.NewKeyPairLocalStorageAdapter,
			services.NewSaltyRTCService,
			controllers.NewSignalingController,
			application.NewMainApplication,

			wire.Bind(new(ports.KeyPairStoragePort), new(*storages.KeyPairLocalStorageAdapter)),
		),
	)
}
