package plugin

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"yora/internal/depends"
	"yora/internal/event"
	"yora/internal/log"
	"yora/internal/matcher"

	"github.com/rs/zerolog"

	onebotevent "yora/protocols/onebot/event"
)

// PluginManager 插件管理器，负责插件的注册、加载和事件分发
type PluginManager struct {
	plugins  map[string]Plugin // 插件映射表，key为插件名称
	logger   zerolog.Logger    // 日志记录器
	mu       sync.RWMutex      // 读写锁，保护并发访问
	baseDeps []depends.Dependent
}

// NewPluginManager 创建新的插件管理器实例
// baseDeps 为插件管理器的基础依赖，可以在handler函数直接使用
func NewPluginManager(baseDeps ...depends.Dependent) *PluginManager {
	if len(baseDeps) == 0 {
		baseDeps = []depends.Dependent{}
	}
	return &PluginManager{
		plugins:  make(map[string]Plugin),
		logger:   log.NewPluginManager("插件管理器"),
		baseDeps: baseDeps,
	}
}

// GetPlugin 根据名称获取插件
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, exists := pm.plugins[name]; exists {
		pm.logger.Debug().Str("插件名", name).Msg("找到插件")
		return p, nil
	}

	pm.logger.Warn().Str("插件名", name).Msg("未找到插件")
	return nil, fmt.Errorf("未找到插件: %s", name)
}

// RegisterPlugin 注册插件到管理器
func (pm *PluginManager) RegisterPlugin(p Plugin) {
	if p == nil {
		panic("无法注册空插件")
	}

	metadata := p.Metadata()
	if metadata.Name == "" {
		panic("插件名称不能为空")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 检查是否已存在同名插件
	if _, exists := pm.plugins[metadata.Name]; exists {
		pm.logger.Warn().
			Str("插件名", metadata.Name).
			Msg("存在同名插件")
		panic("存在同名插件")
	}

	pm.plugins[metadata.Name] = p
	pm.logger.Info().
		Str("插件名", metadata.Name).
		Str("版本", metadata.Version).
		Msg("插件注册成功")

}

// Plugins 返回所有已注册插件的副本
func (pm *PluginManager) Plugins() []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		plugins = append(plugins, p)
	}

	pm.logger.Debug().Int("插件数量", len(plugins)).Msg("获取插件列表")
	return plugins
}

// Dispatch 将事件分发给所有匹配的处理器
// 遍历所有插件的所有匹配器，找到能处理该事件的处理器并执行
// 如果多个处理器出错，返回最后一个错误
func (pm *PluginManager) Dispatch(ctx context.Context, e event.Event) error {
	if ctx == nil {
		return fmt.Errorf("上下文不能为空")
	}
	if e == nil {
		return fmt.Errorf("事件不能为空")
	}

	if ms, ok := e.(*onebotevent.Event); ok {
		pm.logger.Debug().Str("消息内容", ms.ChatID()).Msg("收到消息事件")
	}

	pm.mu.RLock()
	pluginList := make([]Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		pluginList = append(pluginList, p)
	}
	pm.mu.RUnlock()

	pm.logger.Debug().
		Str("事件类型", fmt.Sprintf("%T", e)).
		Int("插件数量", len(pluginList)).
		Msg("分发事件到插件")

	var (
		lastErr      error
		handledCount int
	)

	// 按按优先级排序拿到matcher
	matchers := make([]*matcher.Matcher, 0, len(pluginList)*5)
	for _, p := range pluginList {
		matchers = append(matchers, p.Matchers()...)
	}
	sort.Slice(matchers, func(i, j int) bool {
		return matchers[i].Priority < matchers[j].Priority
	})

	// 遍历所有插件
	for _, m := range matchers {
		if m.Match(ctx, e) {
			pm.logger.Debug().
				Str("匹配器", fmt.Sprintf("%T", m)).
				Msg("匹配器可以处理事件，执行处理器")

			if err := m.Handle(ctx, e, pm.baseDeps...); err != nil {
				pm.logger.Error().
					Err(err).
					Str("匹配器", fmt.Sprintf("%T", m)).
					Msg("处理器执行失败")
				lastErr = err
			} else {
				if m.Block {
					pm.logger.Debug().
						Str("匹配器", fmt.Sprintf("%T", m)).
						Msg("处理器执行成功，阻止事件传播")
					break
				}
				handledCount++
				pm.logger.Debug().
					Str("匹配器", fmt.Sprintf("%T", m)).
					Msg("处理器执行成功")
			}
		}
	}

	pm.logger.Info().
		Int("已处理", handledCount).
		Int("总匹配器", len(matchers)).
		Bool("有错误", lastErr != nil).
		Msg("事件分发完成")

	return lastErr
}

// UnregisterPlugin 注销指定名称的插件
func (pm *PluginManager) UnregisterPlugin(name string) error {
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
