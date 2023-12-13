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
func (s Server) LoadBalancerStatusUpdate(ctx context.Context, loadBalancerID gidx.PrefixedID, oldStatus *metastatus.LoadBalancerStatus, newStatus *metastatus.LoadBalancerStatus) error {
	if config.AppConfig.Metadata.Endpoint == "" {
		s.Logger.Warnln("metadata not configured")
		return nil
	}

	jsonBytes, err := json.Marshal(newStatus)
	if err != nil {
		return err
	}

	if err := s.updateMetering(ctx, loadBalancerID, oldStatus, newStatus); err != nil {
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

func (s Server) updateMetering(ctx context.Context, loadBalancerID gidx.PrefixedID, oldStatus *metastatus.LoadBalancerStatus, newStatus *metastatus.LoadBalancerStatus) error {
	if s.MeteringSubject == "" {
		s.Logger.Warnln("metering subject not configured")
		return nil
	}

	if newStatus.State == metastatus.LoadBalancerStateDeleted || newStatus.State == metastatus.LoadBalancerStateActive {
		changeset := []events.FieldChange{
			{
				Field:         "metadata_status",
				PreviousValue: string(oldStatus.State),
				CurrentValue:  string(newStatus.State),
			},
		}

		eventType := "metadata.update"

		msg := events.ChangeMessage{
			EventType:    eventType,
			SubjectID:    loadBalancerID,
			Timestamp:    time.Now().UTC(),
			FieldChanges: changeset,
		}
		if _, err := s.EventsConnection.PublishChange(ctx, s.MeteringSubject, msg); err != nil {
			return fmt.Errorf("failed to publish change: %w", err)
		}
	}

	return nil
}
