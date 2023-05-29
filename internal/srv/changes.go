package srv

func (s *Server) processLoadBalancerChangeCreate(lb *loadBalancer) error {
	if err := s.newDeployment(lb); err != nil {
		s.Logger.Errorw("handler unable to create loadbalancer", "error", err)
		return err
	}

	return nil
}

func (s *Server) processLoadBalancerChangeDelete(lb *loadBalancer) error {
	if err := s.removeDeployment(lb); err != nil {
		s.Logger.Errorw("handler unable to delete loadbalancer", "error", err, "loadBalancer", lb.loadBalancerID.String())
		return err
	}

	return nil
}

func (s *Server) processLoadBalancerChangeUpdate(lb *loadBalancer) error {
	return nil
}
