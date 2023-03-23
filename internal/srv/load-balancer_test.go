package srv

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/loadbalanceroperator/internal/utils"
	"go.infratographer.com/x/pubsubx"
	"go.uber.org/zap"
)

func (suite srvTestSuite) TestProcessLoadBalancer() {
	type testCase struct {
		name        string
		msg         pubsubx.Message
		expectError bool
	}

	js := utils.GetJetstreamConnection(suite.NATSServer)

	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "TestProcessLB",
		Subjects: []string{"plb.foo", "plb.bar"},
		MaxBytes: 1024,
	})

	dir, cp, ch, pwd := utils.CreateWorkspace("test-process-lb")
	defer os.RemoveAll(dir)

	srv := Server{
		Gin:             gin.Default(),
		Context:         context.TODO(),
		StreamName:      "TestProcessLB",
		Logger:          zap.NewNop().Sugar(),
		KubeClient:      suite.Kubeenv.Config,
		JetstreamClient: js,
		Debug:           false,
		Prefix:          "plb",
		Subjects:        []string{"foo", "bar"},
		Subscriptions:   []*nats.Subscription{},
		Chart:           ch,
		ChartPath:       cp,
		ValuesPath:      pwd + "/../../hack/ci/values.yaml",
	}

	testCases := []testCase{
		{
			name:        "create message",
			expectError: false,
			msg: pubsubx.Message{
				EventType: create,
			},
		},
		{
			name:        "update message",
			expectError: false,
			msg: pubsubx.Message{
				EventType: update,
			},
		},
		{
			name:        "delete message",
			expectError: false,
			msg: pubsubx.Message{
				EventType: delete,
			},
		},
		{
			name:        "unknown message",
			expectError: true,
			msg: pubsubx.Message{
				EventType: "unknown",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			s := srv
			err := s.processLoadBalancer(tc.msg)
			assert.Nil(suite.T(), err)
		})
	}
}

func (suite srvTestSuite) TestLoadBalancerUpdate() {
	s := Server{}
	msg := pubsubx.Message{}
	err := s.processLoadBalancerUpdate(msg)
	assert.Nil(suite.T(), err)
}
