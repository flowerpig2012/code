package teadb

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"testing"
	"time"
)

func TestValueDAO_Insert(t *testing.T) {
	dao := SharedDB().ValueDAO()

	{
		value := agents.NewValue()
		value.AgentId = "local"
		value.AppId = "mysql"
		value.SetTime(time.Now())
		value.Value = 4
		err := dao.Insert(value.AgentId, value)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestValueDAO_Insert2(t *testing.T) {
	jsonString := `
{
    "code": 500,
    "message": "\u8bf7\u8f93\u5165\u6b63\u786e\u7684\u4ee4\u724c\uff08001\uff09",
    "data": {},
    "next": null,
    "errors": []
}`
	v := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonString), &v)
	if err != nil {
		t.Fatal(err)
	}
	value := &agents.Value{
		AppId:       "1",
		AgentId:     "1",
		ItemId:      "1",
		Value:       v,
		Error:       "",
		NoticeLevel: notices.NoticeLevelWarning,
	}
	value.SetTime(time.Now())

	err = SharedDB().ValueDAO().Insert("local", value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestValueDAO_ClearItemValues(t *testing.T) {
	dao := SharedDB().ValueDAO()
	err := dao.ClearItemValues("local", "1", "1", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestValuedAO_FindLatestItemValue(t *testing.T) {
	dao := SharedDB().ValueDAO()
	v, err := dao.FindLatestItemValue("local", "system", "cpu.load")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Log("not found")
		return
	}
	logs.PrintAsJSON(v, t)
	t.Log("createdTime:", timeutil.Format("Y-m-d H:i:s", time.Unix(v.CreatedAt, 0)))
}

func TestValuedAO_FindLatestItemValueNoError(t *testing.T) {
	dao := SharedDB().ValueDAO()
	v, err := dao.FindLatestItemValueNoError("local", "system", "cpu.load")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Log("not found")
		return
	}
	logs.PrintAsJSON(v, t)
	t.Log("createdTime:", timeutil.Format("Y-m-d H:i:s", time.Unix(v.CreatedAt, 0)))
}

func TestValueDAO_ListItemValues(t *testing.T) {
	dao := SharedDB().ValueDAO()
	values, err := dao.ListItemValues("local", "system", "cpu.load", 0, "", 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range values {
		t.Log(v.Id, v.Value, v.NoticeLevel)
	}
}

func TestValueDAO_QueryValues(t *testing.T) {
	dao := SharedDB().ValueDAO()
	q := NewQuery("values.agent.local")
	q.Limit(10)
	values, err := dao.QueryValues(q)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range values {
		t.Log(v)
	}
}

func TestValueDAO_DropAgentTable(t *testing.T) {
	dao := SharedDB().ValueDAO()
	err := dao.DropAgentTable("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestValueDAO_GroupValues(t *testing.T) {
	dao := SharedDB().ValueDAO()

	q := NewQuery("values.agent.local").
		Attr("itemId", "cpu.load")

	values, err := dao.GroupValuesByTime(q, "day", map[string]Expr{
		"load1":  NewAvgExpr("value.load1"),
		"load5":  "value.load5",
		"load15": "value.load15",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range values {
		t.Log(v.TimeFormat.Day, v.Value)
	}
}