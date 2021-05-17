package coreutilitydependencies

import (
	container "github.com/golobby/container/v2"
	coreutilities "github.com/horeekaa/backend/core/utilities"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type MapProcessorUtilityDependency struct{}

func (_ *MapProcessorUtilityDependency) Bind() {
	container.Singleton(
		func() coreutilityinterfaces.MapProcessorUtility {
			mapProcessUtil, _ := coreutilities.NewMapProcessorUtility()
			return mapProcessUtil
		},
	)
}
