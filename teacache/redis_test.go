package teacache

import (
	"github.com/TeaWeb/code/teatesting"
	"testing"
	"time"
)

func TestRedisManager(t *testing.T) {
	if !teatesting.RequireRedis() {
		return
	}

	manager := NewRedisManager()
	manager.Life = 30 * time.Second
	manager.SetOptions(map[string]interface{}{
		"host": "127.0.0.1",
	})

	t.Log(manager.Write("hello", []byte("world")))
	r, err := manager.Read("hello")
	if err != nil {
		t.Fatal("err:", err)
	} else {
		t.Log("read:", string(r))
	}
}

func TestRedisManager_Stat(t *testing.T) {
	if !teatesting.RequireRedis() {
		return
	}

	manager := NewRedisManager()
	manager.SetId("abc")
	manager.Life = 1800 * time.Second
	manager.SetOptions(map[string]interface{}{
		"network": "tcp",
		"host":    "127.0.0.1",
		"port":    "6379",
	})
	t.Log(manager.Write("key1", []byte("value1")))
	t.Log(manager.Write("key2", []byte("value1")))
	t.Log(manager.Write("key3", []byte("value1")))
	t.Log(manager.Stat())
}

func TestRedisManager_Clean(t *testing.T) {
	if !teatesting.RequireRedis() {
		return
	}

	manager := NewRedisManager()
	manager.SetId("abc")
	manager.Life = 1800 * time.Second
	manager.SetOptions(map[string]interface{}{
		"network": "tcp",
		"host":    "127.0.0.1",
		"port":    "6379",
	})
	t.Log(manager.Clean())
}
