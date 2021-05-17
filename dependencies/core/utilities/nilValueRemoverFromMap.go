package coreutilitydependencies

import (
	container "github.com/golobby/container/v2"
	coreutilities "github.com/horeekaa/backend/core/utilities"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type NilValueRemoverFromMapDependency struct{}

func (_ *NilValueRemoverFromMapDependency) Bind() {
	container.Singleton(
		func() coreutilityinterfaces.NilValueRemoverFromMapUtility {
			nilValueRemoverUtil, _ := coreutilities.NewNilValueRemoverFromMapUtility()
			return nilValueRemoverUtil
		},
	)
}
