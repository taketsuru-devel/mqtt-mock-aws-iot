package mock

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
  AfterConnect, PubSubRegister, AfterDisconnectはmockのcallback
  AfterConnect, PubSubRegisterはerrorを返すと、mockの処理は一通り実行後即時切断される
*/

type MockEventCallback func(*Mock) error

type MockEntity interface {
	GetClientId() string
	AfterConnect() MockEventCallback
	PubSubRegister() MockEventCallback
	AfterDisconnect() MockEventCallback
}

type MockEntitySample struct {
	ClientId string
}

func (ms *MockEntitySample) GetClientId() string {
	return ms.ClientId
}

func (ms *MockEntitySample) AfterConnect() MockEventCallback {
	return func(m *Mock) error {
		m.logger.Info("after connect: " + ms.ClientId)
		return nil
	}
}

func (ms *MockEntitySample) PubSubRegister() MockEventCallback {
	return func(m *Mock) error {
		ticker := time.NewTicker(10 * time.Second)
		go func() {
			defer ticker.Stop()
			defer m.logger.Debug("exit")
			for {
				select {
				case <-ticker.C:
					if err := m.Publish("test/p/"+ms.ClientId, map[string]string{"test1": "test2"}); err != nil {
						m.logger.Err(err)
					}
				case <-m.GetContext().Done():
					m.logger.Debug("done received")
					return
				}
			}
		}()

		if err := m.Subscribe("test/s/"+ms.ClientId, func(c mqtt.Client, message mqtt.Message) {
			m.logger.Info(string(message.Payload()))
		}); err != nil {
			return err
		}
		return nil
	}
}

func (ms *MockEntitySample) AfterDisconnect() MockEventCallback {
	return func(m *Mock) error {
		m.logger.Debug("disconnected: " + ms.ClientId)
		return nil
	}
}
