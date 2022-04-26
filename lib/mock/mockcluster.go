package mock

import (
	"context"
	"mqtt-mock-aws-iot/lib/connection"
	"mqtt-mock-aws-iot/lib/logger"
	"sync"
	"time"
)

type MockCluster struct {
	ctx        context.Context
	wg         sync.WaitGroup
	PemCa      string
	PemCrt     string
	PemKey     string
	EndPoint   string
	Mocks      []*Mock
	cancelfunc context.CancelFunc
	logger     logger.LoggerIf
}

func NewMockCluster(PemCa string, PemCrt string, PemKey string, EndPoint string, logger logger.LoggerIf) *MockCluster {
	ctx, cancel := context.WithCancel(context.Background())
	return &MockCluster{
		ctx:        ctx,
		wg:         sync.WaitGroup{},
		cancelfunc: cancel,
		PemCa:      PemCa,
		PemCrt:     PemCrt,
		PemKey:     PemKey,
		EndPoint:   EndPoint,
		logger:     logger,
	}
}

func (mc *MockCluster) AddParallel(entities []MockEntity) {
	entitiesNum := len(entities)
	mc.Mocks = make([]*Mock, 0, entitiesNum)
	for i := 0; i < entitiesNum; i++ {
		entity := entities[i]
		mc.logger.Debug("AWS IoT Connect: " + entity.GetClientId())
		c := *connection.GetConnection(entity.GetClientId(), &mc.PemCa, &mc.PemCrt, &mc.PemKey, &mc.EndPoint)
		mc.Mocks = append(mc.Mocks, NewMock(mc.ctx, mc.wg, c, entity, mc.logger))
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func (mc *MockCluster) Close() {
	mc.cancelfunc()
	mc.wg.Wait()
}
