package plugin

type ErrPlugin uint8

func (e ErrPlugin) Error() string {
	return errMap[e]
}

const (
	ErrPluginNotInstalled ErrPlugin = 200 + iota
	ErrPluginUninstalled
	ErrUnknownPlugin
)

var errMap = map[ErrPlugin]string{
	ErrPluginNotInstalled: "plugin not installed yet",
	ErrPluginUninstalled:  "plugin uninstalled",
	ErrUnknownPlugin:      "unknown plugin",
}
