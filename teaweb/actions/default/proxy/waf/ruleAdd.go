package waf

import (
	"github.com/TeaWeb/code/teaconfigs"
	wafactions "github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/checkpoints"
	"github.com/TeaWeb/code/teawaf/rules"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type RuleAddAction actions.Action

func (this *RuleAddAction) RunGet(params struct {
	WafId   string
	GroupId string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}
	this.Data["config"] = maps.Map{
		"id":   waf.Id,
		"name": waf.Name,
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	this.Data["group"] = group
	this.Data["connectors"] = []maps.Map{
		{
			"name":        "和 (AND)",
			"value":       rules.RuleConnectorAnd,
			"description": "所有规则都满足才视为匹配",
		},
		{
			"name":        "或 (OR)",
			"value":       rules.RuleConnectorOr,
			"description": "任一规则满足了就视为匹配",
		},
	}

	// check points
	this.Data["checkpoints"] = lists.Map(checkpoints.AllCheckPoints, func(k int, v interface{}) interface{} {
		def := v.(*checkpoints.CheckPointDefinition)
		return maps.Map{
			"name":        def.Name,
			"prefix":      def.Prefix,
			"description": def.Description,
			"hasParams":   def.HasParams,
		}
	})

	this.Data["operators"] = lists.Map(rules.AllRuleOperators, func(k int, v interface{}) interface{} {
		def := v.(*rules.RuleOperatorDefinition)
		return maps.Map{
			"name":        def.Name,
			"code":        def.Code,
			"description": def.Description,
		}
	})

	this.Data["actions"] = lists.Map(wafactions.AllActions, func(k int, v interface{}) interface{} {
		def := v.(*wafactions.ActionDefinition)
		return maps.Map{
			"name":        def.Name,
			"description": def.Description,
			"code":        def.Code,
		}
	})

	this.Show()
}

// 提交测试或者保存
func (this *RuleAddAction) RunPost(params struct {
	WafId   string
	GroupId string

	Name string

	RulePrefixes  []string
	RuleParams    []string
	RuleOperators []string
	RuleValues    []string

	Connector string
	Action    string

	Test         bool
	TestPrefixes []string
	TestParams   []string
	TestValues   []string

	Must *actions.Must
}) {

	set := rules.NewRuleSet()
	set.Name = params.Name
	for index, prefix := range params.RulePrefixes {
		if index < len(params.RuleParams) && index < len(params.RuleOperators) && index < len(params.RuleValues) {
			rule := rules.NewRule()
			rule.Operator = params.RuleOperators[index]

			param := params.RuleParams[index]
			if len(param) > 0 {
				rule.Param = "${" + prefix + "." + param + "}"
			} else {
				rule.Param = "${" + prefix + "}"
			}
			rule.Value = params.RuleValues[index]
			set.AddRule(rule)
		}
	}
	set.Connector = params.Connector
	set.Action = params.Action

	// 测试
	if params.Test {
		err := set.Init()
		if err != nil {
			this.Fail("校验错误：" + err.Error())
		}

		matchedIndex := -1
		breakIndex := -1
		matchLogs := []string{"start matching ...", "==="}
	Loop:
		for index, prefix := range params.TestPrefixes {
			if index < len(params.TestParams) && index < len(params.TestValues) {
				param := ""
				if len(params.TestParams[index]) == 0 {
					param = "${" + prefix + "}"
				} else {
					param = "${" + prefix + "." + params.TestParams[index] + "}"
				}

				breakIndex = index

				for _, rule := range set.Rules {
					if rule.Param == param {
						value := params.TestValues[index]
						if rule.Test(value) {
							matchLogs = append(matchLogs, "rule: "+rule.Param+" "+rule.Operator+" "+rule.Value+"\ncompare: "+value+"\nresult:true")

							if set.Connector == rules.RuleConnectorOr {
								matchedIndex = index
								break Loop
							}

							if set.Connector == rules.RuleConnectorAnd {
								matchedIndex = index
							}
						} else {
							matchLogs = append(matchLogs, "rule: "+rule.Param+" "+rule.Operator+" "+rule.Value+"\ncompare: "+value+"\nresult:false")

							if set.Connector == rules.RuleConnectorAnd {
								matchedIndex = -1
								break Loop
							}
						}
					}
				}
			}
		}

		this.Data["matchedIndex"] = matchedIndex
		this.Data["breakIndex"] = breakIndex
		this.Data["matchLogs"] = matchLogs
		this.Success()
	}

	// 保存
	params.Must.
		Field("name", params.Name).
		Require("请输入规则集名称")

	// waf
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}

	group.AddRuleSet(set)
	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}