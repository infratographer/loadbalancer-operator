package srv

import "context"

func (s *Server) processLoadBalancerEventUpdate(lb *loadBalancer) error {
	deployed, err := s.hasDeployment(context.TODO(), lb)
	if err != nil {
		return err
	}

	if deployed {
		return s.updateDeployment(lb)
	}

	return nil
}
