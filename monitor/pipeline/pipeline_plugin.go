package pipeline

import (
	"fmt"
	"github.com/leilei3167/copy_design_pattern/monitor/config"
	"github.com/leilei3167/copy_design_pattern/monitor/filter"
	"github.com/leilei3167/copy_design_pattern/monitor/input"
	"github.com/leilei3167/copy_design_pattern/monitor/output"
	"github.com/leilei3167/copy_design_pattern/monitor/plugin"
	"reflect"
	"sync/atomic"
)

// Type input插件类型
var Type = make(plugin.Types)

func init() {
	Type["simple"] = reflect.TypeOf(SimplePipeline{}) //两种pipeline的实现
	Type["pool"] = reflect.TypeOf(PoolPipeline{})
}

/*
开闭原则: 一个软件系统应该对拓展开放,对修改关闭,新增功能应当通过拓展的方式实现,而不是修改原有代码
根据具体的业务场景识别出那些最有可能变化的点，然后分离出去，抽象成稳定的接口。
 * 后续新增功能时，通过扩展接口，而不是修改已有代码实现

*/

/*
桥接模式
	桥接模式是将抽象部分与它的实现部分分离，使它们都可以独立地变化
*/

// Plugin pipeline由input、filter、output三种插件组成，定义了一个数据处理流程
// 数据流向为 input -> filter -> output
type Plugin interface {
	plugin.Plugin //继承
	SetInput(input input.Plugin)
	SetFilter(filter filter.Plugin)
	SetOutput(output output.Plugin)
}

/*
此处体现里氏替换原则:
	* 1、基类应该设计为一个抽象类（不能直接实例化，只能被继承）。 此处为plugin.Plugin
 * 2、子类应该实现基类的抽象接口，而不是重写基类已经实现的具体方法。
 * 3、子类可以新增功能，但不能改变基类的功能。此处为各个部分插件,在plugin基础上拓展方法
 * 4、子类不能新增约束，包括抛出基类没有声明的异常。
 * 例子：
 * pipeline.NewPlugin中的入参没有使用plugin.Config作为入参类型，符合LSP。否则就需要转型，破坏了LSP

*/

// NewPlugin 工厂
func NewPlugin(config config.Pipeline) (Plugin, error) {
	//根据配置文件的插件名 读取
	pipelineType, ok := Type[config.PluginType]
	if !ok {
		return nil, plugin.ErrUnknownPlugin
	}
	pipelinePlugin := reflect.New(pipelineType) //反射调用SetContext
	pipelinePlugin.MethodByName("SetContext").Call([]reflect.Value{reflect.ValueOf(config.Ctx)})
	//设置input插件

	//设置filter

	//设置output

}

//桥接?
type pipelineTemplate struct {
	input   input.Plugin
	filter  filter.Plugin
	output  output.Plugin
	isClose uint32
	run     func()
}

func (p *pipelineTemplate) Install() {
	p.output.Install()
	p.filter.Install()
	p.input.Install()
	p.run()
}

func (p *pipelineTemplate) Uninstall() {
	p.input.Uninstall()
	p.filter.Uninstall()
	p.output.Uninstall()
	atomic.StoreUint32(&p.isClose, 1)
}

func (p *pipelineTemplate) SetInput(input input.Plugin) {
	p.input = input
}

func (p *pipelineTemplate) SetFilter(filter filter.Plugin) {
	p.filter = filter
}

func (p *pipelineTemplate) SetOutput(output output.Plugin) {
	p.output = output
}

func (p *pipelineTemplate) doRun() {
	for atomic.LoadUint32(&p.isClose) != 1 {
		//获取输入
		event, err := p.input.Input()
		if err != nil {
			fmt.Printf("pipeline input err %s\n", err.Error())
			atomic.StoreUint32(&p.isClose, 1)
			break
		}
		//过滤
		event = p.filter.Filter(event)

		//输出
		if err = p.output.Output(event); err != nil {
			fmt.Printf("pipeline output err %s\n", err.Error())
			atomic.StoreUint32(&p.isClose, 1)
			break
		}

	}
}
