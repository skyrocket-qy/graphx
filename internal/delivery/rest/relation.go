package rest

import (
	"net/http"
	"strconv"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
	usecasedomain "github.com/skyrocketOoO/zanazibar-dag/domain/usecase"

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

// @Summary Query relations based on parameters
// @Description Query relations based on specified parameters.
// @Tags Relation
// @Accept json
// @Produce json
// @Param object-namespace query string false "Object Namespace"
// @Param object-name query string false "Object Name"
// @Param relation query string false "Relation"
// @Param subject-namespace query string false "Subject Namespace"
// @Param subject-name query string false "Subject Name"
// @Param subject-relation query string false "Subject Relation"
// @Param page-token query string false "Page token"
// @Param page-size query string false "Page size"
// @Success 200 {object} delivery.Get.respBody
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/ [get]
func (h *RelationHandler) Get(c *gin.Context) {
	relation := domain.Relation{
		ObjectNamespace:  c.Query("object-namespace"),
		ObjectName:       c.Query("object-name"),
		Relation:         c.Query("relation"),
		SubjectNamespace: c.Query("subject-namespace"),
		SubjectName:      c.Query("subject-name"),
		SubjectRelation:  c.Query("subject-relation"),
	}
	pageToken := c.Query("page-token")
	pageSize, err := strconv.Atoi(c.Query("page-size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	relations, token, err := h.RelationUsecase.Get(relation, usecasedomain.PageOptions{
		PageToken: pageToken,
		PageSize:  pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	type respBody struct {
		Relations []domain.Relation `json:"relations"`
		PageToken string            `json:"page_token"`
	}
	c.JSON(http.StatusOK, respBody{
		Relations: relations,
		PageToken: token,
	})
}

// @Summary Create a new relation
// @Description Create a new relation based on the provided JSON payload.
// @Tags Relation
// @Accept json
// @Produce json
// @Param relation body delivery.Create.requestBody true "Relation object to be created"
// @Success 200
// @Failure 400 {object} domain.ErrResponse
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/ [post]
func (h *RelationHandler) Create(c *gin.Context) {
	type requestBody struct {
		Relation domain.Relation `json:"relation"`
		ExistOk  bool            `json:"exist_ok"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.Create(reqBody.Relation, reqBody.ExistOk); err != nil {
		if _, ok := err.(domain.CauseCycleError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		} else if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Delete a relation
// @Description Delete a relation based on the provided JSON payload.
// @Tags Relation
// @Accept json
// @Produce json
// @Param relation body delivery.Delete.requestBody true "Relation object to be deleted"
// @Success 200
// @Failure 400 {object} domain.ErrResponse
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/ [delete]
func (h *RelationHandler) Delete(c *gin.Context) {
	type requestBody struct {
		Relation domain.Relation `json:"relation"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.Delete(reqBody.Relation); err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RelationHandler) DeleteByQueries(c *gin.Context) {
	type requestBody struct {
		Queries []domain.Relation `json:"queries"`
	}
	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.DeleteByQueries(body.Queries); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *RelationHandler) BatchOperation(c *gin.Context) {
	type requestBody struct {
		Operations []domain.Operation `json:"operations"`
	}
	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.RelationUsecase.BatchOperation(body.Operations); err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get all unique namespaces
// @Description Retrieve all unique namespaces for relations.
// @Tags Relation
// @Produce json
// @Success 200 {object} domain.StringsResponse
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-all-namespaces [post]
func (h *RelationHandler) GetAllNamespaces(c *gin.Context) {
	namespaces, err := h.RelationUsecase.GetAllNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.StringsResponse{
		Data: namespaces,
	})
}

// @Summary Check if a relation link exists
// @Description Check if a relation link exists between two entities
// @Tags Relation
// @Accept json
// @Produce json
// @Param relation body delivery.Check.requestBody true "comment"
// @Success 200
// @Failure 400 {object} domain.ErrResponse
// @Failure 403
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/check [post]
func (h *RelationHandler) Check(c *gin.Context) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject" binding:"required"`
		Object          domain.Node            `json:"object" binding:"required"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	ok, err := h.RelationUsecase.Check(body.Subject, body.Object, body.SearchCondition)
	if err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
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
// @Param relation body delivery.GetShortestPath.requestBody true "comment"
// @Success 200 {object} domain.DataResponse "Shortest path between entities"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-shortest-path [post]
func (h *RelationHandler) GetShortestPath(c *gin.Context) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject" binding:"required"`
		Object          domain.Node            `json:"object" binding:"required"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.RelationUsecase.GetShortestPath(body.Subject, body.Object, body.SearchCondition)
	if err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if len(paths) == 0 {
		c.Status(http.StatusForbidden)
	}
	c.JSON(http.StatusOK, domain.RelationsResponse{
		Relations: paths,
	})
}

// @Summary Get all paths between two entities in a relation graph
// @Description Get all paths between two entities in a relation graph
// @Tags Relation
// @Accept json
// @Produce json
// @Param relation body delivery.GetAllPaths.requestBody true "Relation object specifying the entities"
// @Success 200 {object} delivery.GetAllPaths.response "All paths between entities"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-all-paths [post]
func (h *RelationHandler) GetAllPaths(c *gin.Context) {
	type requestBody struct {
		Subject         domain.Node            `json:"subject" binding:"required"`
		Object          domain.Node            `json:"object" binding:"required"`
		SearchCondition domain.SearchCondition `json:"search_condition"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.RelationUsecase.GetAllPaths(body.Subject, body.Object, body.SearchCondition)
	if err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
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

// @Summary Get all relations for a given object
// @Description Get all relations for a given object specified by namespace, name, and relation
// @Tags Relation
// @Accept json
// @Produce json
// @Param subject body delivery.GetAllObjectRelations.requestBody true "Object information (namespace, name, relation)"
// @Success 200 {object} domain.DataResponse "All relations for the specified object"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-all-object-relations [post]
func (h *RelationHandler) GetAllObjectRelations(c *gin.Context) {
	type requestBody struct {
		Subject          domain.Node             `json:"subject" binding:"required"`
		SearchCondition  domain.SearchCondition  `json:"search_condition"`
		CollectCondition domain.CollectCondition `json:"collect_condition"`
		MaxDepth         int                     `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	relations, err := h.RelationUsecase.GetAllObjectRelations(
		domain.Node(body.Subject),
		body.SearchCondition,
		body.CollectCondition,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.RelationsResponse{
		Relations: relations,
	})
}

// @Summary Get all relations for a given subject
// @Description Get all relations for a given subject specified by namespace, name, and relation
// @Tags Relation
// @Accept json
// @Produce json
// @Param object body delivery.GetAllSubjectRelations.requestBody true "Subject information (namespace, name, relation)"
// @Success 200 {object} domain.DataResponse "All relations for the specified subject"
// @Failure 400 {object} domain.ErrResponse
// @Failure 403
// @Failure 500 {object} domain.ErrResponse
// @Router /relation/get-all-subject-relations [post]
func (h *RelationHandler) GetAllSubjectRelations(c *gin.Context) {
	type requestBody struct {
		Object           domain.Node             `json:"object" binding:"required"`
		SearchCondition  domain.SearchCondition  `json:"search_condition"`
		CollectCondition domain.CollectCondition `json:"collect_condition"`
		MaxDepth         int                     `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	relations, err := h.RelationUsecase.GetAllSubjectRelations(
		domain.Node(body.Object),
		body.SearchCondition,
		body.CollectCondition,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.RequestBodyError); ok {
			c.JSON(http.StatusBadRequest, domain.ErrResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.RelationsResponse{
		Relations: relations,
	})
}

// @Summary Clear all relations
// @Description Clear all relations in the system
// @Tags Relation
// @Accept json
// @Produce json
// @Success 200
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
