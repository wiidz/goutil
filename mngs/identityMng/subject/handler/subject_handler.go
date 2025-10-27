package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wiidz/goutil/mngs/identityMng/subject/dto"
	"github.com/wiidz/goutil/mngs/identityMng/subject/service"
)

type Handler struct{ S *service.Service }

func New(s *service.Service) *Handler { return &Handler{S: s} }

func (h *Handler) List(c *gin.Context) {
	var req dto.ListSubjectsRequest
	_ = c.ShouldBindQuery(&req)
	items, total, err := h.S.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "LIST_FAILED", "message": err.Error()})
		return
	}
	resp := make([]dto.SubjectResponse, 0, len(items))
	for _, m := range items {
		resp = append(resp, dto.SubjectResponse{ID: m.ID, SubjectType: m.SubjectType, LoginID: m.LoginID, ExternalID: m.ExternalID, Status: m.Status})
	}
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	m, err := h.S.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SubjectResponse{ID: m.ID, SubjectType: m.SubjectType, LoginID: m.LoginID, ExternalID: m.ExternalID, Status: m.Status})
}
