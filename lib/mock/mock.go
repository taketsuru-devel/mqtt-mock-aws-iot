package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"mqtt-mock-aws-iot/lib/logger"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Mock struct {
	connection MQTT.Client
	ctx        context.Context
	entity     MockEntity
	logger     logger.LoggerIf
}

func NewMock(ctx context.Context, wg sync.WaitGroup, c MQTT.Client, entity MockEntity, logger logger.LoggerIf) *Mock {
	childCtx, cancelFunc := context.WithCancel(ctx)
	ret := &Mock{
		connection: c,
		ctx:        childCtx,
		logger:     logger,
	}
	wg.Add(1)
	if err := entity.AfterConnect()(ret); err != nil {
		ret.logger.Error(fmt.Sprintf("%s Error on AfterConnect(): %s", entity.GetClientId(), err.Error()))
		cancelFunc()
	}
	if err := entity.PubSubRegister()(ret); err != nil {
		ret.logger.Error(fmt.Sprintf("%s Error on PubSubRegister(): %s\n", entity.GetClientId(), err.Error()))
		cancelFunc()
	}
	//終了処理
	go func(ctx context.Context, clientID string) {
		defer wg.Done()
		defer cancelFunc()
		<-childCtx.Done()
		ret.logger.Debug("AWS IoT Disconnect: " + clientID)
		ret.connection.Disconnect(1000)
		if err := entity.AfterDisconnect()(ret); err != nil {
			ret.logger.Error(fmt.Sprintf("%s Error on AfterDisconnect(): %s", entity.GetClientId(), err.Error()))
		}
	}(ctx, entity.GetClientId())
	return ret
}

//PubSubRegisterでtimerとか使いたい場合
func (m *Mock) GetContext() context.Context {
	return m.ctx
}

func (m *Mock) Subscribe(topic string, callback MQTT.MessageHandler) error {
	m.logger.Debug("accepting to subscribe: " + topic)
	subscribeToken := m.connection.Subscribe(topic, 0, callback)
	subscribeToken.Wait()
	return subscribeToken.Error()
}

//ひとまずQos=0
func (m *Mock) Publish(topic string, message interface{}) error {
	m.logger.Debug("accepting to publish: " + topic)
	messageByte, err := json.Marshal(message)
	if err != nil {
		return err
	}
	token := m.connection.Publish(topic, 0, false, messageByte)
	token.Wait()
	return token.Error()
}
