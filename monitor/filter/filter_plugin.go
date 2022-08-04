package filter

import "github.com/leilei3167/copy_design_pattern/monitor/plugin"

// Plugin 过滤插件
type Plugin interface {
	plugin.Plugin
	Filter(event *plugin.Event) *plugin.Event
}
