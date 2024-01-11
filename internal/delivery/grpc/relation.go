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

func (h *GrpcHandler) Get(c context.Context, relation *proto.Relation) (*proto.RelationsResponse, error) {
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
	response := &proto.RelationsResponse{
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
	queries := make([]domain.Relation, len(req.Queries))
	for i, q := range req.Queries {
		queries[i] = domain.Relation{
			ObjectNamespace:  q.ObjectNamespace,
			ObjectName:       q.ObjectName,
			Relation:         q.Relation,
			SubjectNamespace: q.SubjectNamespace,
			SubjectName:      q.SubjectName,
			SubjectRelation:  q.SubjectRelation,
		}
	}
	err := h.RelationUsecase.DeleteByQueries(queries)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (h *GrpcHandler) BatchOperation(c context.Context, req *proto.BatchOperationRequest) (*proto.Empty, error) {
	operations := make([]domain.Operation, len(req.Operations))
	for i, o := range req.Operations {
		operations[i] = domain.Operation{
			Type: domain.Action(o.Type),
			Relation: domain.Relation{
				ObjectNamespace:  o.Relation.ObjectNamespace,
				ObjectName:       o.Relation.ObjectName,
				Relation:         o.Relation.Relation,
				SubjectNamespace: o.Relation.SubjectNamespace,
				SubjectName:      o.Relation.SubjectName,
				SubjectRelation:  o.Relation.SubjectRelation,
			},
		}
	}
	err := h.RelationUsecase.BatchOperation(operations)
	if err != nil {
		if _, ok := err.(domain.CauseCycleError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (h *GrpcHandler) GetAllNamespaces(c context.Context, empty *proto.Empty) (*proto.StringsResponse, error) {
	namespaces, err := h.RelationUsecase.GetAllNamespaces()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	req := proto.StringsResponse{
		Strings: namespaces,
	}
	return &req, nil
}

func (h *GrpcHandler) Check(c context.Context, req *proto.CheckRequest) (*proto.Empty, error) {
	subject := domain.Node{
		Namespace: req.Subject.Namespace,
		Name:      req.Subject.Name,
		Relation:  req.Subject.Relation,
	}
	object := domain.Node{
		Namespace: req.Object.Namespace,
		Name:      req.Object.Name,
		Relation:  req.Object.Relation,
	}
	searchCondition := domain.SearchCondition{
		In: domain.Compare{
			Namespaces: req.SearchCondition.In.Namespaces,
			Names:      req.SearchCondition.In.Name,
			Relations:  req.SearchCondition.In.Relation,
		},
	}
	ok, err := h.RelationUsecase.Check(subject, object, searchCondition)
	if err != nil {
		if _, ok := err.(domain.CauseCycleError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "")
	}
	return nil, nil
}

func (h *GrpcHandler) GetShortestPath(c context.Context, req *proto.GetShortestPathRequest) (*proto.PathResponse, error) {
	subject := domain.Node{
		Namespace: req.Subject.Namespace,
		Name:      req.Subject.Name,
		Relation:  req.Subject.Relation,
	}
	object := domain.Node{
		Namespace: req.Object.Namespace,
		Name:      req.Object.Name,
		Relation:  req.Object.Relation,
	}
	searchCondition := domain.SearchCondition{
		In: domain.Compare{
			Namespaces: req.SearchCondition.In.Namespaces,
			Names:      req.SearchCondition.In.Name,
			Relations:  req.SearchCondition.In.Relation,
		},
	}
	paths, err := h.RelationUsecase.GetShortestPath(subject, object, searchCondition)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoPaths := make([]*proto.Relation, len(paths))
	for i, path := range paths {
		protoPaths[i] = &proto.Relation{
			ObjectNamespace:  path.ObjectNamespace,
			ObjectName:       path.ObjectName,
			Relation:         path.Relation,
			SubjectNamespace: path.SubjectNamespace,
			SubjectName:      path.SubjectName,
			SubjectRelation:  path.SubjectRelation,
		}
	}
	resp := proto.PathResponse{
		Relations: protoPaths,
	}
	return &resp, nil
}

func (h *GrpcHandler) GetAllPaths(c context.Context, req *proto.GetAllPathsRequest) (*proto.PathsResponse, error) {
	subject := domain.Node{
		Namespace: req.Subject.Namespace,
		Name:      req.Subject.Name,
		Relation:  req.Subject.Relation,
	}
	object := domain.Node{
		Namespace: req.Object.Namespace,
		Name:      req.Object.Name,
		Relation:  req.Object.Relation,
	}
	searchCondition := domain.SearchCondition{
		In: domain.Compare{
			Namespaces: req.SearchCondition.In.Namespaces,
			Names:      req.SearchCondition.In.Name,
			Relations:  req.SearchCondition.In.Relation,
		},
	}
	allPaths, err := h.RelationUsecase.GetAllPaths(subject, object, searchCondition)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	paths := make([]*proto.PathResponse, len(allPaths))
	for i, path := range allPaths {
		relations := make([]*proto.Relation, len(path))
		for j, rel := range path {
			relations[j] = &proto.Relation{
				ObjectNamespace:  rel.ObjectNamespace,
				ObjectName:       rel.ObjectName,
				Relation:         rel.Relation,
				SubjectNamespace: rel.SubjectNamespace,
				SubjectName:      rel.SubjectName,
				SubjectRelation:  rel.SubjectRelation,
			}
		}
		paths[i] = &proto.PathResponse{
			Relations: relations,
		}
	}
	resp := proto.PathsResponse{
		Path: paths,
	}
	return &resp, nil
}

func (h *GrpcHandler) GetAllObjectRelations(c context.Context, req *proto.GetAllObjectRelationsRequest) (*proto.RelationsResponse, error) {
	subject := domain.Node{
		Namespace: req.Subject.Namespace,
		Name:      req.Subject.Name,
		Relation:  req.Subject.Relation,
	}
	searchCondition := domain.SearchCondition{
		In: domain.Compare{
			Namespaces: req.SearchCondition.In.Namespaces,
			Names:      req.SearchCondition.In.Name,
			Relations:  req.SearchCondition.In.Relation,
		},
	}
	collectCondition := domain.CollectCondition{
		In: domain.Compare{
			Namespaces: req.CollectCondition.In.Namespaces,
			Names:      req.CollectCondition.In.Name,
			Relations:  req.CollectCondition.In.Relation,
		},
	}
	relations, err := h.RelationUsecase.GetAllObjectRelations(subject, searchCondition, collectCondition, int(req.MaxDepth))
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
	resp := proto.RelationsResponse{
		Relations: protoRelations,
	}
	return &resp, nil
}

func (h *GrpcHandler) GetAllSubjectRelations(c context.Context, req *proto.GetAllSubjectRelationsRequest) (*proto.RelationsResponse, error) {
	object := domain.Node{
		Namespace: req.Object.Namespace,
		Name:      req.Object.Name,
		Relation:  req.Object.Relation,
	}
	searchCondition := domain.SearchCondition{
		In: domain.Compare{
			Namespaces: req.SearchCondition.In.Namespaces,
			Names:      req.SearchCondition.In.Name,
			Relations:  req.SearchCondition.In.Relation,
		},
	}
	collectCondition := domain.CollectCondition{
		In: domain.Compare{
			Namespaces: req.CollectCondition.In.Namespaces,
			Names:      req.CollectCondition.In.Name,
			Relations:  req.CollectCondition.In.Relation,
		},
	}
	relations, err := h.RelationUsecase.GetAllObjectRelations(object, searchCondition, collectCondition, int(req.MaxDepth))
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
	resp := proto.RelationsResponse{
		Relations: protoRelations,
	}
	return &resp, nil
}

func (h *GrpcHandler) ClearAllRelations(c context.Context, empty *proto.Empty) (*proto.Empty, error) {
	err := h.RelationUsecase.ClearAllRelations()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return nil, nil
}

func (*GrpcHandler) mustEmbedUnimplementedRelationServiceServer() {
}
