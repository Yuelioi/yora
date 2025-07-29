package hook

type HookBuilder struct {
	hookType HookType
	handler  HookHandler
	priority HookPriority
	id       string
	once     bool
}

func NewHookBuilder(hookType HookType) *HookBuilder {
	return &HookBuilder{
		hookType: hookType,
		priority: PriorityNormal,
	}
}

func (hb *HookBuilder) Handler(handler HookHandler) *HookBuilder {
	hb.handler = handler
	return hb
}

func (hb *HookBuilder) Priority(priority HookPriority) *HookBuilder {
	hb.priority = priority
	return hb
}

func (hb *HookBuilder) ID(id string) *HookBuilder {
	hb.id = id
	return hb
}

func (hb *HookBuilder) Once() *HookBuilder {
	hb.once = true
	return hb
}

func (hb *HookBuilder) Register(manager *HookManager) string {
	return manager.AddHookWithOptions(hb.hookType, hb.handler, hb.priority, hb.id, hb.once)
}
