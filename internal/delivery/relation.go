package delivery

import (
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
func (h *RelationHandler) GetAll(c *gin.Context) {
	relations, err := h.RelationUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: relations,
	})
}

func (h *RelationHandler) Query(c *gin.Context) {
	query := domain.Relation{
		ObjectNamespace:  c.Query("object-namespace"),
		ObjectName:       c.Query("object-name"),
		Relation:         c.Query("relation"),
		SubjectNamespace: c.Query("subject-namespace"),
		SubjectName:      c.Query("subject-name"),
		SubjectRelation:  c.Query("subject-relation"),
	}

	relations, err := h.RelationUsecase.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}

	type Response struct {
		Data []domain.Relation `json:"data"`
	}
	c.JSON(http.StatusOK, Response{
		Data: relations,
	})
}

func (h *RelationHandler) Create(c *gin.Context) {
	relation := domain.Relation{}
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.Create(relation); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RelationHandler) Delete(c *gin.Context) {
	relation := domain.Relation{}
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.Delete(relation); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RelationHandler) GetAllNamespaces(c *gin.Context) {
	namespaces, err := h.RelationUsecase.GetAllNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	type response struct {
		Data []string `json:"data"`
	}
	c.JSON(http.StatusOK, response{
		Data: namespaces,
	})
}

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
func (h *RelationHandler) Check(c *gin.Context) {
	relation := domain.Relation{}
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	ok, err := h.RelationUsecase.Check(
		domain.Node{
			Namespace: relation.SubjectNamespace,
			Name:      relation.SubjectName,
			Relation:  relation.SubjectRelation,
		},
		domain.Node{
			Namespace: relation.ObjectNamespace,
			Name:      relation.ObjectName,
			Relation:  relation.Relation,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if !ok {
		c.Status(http.StatusForbidden)
	}
	c.Status(http.StatusOK)
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
func (h *RelationHandler) GetShortestPath(c *gin.Context) {
	relation := domain.Relation{}
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.RelationUsecase.GetShortestPath(
		domain.Node{
			Namespace: relation.SubjectNamespace,
			Name:      relation.SubjectName,
			Relation:  relation.SubjectRelation,
		},
		domain.Node{
			Namespace: relation.ObjectNamespace,
			Name:      relation.ObjectName,
			Relation:  relation.Relation,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if len(paths) == 0 {
		c.Status(http.StatusForbidden)
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: paths,
	})
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
func (h *RelationHandler) GetAllPaths(c *gin.Context) {
	relation := domain.Relation{}
	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.RelationUsecase.GetAllPaths(
		domain.Node{
			Namespace: relation.SubjectNamespace,
			Name:      relation.SubjectName,
			Relation:  relation.SubjectRelation,
		},
		domain.Node{
			Namespace: relation.ObjectNamespace,
			Name:      relation.ObjectName,
			Relation:  relation.Relation,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if len(paths) == 0 {
		c.Status(http.StatusForbidden)
	}
	type response struct {
		Data [][]domain.Relation `json:"data"`
	}
	c.JSON(http.StatusOK, response{
		Data: paths,
	})
}

func (h *RelationHandler) GetAllObjectRelations(c *gin.Context) {
	type request struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
		Relation  string `json:"relation"`
	}
	subject := request{}
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	relations, err := h.RelationUsecase.GetAllObjectRelations(
		domain.Node(subject),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: relations,
	})
}

func (h *RelationHandler) GetAllSubjectRelations(c *gin.Context) {
	type request struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
		Relation  string `json:"relation"`
	}
	object := request{}
	if err := c.ShouldBindJSON(&object); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	relations, err := h.RelationUsecase.GetAllSubjectRelations(
		domain.Node(object),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: relations,
	})
}

// @Summary Clear all relations
// @Description Clear all relations in the system
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200 {string} string "All relations cleared"
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/clear-all-relations [post]
func (h *RelationHandler) ClearAllRelations(c *gin.Context) {
	err := h.RelationUsecase.ClearAllRelations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
