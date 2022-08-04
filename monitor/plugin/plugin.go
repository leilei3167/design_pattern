package plugin

// Plugin 插件接口，在监控系统中，一切皆为插件
type Plugin interface {
	// Install 安装插件，只有在安装之后才能运行
	Install()
	// Uninstall 卸载插件，卸载后停止运行
	Uninstall()
	// SetContext 插件上下文设置
	SetContext(ctx Context)
}

// Event 定义插件间的通信事件,建造者模式逐步构建
type Event struct {
	headers map[string]string
	payload any
}

func NewEvent(payload any) *Event {
	return &Event{
		headers: make(map[string]string),
		payload: payload,
	}
}

func (e *Event) AddHeader(key, value string) *Event {
	e.headers[key] = value
	return e
}

func (e *Event) Payload() any {
	return e.payload
}

func (e *Event) Header(key string) (string, bool) {
	val, ok := e.headers[key]
	return val, ok
}
