package eventbus

func RegisterEventHandlerWithEventBus[Event any](eventBus EventBus, eventHandler EventHandler[Event]) error {
	return eventBus.Register(eventHandler.Handle)
}

func RegisterEventHandlerFuncWithEventBus[Event any](eventBus EventBus, eventHandlerFunc EventHandlerFunc[*Event]) error {
	return eventBus.Register(eventHandlerFunc)
}
