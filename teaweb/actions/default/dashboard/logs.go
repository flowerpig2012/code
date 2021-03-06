package dashboard

import (
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type LogsAction actions.Action

// 实时日志
func (this *LogsAction) Run(params struct{}) {
	ones, err := teadb.AccessLogDAO().ListTopAccessLogs(timeutil.Format("Ymd"), 10)
	if err != nil {
		if err != teadb.ErrorDBUnavailable {
			logs.Error(err)
		}
		this.Data["logs"] = []*accesslogs.AccessLog{}
	} else {
		this.Data["logs"] = ones
	}

	this.Data["qps"] = teaproxy.QPS

	this.Success()
}
