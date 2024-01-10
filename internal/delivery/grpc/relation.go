package grpc

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain/delivery/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
}

func (h *GrpcHandler) Get(c context.Context, relation *proto.Relation) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (h *GrpcHandler) Create(c context.Context, relation *proto.RelationCreateRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (h *GrpcHandler) Delete(c context.Context, relation *proto.Relation) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (h *GrpcHandler) DeleteByQueries(c context.Context, *proto.DeleteByQueriesRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteByQueries not implemented")
}
func (h *GrpcHandler) BatchOperation(c context.Context, req *proto.BatchOperationRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchOperation not implemented")
}
func (h *GrpcHandler) GetAllNamespaces(c context.Context, *proto.Empty) (*proto.StringsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNamespaces not implemented")
}
func (h *GrpcHandler) Check(c context.Context, req *proto.CheckRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (h *GrpcHandler) GetShortestPath(c context.Context, *proto.GetShortestPathRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortestPath not implemented")
}
func (h *GrpcHandler) GetAllPaths(c context.Context, *proto.GetAllPathsRequest) (*proto.PathsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPaths not implemented")
}
func (h *GrpcHandler) GetAllObjectRelations(c context.Context, *proto.GetAllObjectRelationsRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllObjectRelations not implemented")
}
func (h *GrpcHandler) GetAllSubjectRelations(c context.Context, *proto.GetAllSubjectRelationsRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllSubjectRelations not implemented")
}
func (h *GrpcHandler) ClearAllRelations(c context.Context, *proto.Empty) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAllRelations not implemented")
}
