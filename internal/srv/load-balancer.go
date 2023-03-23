package srv

import (
	"github.com/spf13/viper"
	"go.infratographer.com/x/pubsubx"
	"go.infratographer.com/x/urnx"

	events "go.infratographer.com/loadbalanceroperator/pkg/events/v1alpha1"
)

func (s *Server) processLoadBalancer(msg pubsubx.Message) error {
	switch msg.EventType {
	case create:
		if err := s.processLoadBalancerCreate(msg); err != nil {
			return err
		}
	case update:
		if err := s.processLoadBalancerUpdate(msg); err != nil {
			return err
		}
	case delete:
		if err := s.processLoadBalancerDelete(msg); err != nil {
			return err
		}
	default:
		s.Logger.Errorw("Unknown action: %s", "action", msg.EventType)
		return errUnknownEventType
	}

	return nil
}

func (s *Server) processLoadBalancerCreate(msg pubsubx.Message) error {
	lbdata := events.LoadBalancerData{}

	lbURN, err := urnx.Parse(msg.SubjectURN)
	if err != nil {
		s.Logger.Errorw("unable to parse load-balancer URN", "error", err)
		return err
	}

	if err := s.parseLBData(&msg.AdditionalData, &lbdata); err != nil {
		s.Logger.Errorw("handler unable to parse loadbalancer data", "error", err)
		return err
	}

	if _, err := s.CreateNamespace(lbURN.ResourceID.String()); err != nil {
		s.Logger.Errorw("handler unable to create required namespace", "error", err)
		return err
	}

	overrides := []valueSet{}
	for _, cpuFlag := range viper.GetStringSlice("helm-cpu-flag") {
		overrides = append(overrides, valueSet{
			helmKey: cpuFlag,
			value:   lbdata.Resources.CPU,
		})
	}

	for _, memFlag := range viper.GetStringSlice("helm-memory-flag") {
		overrides = append(overrides, valueSet{
			helmKey: memFlag,
			value:   lbdata.Resources.Memory,
		})
	}

	if err := s.newDeployment(lbURN.ResourceID.String(), overrides); err != nil {
		s.Logger.Errorw("handler unable to create loadbalancer", "error", err)
		return err
	}

	return nil
}

func (s *Server) processLoadBalancerDelete(msg pubsubx.Message) error {
	lbURN, err := urnx.Parse(msg.SubjectURN)
	if err != nil {
		s.Logger.Errorw("unable to parse load-balancer URN", "error", err)
		return err
	}

	if err := s.removeDeployment(lbURN.ResourceID.String()); err != nil {
		s.Logger.Errorw("handler unable to delete loadbalancer", "error", err)
		return err
	}

	return nil
}

func (s *Server) processLoadBalancerUpdate(msg pubsubx.Message) error {
	return nil
}
