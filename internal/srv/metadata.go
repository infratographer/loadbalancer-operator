package srv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"

	metastatus "go.infratographer.com/load-balancer-api/pkg/metadata"
	metacli "go.infratographer.com/metadata-api/pkg/client"

	"go.infratographer.com/load-balancer-operator/internal/config"
)

// LoadBalancerStatusUpdate updates the state of a load balancer in the metadata service
func (s Server) LoadBalancerStatusUpdate(ctx context.Context, loadBalancerID gidx.PrefixedID, status *metastatus.LoadBalancerStatus) error {
	if config.AppConfig.Metadata.Endpoint == "" {
		s.Logger.Warnln("metadata not configured")
		return nil
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return err
	}

	if err := s.updateMetering(ctx, loadBalancerID, status); err != nil {
		return err
	}

	if _, err := s.MetadataClient.StatusUpdate(ctx, &metacli.StatusUpdateInput{
		NodeID:      loadBalancerID.String(),
		NamespaceID: config.AppConfig.Metadata.StatusNamespaceID.String(),
		Source:      config.AppConfig.Metadata.Source,
		Data:        json.RawMessage(jsonBytes),
	}); err != nil {
		return err
	}

	return nil
}

func (s Server) updateMetering(ctx context.Context, loadBalancerID gidx.PrefixedID, status *metastatus.LoadBalancerStatus) error {
	if s.MeteringSubject == "" {
		s.Logger.Warnln("metering subject not configured")
		return nil
	}

	if status.State == metastatus.LoadBalancerStateDeleted || status.State == metastatus.LoadBalancerStateActive {
		eventType := "metadata.status"
		if status.State == metastatus.LoadBalancerStateDeleted {
			eventType += ".deleted"
		} else {
			eventType += ".active"
		}

		msg := events.ChangeMessage{
			EventType: eventType,
			SubjectID: loadBalancerID,
			Timestamp: time.Now().UTC(),
		}
		if _, err := s.EventsConnection.PublishChange(ctx, s.MeteringSubject, msg); err != nil {
			return fmt.Errorf("failed to publish change: %w", err)
		}
	}

	return nil
}
