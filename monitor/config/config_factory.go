package config

// Factory 关键点1: 定义抽象工厂接口，里面定义了产品族中各个产品的工厂方法
// 工厂应该具有的方法,此处是创建不同组件的配置选项
type Factory interface {
	CreateInputConfig() Input
	CreateFilterConfig() Filter
	CreateOutputConfig() Output
	CreatePipelineConfig() Pipeline
}

//对于使用方就依赖此处的抽象工厂接口,通过New函数依赖注入(json,yaml)具体的实现
