package pluginmanager

import (
	"fmt"
	"reflect"
)

type Plugin interface {
	Run()
}

type MyPlugin1 struct {
	Name string
}

func (p MyPlugin1) Run() {
	fmt.Println("MyPlugin1 ", p.Name, "Running")
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// RunPlugins 执行所有注册的插件
func (pm *PluginManager) RunPlugins() {
	for _, plugin := range pm.plugins {
		// 使用反射调用插件的 Run 方法
		run := reflect.ValueOf(plugin).MethodByName("Run")
		if run.IsValid() {
			run.Call(nil)
		}
	}
}

// RunPlugins 执行所有注册的插件
func (pm *PluginManager) RunPlugin(name string) {
	for _, plugin := range pm.plugins {
		// 使用反射获取插件的name
		if name != reflect.ValueOf(plugin).Elem().FieldByName("Name").String() {
			continue
		}
		run := reflect.ValueOf(plugin).MethodByName("Run")
		if run.IsValid() {
			run.Call(nil)
			break
		}
	}
}
