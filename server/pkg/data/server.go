package data

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

type Server struct {
}

func (s *Server) GetData(context.Context, *api.GetDataRequest) (*api.GetDataResponse, error) {
	panic("implement me")
}

func (s *Server) UpdateDatum(context.Context, *api.UpdateDatumRequest) (*api.UpdateDatumResponse, error) {
	panic("implement me")
}

func (s *Server) CreateDatum(context.Context, *api.CreateDatumRequest) (*api.CreateDatumResponse, error) {
	panic("implement me")
}

func (s *Server) DeleteDatum(context.Context, *api.DeleteDatumRequest) (*api.DeleteDatumResponse, error) {
	panic("implement me")
}

func (s *Server) ShareDatum(context.Context, *api.ShareDatumRequest) (*api.ShareDatumResponse, error) {
	panic("implement me")
}
