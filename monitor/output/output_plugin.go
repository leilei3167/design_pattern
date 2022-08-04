package output

import "github.com/leilei3167/copy_design_pattern/monitor/plugin"

// Plugin 输出插件
type Plugin interface {
	plugin.Plugin
	Output(event *plugin.Event) error
}
