package mavbake

import (
	"github.com/mavryk-network/mavpeak/core/common"
)

var (
	rightsEventSource = common.NewEventSource[*RightsStatus](func(rights *RightsStatus) bool {
		return rights.Level > lastProcessedRightsHeight
	})
	lastProcessedRightsHeight = int64(0)
)

func init() {
	go rightsEventSource.Run()
}
