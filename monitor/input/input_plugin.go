package input

import (
	"github.com/leilei3167/copy_design_pattern/monitor/config"
	"github.com/leilei3167/copy_design_pattern/monitor/plugin"
	"reflect"
)

/*
策略模式
*/
var Type = make(plugin.Types)

func init() {
	Type["memory_mq"] = reflect.TypeOf(MemoryMqInput{})
	Type["socket"] = reflect.TypeOf(SocketInput{})
}

type Plugin interface {
	plugin.Plugin
	Input() (*plugin.Event, error)
}

func NewPlugin(config config.Input) (Plugin, error) {
	inputType, ok := Type[config.PluginType]
	if !ok {
		return nil, plugin.ErrUnknownPlugin
	}
	inputPlugin := reflect.New(inputType)
	ctx := reflect.ValueOf(config.Ctx)
	inputPlugin.MethodByName("SetContext").Call([]reflect.Value{ctx})
	return inputPlugin.Interface().(Plugin), nil
}
