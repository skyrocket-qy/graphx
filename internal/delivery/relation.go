package delivery

import (
	"fmt"
	"net/http"
	"zanzibar-dag/domain"
	usecasedomain "zanzibar-dag/domain/usecase"

	"github.com/gin-gonic/gin"
)

type RelationHandler struct {
	RelationUsecase usecasedomain.RelationUsecase
}

func NewRelationHandler(permissionUsecase usecasedomain.RelationUsecase) *RelationHandler {
	return &RelationHandler{
		RelationUsecase: permissionUsecase,
	}
}

// @Summary Get all relations
// @Description Get a list of all relations
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {object} domain.DataResponse
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-all-relations [get]
func (h *RelationHandler) GetAll(c *gin.Context) error {
	relations, err := h.RelationUsecase.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(domain.DataResponse{
		Data: relations,
	})
}

func (h *RelationHandler) Query(c *gin.Context) error

func (h *RelationHandler) Create(c *gin.Context) error

func (h *RelationHandler) Delete(c *gin.Context) error

// @Summary Check if a relation link exists
// @Description Check if a relation link exists between two entities
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {string} string "Relation link exists"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403 {string} string "Relation link does not exist"
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/check [post]
func (h *RelationHandler) Check(c *gin.Context) error {
	type request struct {
		ObjectNamespace  string `json:"object_namespace"`
		ObjectName       string `json:"object_name"`
		Relation         string `json:"relation"`
		SubjectNamespace string `json:"subject_namespace"`
		SubjectName      string `json:"subject_name"`
		SubjectRelation  string `json:"subject_relation"`
	}

	req := request{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrResponse{
			Error: fmt.Sprintf("Parse body error: %s", err.Error()),
		})
	}

	ok, err := h.RelationUsecase.Check(domain.RelationTuple{
		ObjectNamespace:  req.ObjectNamespace,
		ObjectName:       req.ObjectName,
		Relation:         req.Relation,
		SubjectNamespace: req.SubjectNamespace,
		SubjectName:      req.SubjectName,
		SubjectRelation:  req.SubjectRelation,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: err.Error(),
		})
	}
	if ok {
		return nil
	}
	return c.SendStatus(http.StatusForbidden)
}

// @Summary Get the shortest path between two entities in a relation graph
// @Description Get the shortest path between two entities in a relation graph
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {object} domain.DataResponse "Shortest path between entities"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403 {string} string "No path found"
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/path [post]
func (h *RelationHandler) GetShortestPath(c *gin.Context) error {
	type request struct {
		ObjectNamespace  string `json:"object_namespace"`
		ObjectName       string `json:"object_name"`
		Relation         string `json:"relation"`
		SubjectNamespace string `json:"subject_namespace"`
		SubjectName      string `json:"subject_name"`
		SubjectRelation  string `json:"subject_relation"`
	}

	req := request{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrResponse{
			Error: fmt.Sprintf("Parse body error: %s", err.Error()),
		})
	}

	path, err := h.RelationUsecase.GetShortestPath(domain.RelationTuple{
		ObjectNamespace:  req.ObjectNamespace,
		ObjectName:       req.ObjectName,
		Relation:         req.Relation,
		SubjectNamespace: req.SubjectNamespace,
		SubjectName:      req.SubjectName,
		SubjectRelation:  req.SubjectRelation,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: err.Error(),
		})
	}
	if len(path) > 0 {
		type response struct {
			Data []string
		}
		return c.JSON(response{
			Data: path,
		})
	}
	return c.SendStatus(http.StatusForbidden)
}

// @Summary Get the shortest path between two entities in a relation graph
// @Description Get the shortest path between two entities in a relation graph
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {object} domain.DataResponse "Shortest path between entities"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403 {string} string "No path found"
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/path [post]
func (h *RelationHandler) GetAllPaths(c *gin.Context) error {
	type request struct {
		ObjectNamespace  string `json:"object_namespace"`
		ObjectName       string `json:"object_name"`
		Relation         string `json:"relation"`
		SubjectNamespace string `json:"subject_namespace"`
		SubjectName      string `json:"subject_name"`
		SubjectRelation  string `json:"subject_relation"`
	}

	req := request{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrResponse{
			Error: fmt.Sprintf("Parse body error: %s", err.Error()),
		})
	}

	path, err := h.RelationUsecase.GetShortestPath(domain.RelationTuple{
		ObjectNamespace:  req.ObjectNamespace,
		ObjectName:       req.ObjectName,
		Relation:         req.Relation,
		SubjectNamespace: req.SubjectNamespace,
		SubjectName:      req.SubjectName,
		SubjectRelation:  req.SubjectRelation,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: err.Error(),
		})
	}
	if len(path) > 0 {
		type response struct {
			Data []string
		}
		return c.JSON(response{
			Data: path,
		})
	}
	return c.SendStatus(http.StatusForbidden)
}

func (h *RelationHandler) GetObjectRelations(c *gin.Context) error

// @Summary Clear all relations
// @Description Clear all relations in the system
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {string} string "All relations cleared"
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/clear-all-relations [post]
func (h *RelationHandler) ClearAllRelations(c *gin.Context) error {
	err := h.RelationUsecase.ClearAllRelations()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrResponse{
			Error: err.Error(),
		})
	}
	return nil
}
