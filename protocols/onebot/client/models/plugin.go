package models

// PluginInfo 插件信息
type Info struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	HomePage    string            `json:"home_page"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}
