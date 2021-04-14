package coreutilitydependencies

import (
	"github.com/golobby/container/v2"
	coreutilities "github.com/horeekaa/backend/core/utilities"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type StructFieldIteratorDependency struct{}

func (structFieldDpdcy *StructFieldIteratorDependency) bind() {
	container.Singleton(
		func() coreutilityinterfaces.StructFieldIteratorUtility {
			structFieldIteratorUtility, _ := coreutilities.NewStructFieldIteratorUtility()
			return structFieldIteratorUtility
		},
	)
}
