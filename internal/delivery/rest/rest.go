package rest

import (
	"net/http"

	"github.com/skyrocketOoO/zanazibar-dag/domain"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	usecase domain.Usecase
}

func NewDelivery(usecase domain.Usecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}

// @Summary Check the server started
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Router /ping [get]
func (d *Delivery) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{Message: "pong"})
}

// @Summary Check the server healthy
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Failure 503 {obj} domain.Response
// @Router /healthy [get]
func (d *Delivery) Healthy(c *gin.Context) {
	// do something check
	if err := d.usecase.Healthy(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, domain.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response{Message: "healthy"})
}

// @Summary Query edges based on parameters
// @Description Query edges based on specified parameters.
// @Tags Edge
// @Accept json
// @Produce json
// @Param obj-namespace query string false "Obj Namespace"
// @Param obj-name query string false "Obj Name"
// @Param edge query string false "Edge"
// @Param sbj-namespace query string false "Sbj Namespace"
// @Param sbj-name query string false "Sbj Name"
// @Param sbj-edge query string false "Sbj Edge"
// @Param page-token query string false "Page token"
// @Param page-size query string false "Page size"
// @Success 200 {obj} delivery.Get.respBody
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [get]
func (h *Delivery) Get(c *gin.Context) {
	edge := domain.Edge{
		ObjNs:   c.Query("obj-ns"),
		ObjName: c.Query("obj-name"),
		ObjRel:  c.Query("obj-rel"),
		SbjNs:   c.Query("sbj-ns"),
		SbjName: c.Query("sbj-name"),
		SbjRel:  c.Query("sbj-rel"),
	}
	// pageToken := c.Query("page-token")
	// pageSize, err := strconv.Atoi(c.Query("page-size"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, domain.ErrResponse{
	// 		Error: err.Error(),
	// 	})
	// 	return
	// }
	edges, token, err := h.usecase.Get(edge, domain.PageOptions{
		// PageToken: pageToken,
		// PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	type respBody struct {
		Edges     []domain.Edge `json:"edges"`
		PageToken string        `json:"page_token"`
	}
	c.JSON(http.StatusOK, respBody{
		Edges:     edges,
		PageToken: token,
	})
}

// @Summary Create a new edge
// @Description Create a new edge based on the provided JSON payload.
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Create.requestBody true "Edge obj to be created"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [post]
func (h *Delivery) Create(c *gin.Context) {
	type requestBody struct {
		Edge    domain.Edge `json:"edge"`
		ExistOk bool        `json:"exist_ok"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.usecase.Create(reqBody.Edge, reqBody.ExistOk); err != nil {
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

// @Summary Delete a edge
// @Description Delete a edge based on the provided JSON payload.
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Delete.requestBody true "Edge obj to be deleted"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [delete]
func (h *Delivery) Delete(c *gin.Context) {
	type requestBody struct {
		Edge domain.Edge `json:"edge"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.usecase.Delete(reqBody.Edge); err != nil {
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

func (h *Delivery) DeleteByQueries(c *gin.Context) {
	type requestBody struct {
		Queries []domain.Edge `json:"queries"`
	}
	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	if err := h.usecase.DeleteByQueries(body.Queries); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Delivery) BatchOperation(c *gin.Context) {
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
	if err := h.usecase.BatchOperation(body.Operations); err != nil {
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
// @Description Retrieve all unique namespaces for edges.
// @Tags Edge
// @Produce json
// @Success 200 {obj} domain.StringsResponse
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-namespaces [post]
func (h *Delivery) GetAllNamespaces(c *gin.Context) {
	namespaces, err := h.usecase.GetAllNs()
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

// @Summary Check if a edge link exists
// @Description Check if a edge link exists between two entities
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Check.requestBody true "comment"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/check [post]
func (h *Delivery) Check(c *gin.Context) {
	type requestBody struct {
		Sbj        domain.Vertex     `json:"sbj" binding:"required"`
		Obj        domain.Vertex     `json:"obj" binding:"required"`
		SearchCond domain.SearchCond `json:"search_cond"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	ok, err := h.usecase.Check(body.Sbj, body.Obj, body.SearchCond)
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

// @Summary Get the shortest path between two entities in a edge graph
// @Description Get the shortest path between two entities in a edge graph
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.GetShortestPath.requestBody true "comment"
// @Success 200 {obj} domain.DataResponse "Shortest path between entities"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-shortest-path [post]
func (h *Delivery) GetShortestPath(c *gin.Context) {
	type requestBody struct {
		Sbj        domain.Vertex     `json:"sbj" binding:"required"`
		Obj        domain.Vertex     `json:"obj" binding:"required"`
		SearchCond domain.SearchCond `json:"search_cond"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.usecase.GetShortestPath(body.Sbj, body.Obj, body.SearchCond)
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
	c.JSON(http.StatusOK, domain.EdgesResponse{
		Edges: paths,
	})
}

// @Summary Get all paths between two entities in a edge graph
// @Description Get all paths between two entities in a edge graph
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.GetAllPaths.requestBody true "Edge obj specifying the entities"
// @Success 200 {obj} delivery.GetAllPaths.response "All paths between entities"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-paths [post]
func (h *Delivery) GetAllPaths(c *gin.Context) {
	type requestBody struct {
		Sbj        domain.Vertex     `json:"sbj" binding:"required"`
		Obj        domain.Vertex     `json:"obj" binding:"required"`
		SearchCond domain.SearchCond `json:"search_cond"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	paths, err := h.usecase.GetAllPaths(body.Sbj, body.Obj, body.SearchCond)
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
		Data [][]domain.Edge `json:"data"`
	}
	c.JSON(http.StatusOK, response{
		Data: paths,
	})
}

// @Summary Get all edges for a given obj
// @Description Get all edges for a given obj specified by namespace, name, and edge
// @Tags Edge
// @Accept json
// @Produce json
// @Param sbj body delivery.GetAllObjEdges.requestBody true "Obj information (namespace, name, edge)"
// @Success 200 {obj} domain.DataResponse "All edges for the specified obj"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-obj-edges [post]
func (h *Delivery) GetAllObjEdges(c *gin.Context) {
	type requestBody struct {
		Sbj         domain.Vertex      `json:"sbj" binding:"required"`
		SearchCond  domain.SearchCond  `json:"search_cond"`
		CollectCond domain.CollectCond `json:"collect_cond"`
		MaxDepth    int                `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	edges, err := h.usecase.GetAllObjRels(
		domain.Vertex(body.Sbj),
		body.SearchCond,
		body.CollectCond,
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
	c.JSON(http.StatusOK, domain.EdgesResponse{
		Edges: edges,
	})
}

// @Summary Get all edges for a given sbj
// @Description Get all edges for a given sbj specified by namespace, name, and edge
// @Tags Edge
// @Accept json
// @Produce json
// @Param obj body delivery.GetAllSbjEdges.requestBody true "Sbj information (namespace, name, edge)"
// @Success 200 {obj} domain.DataResponse "All edges for the specified sbj"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-sbj-edges [post]
func (h *Delivery) GetAllSbjEdges(c *gin.Context) {
	type requestBody struct {
		Obj         domain.Vertex      `json:"obj" binding:"required"`
		SearchCond  domain.SearchCond  `json:"search_cond"`
		CollectCond domain.CollectCond `json:"collect_cond"`
		MaxDepth    int                `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	edges, err := h.usecase.GetAllSbjRels(
		domain.Vertex(body.Obj),
		body.SearchCond,
		body.CollectCond,
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
	c.JSON(http.StatusOK, domain.EdgesResponse{
		Edges: edges,
	})
}

func (h *Delivery) GetTree(c *gin.Context) {
	type requestBody struct {
		Sbj      domain.Vertex `json:"sbj" binding:"required"`
		MaxDepth int           `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	tree, err := h.usecase.GetTree(
		body.Sbj,
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
	type response struct {
		Tree domain.TreeNode `json:"tree"`
	}
	c.JSON(http.StatusOK, response{
		Tree: *tree,
	})
}

// @Summary Clear all edges
// @Description Clear all edges in the system
// @Tags Edge
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/clear-all-edges [post]
func (h *Delivery) ClearAllEdges(c *gin.Context) {
	err := h.usecase.ClearAllEdges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrResponse{
			Error: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
