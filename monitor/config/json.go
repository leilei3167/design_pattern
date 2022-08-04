package config

import "encoding/json"

func loadJson(conf string, item any) error { //实现多态
	return json.Unmarshal([]byte(conf), item)
}

type JsonFactory struct {
}

func NewJsonFactory() *JsonFactory {
	return &JsonFactory{}
}

//实现抽象工厂接口

// CreateInputConfig 例子 {"name":"input1", "type":"memory_mq", "context":{"topic":"monitor",...}}
func (j JsonFactory) CreateInputConfig() Input {
	return Input{loadConf: loadJson}
}

// CreateFilterConfig 例子 [{"name":"filter1", "type":"to_json"},{"name":"filter2", "type":"add_timestamp"},...]
func (j JsonFactory) CreateFilterConfig() Filter {
	return Filter{loadConf: loadJson}
}

// CreateOutputConfig 例子 {"name":"output1", "type":"memory_db", "context":{"tableName":"test",...}}
func (j JsonFactory) CreateOutputConfig() Output {
	return Output{loadConf: loadJson}
}

// CreatePipelineConfig 例子 {"name":"pipline1", "type":"simple", "input":{...}, "filter":{...}, "output":{...}}
func (j JsonFactory) CreatePipelineConfig() Pipeline {
	pipeline := Pipeline{}
	pipeline.loadConf = loadJson
	return pipeline
}
