package coreutilitydependencies

import (
	container "github.com/golobby/container/v2"
	coreutilities "github.com/horeekaa/backend/core/utilities"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type StructComparisonDependency struct{}

func (strctCompareDpdcy *StructComparisonDependency) Bind() {
	container.Singleton(
		func() coreutilityinterfaces.StructComparisonUtility {
			structComparisonUtility, _ := coreutilities.NewStructComparisonUtility()
			return structComparisonUtility
		},
	)
}
