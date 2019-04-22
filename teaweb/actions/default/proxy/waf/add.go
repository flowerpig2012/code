package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teawaf"
	"github.com/TeaWeb/code/teawaf/groups"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

type AddAction actions.Action

// 添加策略
func (this *AddAction) RunGet(params struct{}) {
	this.Data["groups"] = lists.Map(groups.InternalGroups, func(k int, v interface{}) interface{} {
		g := v.(*rules.RuleGroup)
		return maps.Map{
			"name": g.Name,
			"code": g.Code,
		}
	})

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	Name       string
	GroupCodes []string
	Must       *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	waf := teawaf.NewWAF()
	waf.Name = params.Name

	for _, groupCode := range params.GroupCodes {
		for _, g := range groups.InternalGroups {
			if g.Code == groupCode {
				newGroup := rules.NewRuleGroup()
				newGroup.Id = stringutil.Rand(16)
				newGroup.On = g.On
				newGroup.Code = g.Code
				newGroup.Name = g.Name
				newGroup.RuleSets = g.RuleSets
				waf.AddRuleGroup(newGroup)
			}
		}
	}

	filename := "waf." + waf.Id + ".conf"
	err := waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	wafList := teaconfigs.SharedWAFList()
	wafList.AddFile(filename)
	err = wafList.Save()
	if err != nil {
		err1 := files.NewFile(Tea.ConfigFile(filename)).Delete()
		if err1 != nil {
			logs.Error(err1)
		}

		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}