package plugin

import (
	"errors"
	"fmt"
	"sync"
	"yora/pkg/log"

	"github.com/rs/zerolog"
)

var (
	ponce          sync.Once
	pluginRegistry *PluginRegistry
)

// 插件管理器，负责插件的注册、加载和事件分发
type PluginRegistry struct {
	plugins map[string]Plugin // 插件映射表，key为插件名称
	logger  zerolog.Logger    // 日志记录器
	mu      sync.RWMutex      // 读写锁，保护并发访问
	mr      *MatcherRegistry  // 匹配器管理器
}

func GetPluginRegistry() *PluginRegistry {
	ponce.Do(func() {
		pluginRegistry = newPluginRegistry()
	})
	return pluginRegistry
}

// 创建新的插件管理器实例
func newPluginRegistry() *PluginRegistry {

	return &PluginRegistry{
		plugins: make(map[string]Plugin),
		logger:  log.NewPluginRegistry("插件管理器"),
		mr:      GetMatcherRegistry(),
	}
}

// 根据名称获取插件
func (pm *PluginRegistry) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, exists := pm.plugins[name]; exists {
		pm.logger.Debug().Str("插件名", name).Msg("找到插件")
		return p, nil
	}

	pm.logger.Warn().Str("插件名", name).Msg("未找到插件")
	return nil, fmt.Errorf("未找到插件: %s", name)
}

// 册插件到管理器
func (pm *PluginRegistry) RegisterPlugins(ps ...Plugin) error {

	matchers := make([]*Matcher, 0, len(ps))

	for i, p := range ps {
		if p == nil {
			return errors.New("无法注册空插件")

		}

		// 初始化插件
		if err := p.Init(); err != nil {
			pm.logger.Error().
				Err(err).Int("插件索引", i).Msg("插件初始化失败，跳过")
			return err
		}

		// 加载插件
		if err := p.Load(); err != nil {
			pm.logger.Error().Err(err).Int("插件索引", i).Msg("插件加载失败，跳过")
			return err
		}

		// 获取元数据
		metadata := p.Metadata()
		if metadata.Name == "" {
			panic("插件名称不能为空")
		}

		// 检查是否已存在同名插件
		if _, exists := pm.plugins[metadata.Name]; exists {
			pm.logger.Warn().Str("插件名", metadata.Name).Msg("存在同名插件")
			panic("存在同名插件")
		}

		pm.plugins[metadata.Name] = p
		pm.logger.Info().Str("插件名", metadata.Name).Str("版本", metadata.Version).Msg("插件注册成功")

		for _, m := range p.Matchers() {
			matchers = append(matchers, m)
		}
	}

	// 注册匹配器
	pm.mr.RegisterMatchers(matchers...)

	return nil

}

// Plugins 返回所有已注册插件的副本
func (pm *PluginRegistry) Plugins() []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		plugins = append(plugins, p)
	}

	pm.logger.Debug().Int("插件数量", len(plugins)).Msg("获取插件列表")
	return plugins
}

// UnregisterPlugin 注销指定名称的插件
func (pm *PluginRegistry) UnregisterPlugin(name string) error {
	if name == "" {
		return fmt.Errorf("插件名称不能为空")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.plugins[name]; !exists {
		pm.logger.Warn().Str("插件名", name).Msg("尝试注销不存在的插件")
		return fmt.Errorf("未找到插件: %s", name)
	}

	delete(pm.plugins, name)
	pm.logger.Info().Str("插件名", name).Msg("插件注销成功")

	return nil
}
