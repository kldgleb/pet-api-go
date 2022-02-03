package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test-api/pkg/entity"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input entity.TodoList
	if err = c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listId, err := h.services.Create(input, userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "error while creating todoList")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"listId": listId,
	})
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, _ := c.Get("userId")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
}

func (h *Handler) getListById(c *gin.Context) {

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
