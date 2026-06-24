package service

import (
	"context"
	"encoding/json"

	"github.com/gofreego/opengate/api/opengate_v1"
	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/pkg/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetCORSConfig returns the live CORS configuration for use by middleware.
func (s *Service) GetCORSConfig() *utils.CORSConfig {
	return s.settingsMgr.GetCORSConfig()
}

// GetAppSettings implements the gRPC OpenGateServiceServer interface.
func (s *Service) GetAppSettings(_ context.Context, _ *opengate_v1.GetAppSettingsRequest) (*opengate_v1.GetAppSettingsResponse, error) {
	raw := s.settingsMgr.GetAll()
	fields := make(map[string]*structpb.Value, len(raw))
	for k, v := range raw {
		var iface interface{}
		if err := json.Unmarshal(v, &iface); err != nil {
			return nil, err
		}
		protoVal, err := structpb.NewValue(iface)
		if err != nil {
			return nil, err
		}
		fields[k] = protoVal
	}
	return &opengate_v1.GetAppSettingsResponse{
		Settings: &structpb.Struct{Fields: fields},
		Message:  "App settings retrieved successfully",
	}, nil
}

// UpsertAppSetting implements the gRPC OpenGateServiceServer interface.
func (s *Service) UpsertAppSetting(ctx context.Context, req *opengate_v1.UpsertAppSettingRequest) (*opengate_v1.UpsertAppSettingResponse, error) {
	valueBytes, err := req.GetValue().MarshalJSON()
	if err != nil {
		return nil, err
	}
	setting := &models.AppSetting{Key: req.GetKey(), Value: string(valueBytes)}
	if err := s.repo.UpsertAppSetting(ctx, setting); err != nil {
		return nil, err
	}
	s.settingsMgr.Refresh(ctx)
	return &opengate_v1.UpsertAppSettingResponse{Message: "Setting updated successfully"}, nil
}
