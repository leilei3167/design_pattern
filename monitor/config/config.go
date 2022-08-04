package config

import "github.com/leilei3167/copy_design_pattern/monitor/plugin"

/*抽象工厂 用于生产某一产品族*/

//配置各个产品对象

// 配置基础结构
type item struct {
	Name       string                                    `json:"name" yaml:"name"`
	PluginType string                                    `json:"type" yaml:"type"`
	Ctx        plugin.Context                            `json:"context" yaml:"context"`
	loadConf   func(conf string, item interface{}) error // 封装不同配置文件的加载逻辑，实现多态的关键
}

// Input配置对象
type Input item

func (i *Input) Load(conf string) error {
	return i.loadConf(conf, i)
}

// Filter配置对象
type Filter item

func (f *Filter) Load(conf string) error {
	return f.loadConf(conf, f)
}

// Output配置对象
type Output item

func (o *Output) Load(conf string) error {
	return o.loadConf(conf, o)
}

// Pipeline配置对象
type Pipeline struct {
	item    `yaml:",inline"` // yaml嵌套时需要加上,inline
	Input   Input            `json:"input" yaml:"input"`
	Filters []Filter         `json:"filters" yaml:"filters,flow"`
	Output  Output           `json:"output" yaml:"output"`
}

func (p *Pipeline) Load(conf string) error {
	return p.loadConf(conf, p)
}
