package tealogs

import (
	"github.com/TeaWeb/code/teaevents"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		accessLogger = NewAccessLogger()

		teaevents.On(teaevents.EventTypeReload, func(event teaevents.EventInterface) {
			ResetAllPolicies()
		})
	})
}
