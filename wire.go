//+build wireinject

package signaling_server

import (
	"github.com/google/wire"
	"github.com/pipe-network/signaling-server/application"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/application/services"
	"github.com/pipe-network/signaling-server/infrastructure/database/repositories"
	"github.com/pipe-network/signaling-server/infrastructure/providers"
	infrastructureServices "github.com/pipe-network/signaling-server/infrastructure/services"
	"github.com/pipe-network/signaling-server/infrastructure/storages"
	"github.com/pipe-network/signaling-server/interface/controllers"
)

var Providers = wire.NewSet(
	providers.ProvideUpgrader,
	providers.DatabaseProvider,
)

func InitializeMainApplication() (application.MainApplication, error) {
	panic(
		wire.Build(
			Providers,
			services.NewFlagServiceImpl,
			infrastructureServices.NewFCMNotificationService,
			services.NewAddDeviceServiceImpl,
			services.NewSaltyRTCServiceImpl,
			repositories.NewDeviceTokenDatabaseRepository,
			storages.NewKeyPairLocalStorageAdapter,
			controllers.NewAddDeviceController,
			controllers.NewSignalingController,
			application.NewMainApplication,

			wire.Bind(new(ports.KeyPairStorage), new(*storages.KeyPairLocalStorageAdapter)),
		),
	)
}
