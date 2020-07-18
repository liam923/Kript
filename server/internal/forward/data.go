package forward

import (
	"context"
	"github.com/liam923/Kript/server/pkg/proto/kript/api"
)

type dataForwarder struct {
	client api.DataServiceClient
}

func (df *dataForwarder) GetData(ctx context.Context, request *api.GetDataRequest) (*api.GetDataResponse, error) {
	return df.client.GetData(ctx, request)
}

func (df *dataForwarder) UpdateDatum(ctx context.Context, request *api.UpdateDatumRequest) (*api.UpdateDatumResponse, error) {
	return df.client.UpdateDatum(ctx, request)
}

func (df *dataForwarder) CreateDatum(ctx context.Context, request *api.CreateDatumRequest) (*api.CreateDatumResponse, error) {
	return df.client.CreateDatum(ctx, request)
}

func (df *dataForwarder) DeleteDatum(ctx context.Context, request *api.DeleteDatumRequest) (*api.DeleteDatumResponse, error) {
	return df.client.DeleteDatum(ctx, request)
}

func (df *dataForwarder) ShareDatum(ctx context.Context, request *api.ShareDatumRequest) (*api.ShareDatumResponse, error) {
	return df.client.ShareDatum(ctx, request)
}

func NewDataForwarder(client api.DataServiceClient) api.DataServiceServer {
	return &dataForwarder{client: client}
}
