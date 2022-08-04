package monitor

import (
	"fmt"
	"github.com/leilei3167/copy_design_pattern/monitor/config"
	"github.com/leilei3167/copy_design_pattern/monitor/plugin"
)

/*
监控系统模块,模拟实现服务日志手机,分析储存等功能
	主要用到
	- 桥接模式
	- 策略模式


*/

type System struct {
	plugins       map[string]plugin.Plugin //监控系统中一切皆插件
	configFactory config.Factory           //json或yaml的配置文件解析接口实例
}

func NewSystem(configFactory config.Factory) *System { //在调用处进行依赖注入 配置文件解析
	return &System{
		plugins:       make(map[string]plugin.Plugin),
		configFactory: configFactory,
	}
}

func (s *System) LoadConf(conf string) error {
	pipelineConf := s.configFactory.CreatePipelineConfig()
	if err := pipelineConf.Load(conf); err != nil {
		return err
	}
	//创建pipeline插件

	return nil
}

func (s *System) Start() {
	//依次安装插件
	for name, p := range s.plugins {
		p.Install()
		fmt.Printf("plugin:%v install success!\n", name)
	}
}

func (s *System) ShutDown() {
	//依次卸载插件
	for name, p := range s.plugins {
		p.Uninstall()
		fmt.Printf("plugin:%v uninstall success!\n", name)
	}

}
