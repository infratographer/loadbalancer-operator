package srv

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
				EventType:  create,
				SubjectURN: "urn:infratographer:load-balancer:07442309-182E-4498-BC41-EBD679A9A792",
			},
		},
		{
			name:        "create failure",
			expectError: true,
			msg: pubsubx.Message{
				EventType:  create,
				SubjectURN: "thisisnonsense",
			},
		},
		{
			name:        "update message",
			expectError: false,
			msg: pubsubx.Message{
				EventType:  update,
				SubjectURN: "urn:infratographer:load-balancer:07442309-182E-4498-BC41-EBD679A9A792",
			},
		},
		{
			name:        "update failure",
			expectError: true,
			msg: pubsubx.Message{
				EventType:  update,
				SubjectURN: "thisisnonsense",
			},
		},
		{
			name:        "delete message",
			expectError: false,
			msg: pubsubx.Message{
				EventType:  delete,
				SubjectURN: "urn:infratographer:load-balancer:07442309-182E-4498-BC41-EBD679A9A792",
			},
		},
		{
			name:        "delete failure",
			expectError: true,
			msg: pubsubx.Message{
				EventType:  delete,
				SubjectURN: "thisisnonsense",
			},
		},
		{
			name:        "unknown message",
			expectError: true,
			msg: pubsubx.Message{
				EventType:  "unknown",
				SubjectURN: "urn:infratographer:load-balancer:" + uuid.NewString(),
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			s := srv
			err := s.processLoadBalancer(tc.msg)

			if tc.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.Nil(suite.T(), err)
			}

		})
	}
}

func (suite srvTestSuite) TestProcessLoadBalancerCreate() {
	type testCase struct {
		name        string
		msg         pubsubx.Message
		expectError bool
	}

	js := utils.GetJetstreamConnection(suite.NATSServer)

	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "TestProcessLB",
		Subjects: []string{"clb.foo", "clb.bar"},
		MaxBytes: 1024,
	})

	dir, cp, ch, pwd := utils.CreateWorkspace("test-process-lb")
	defer os.RemoveAll(dir)

	srv := Server{
		Gin:             gin.Default(),
		Context:         context.TODO(),
		StreamName:      "TestCreateLB",
		Logger:          zap.NewNop().Sugar(),
		KubeClient:      suite.Kubeenv.Config,
		JetstreamClient: js,
		Debug:           false,
		Prefix:          "clb",
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
				EventType:  create,
				SubjectURN: "urn:infratographer:load-balancer:07442309-182E-4498-BC41-EBD679A9A793",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			s := srv
			err := s.processLoadBalancerCreate(tc.msg)

			if tc.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.Nil(suite.T(), err)
			}

		})
	}

}
