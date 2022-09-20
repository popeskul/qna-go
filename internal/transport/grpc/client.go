package grpc

import (
	"context"
	"fmt"

	"github.com/popeskul/audit-logger/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn       *grpc.ClientConn
	grpcClient domain.AuditServiceClient
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) domain.AuditServiceClient {
	return domain.NewAuditServiceClient(conn)
}

func NewClient(host string, port int) (*Client, error) {
	conn, err := setupGrpcConn(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:       conn,
		grpcClient: getUserServiceClient(conn),
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req *domain.LogItem) (*domain.Empty, error) {
	action, err := domain.ToPbAction(req.Action)
	if err != nil {
		return nil, err
	}

	entity, err := domain.ToPbEntity(req.Entity)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.Log(ctx, &domain.LogRequest{
		Action:    action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})
}
