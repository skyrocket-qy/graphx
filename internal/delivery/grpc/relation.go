package grpc

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"github.com/skyrocketOoO/zanazibar-dag/domain/delivery/proto"
	usecasedomain "github.com/skyrocketOoO/zanazibar-dag/domain/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	RelationUsecase usecasedomain.RelationUsecase
}

func NewRelationHandler(relationUsecase usecasedomain.RelationUsecase) *GrpcHandler {
	return &GrpcHandler{
		RelationUsecase: relationUsecase,
	}
}

func (h *GrpcHandler) Get(c context.Context, relation *proto.Relation) (*proto.DataResponse, error) {
	requestRelation := domain.Relation{
		ObjectNamespace:  relation.ObjectNamespace,
		ObjectName:       relation.ObjectName,
		Relation:         relation.Relation,
		SubjectNamespace: relation.SubjectNamespace,
		SubjectName:      relation.SubjectName,
		SubjectRelation:  relation.SubjectRelation,
	}
	relations, err := h.RelationUsecase.Get(requestRelation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoRelations := make([]*proto.Relation, len(relations))
	for i, rel := range relations {
		protoRelations[i] = &proto.Relation{
			ObjectNamespace:  rel.ObjectNamespace,
			ObjectName:       rel.ObjectName,
			Relation:         rel.Relation,
			SubjectNamespace: rel.SubjectNamespace,
			SubjectName:      rel.SubjectName,
			SubjectRelation:  rel.SubjectRelation,
		}
	}
	response := &proto.DataResponse{
		Relations: protoRelations,
	}
	return response, nil
}
func (h *GrpcHandler) Create(c context.Context, req *proto.RelationCreateRequest) (*proto.Empty, error) {
	requestRelation := domain.Relation{
		ObjectNamespace:  req.Relation.ObjectNamespace,
		ObjectName:       req.Relation.ObjectName,
		Relation:         req.Relation.Relation,
		SubjectNamespace: req.Relation.SubjectNamespace,
		SubjectName:      req.Relation.SubjectName,
		SubjectRelation:  req.Relation.SubjectRelation,
	}
	err := h.RelationUsecase.Create(requestRelation, req.ExistOk)
	if err != nil {
		if _, ok := err.(domain.CauseCycleError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (h *GrpcHandler) Delete(c context.Context, relation *proto.Relation) (*proto.Empty, error) {
	requestRelation := domain.Relation{
		ObjectNamespace:  relation.ObjectNamespace,
		ObjectName:       relation.ObjectName,
		Relation:         relation.Relation,
		SubjectNamespace: relation.SubjectNamespace,
		SubjectName:      relation.SubjectName,
		SubjectRelation:  relation.SubjectRelation,
	}
	err := h.RelationUsecase.Delete(requestRelation)
	if err != nil {
		if _, ok := err.(domain.CauseCycleError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (h *GrpcHandler) DeleteByQueries(c context.Context, req *proto.DeleteByQueriesRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteByQueries not implemented")
}

func (h *GrpcHandler) BatchOperation(c context.Context, req *proto.BatchOperationRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchOperation not implemented")
}
func (h *GrpcHandler) GetAllNamespaces(c context.Context, empty *proto.Empty) (*proto.StringsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNamespaces not implemented")
}
func (h *GrpcHandler) Check(c context.Context, req *proto.CheckRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (h *GrpcHandler) GetShortestPath(c context.Context, req *proto.GetShortestPathRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortestPath not implemented")
}
func (h *GrpcHandler) GetAllPaths(c context.Context, req *proto.GetAllPathsRequest) (*proto.PathsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPaths not implemented")
}
func (h *GrpcHandler) GetAllObjectRelations(c context.Context, req *proto.GetAllObjectRelationsRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllObjectRelations not implemented")
}
func (h *GrpcHandler) GetAllSubjectRelations(c context.Context, req *proto.GetAllSubjectRelationsRequest) (*proto.DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllSubjectRelations not implemented")
}
func (h *GrpcHandler) ClearAllRelations(c context.Context, empty *proto.Empty) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAllRelations not implemented")
}

func (*GrpcHandler) mustEmbedUnimplementedRelationServiceServer() {
}
