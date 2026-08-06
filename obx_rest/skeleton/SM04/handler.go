package SM04

import (
	"net/http"
	"obx_rest/mechanic"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) ListSignatureType(c *gin.Context) {
	list, err := h.usecase.ListSignatureType()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature types retrieved",
		"data":    list,
	})
}

func (h *Handler) CreateSignatureType(c *gin.Context) {
	var req TypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateSignatureType(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signature type created",
	})
}

func (h *Handler) UpdateSignatureType(c *gin.Context) {
	id := c.Param("id")
	var req TypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateSignatureType(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature type updated",
	})
}

func (h *Handler) ListUser(c *gin.Context) {
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	search := c.Query("search")
	list, err := h.usecase.ListUser(companyID, search)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved",
		"data":    list,
	})
}

func (h *Handler) ListRequest(c *gin.Context) {
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListRequest(companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Requests retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateRequest(c *gin.Context) {
	var req RequestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateRequest(c.GetString("companyId"), req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Request created successfully",
	})
}

func (h *Handler) UpdateRequest(c *gin.Context) {
	id := c.Param("id")
	var req RequestUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateRequest(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Request updated successfully",
	})
}

func (h *Handler) DeleteRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteRequest(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Request deleted successfully",
	})
}

func (h *Handler) ListDocument(c *gin.Context) {
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListDocument(companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Documents retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateDocument(c *gin.Context) {
	var req DocumentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	userID, _ := c.Get("userId")
	uid, _ := userID.(string)
	if err := h.usecase.CreateDocument(c.GetString("companyId"), uid, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Document created successfully",
	})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	var req DocumentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateDocument(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document updated successfully",
	})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteDocument(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}

func (h *Handler) ListDocumentCategory(c *gin.Context) {
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	list, err := h.usecase.ListDocumentCategory(companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document categories retrieved",
		"data":    list,
	})
}

func (h *Handler) ListDocumentRevision(c *gin.Context) {
	documentID := c.Param("id")
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListDocumentRevision(documentID, companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document revisions retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateDocumentRevision(c *gin.Context) {
	documentID := c.Param("id")
	var req DocumentRevisionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	userID, _ := c.Get("userId")
	uid, _ := userID.(string)
	if err := h.usecase.CreateDocumentRevision(documentID, c.GetString("companyId"), uid, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Document revision created successfully",
	})
}

func (h *Handler) UpdateDocumentRevision(c *gin.Context) {
	revisionID := c.Param("revisionId")
	var req DocumentRevisionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateDocumentRevision(revisionID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document revision updated successfully",
	})
}

func (h *Handler) DeleteDocumentRevision(c *gin.Context) {
	revisionID := c.Param("revisionId")
	if err := h.usecase.DeleteDocumentRevision(revisionID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document revision deleted successfully",
	})
}

func (h *Handler) ListDocumentEvidence(c *gin.Context) {
	documentID := c.Param("id")
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListDocumentEvidence(documentID, companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document evidence retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateDocumentEvidence(c *gin.Context) {
	documentID := c.Param("id")
	var req DocumentEvidenceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	userID, _ := c.Get("userId")
	uid, _ := userID.(string)
	if err := h.usecase.CreateDocumentEvidence(documentID, c.GetString("companyId"), uid, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Document evidence created successfully",
	})
}

func (h *Handler) UpdateDocumentEvidence(c *gin.Context) {
	evidenceID := c.Param("evidenceId")
	var req DocumentEvidenceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateDocumentEvidence(evidenceID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document evidence updated successfully",
	})
}

func (h *Handler) DeleteDocumentEvidence(c *gin.Context) {
	evidenceID := c.Param("evidenceId")
	if err := h.usecase.DeleteDocumentEvidence(evidenceID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Document evidence deleted successfully",
	})
}

func (h *Handler) ListForm(c *gin.Context) {
	requestID := c.Param("id")
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	list, err := h.usecase.ListForm(requestID, companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature forms retrieved",
		"data":    list,
	})
}

func (h *Handler) GenerateForms(c *gin.Context) {
	requestID := c.Param("id")
	count, err := h.usecase.GenerateForms(requestID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signature forms generated",
		"count":   count,
	})
}

func (h *Handler) FlagAction(c *gin.Context) {
	flagID := c.Param("flagId")
	var req FlagActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	userID, _ := c.Get("userId")
	uid, _ := userID.(string)
	if err := h.usecase.FlagAction(flagID, uid, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature flag updated",
	})
}
