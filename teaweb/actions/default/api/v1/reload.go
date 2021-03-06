package v1

import (
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type ReloadAction actions.Action

// 重载代理服务
func (this *ReloadAction) RunGet(params struct{}) {
	err := teaproxy.SharedManager.Restart()
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	proxyutils.FinishChange()

	apiutils.SuccessOK(this)
}
