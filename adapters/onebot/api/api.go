package api

import (
	"context"
	"sync"
	"yora/adapters/onebot/client"
)

var (
	once     sync.Once
	instance *API
)

type API struct {
	client *client.Client
}

func newAPI() *API {
	return &API{client: client.GetClient(context.Background())}
}

func GetAPI() *API {
	once.Do(func() {
		instance = newAPI()
	})
	return instance
}
