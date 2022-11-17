package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/handler"
)

type userHandler struct {
	userService handler.UserService
}

func NewUserHandler(usr handler.UserService) *userHandler {
	return &userHandler{userService: usr}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})

}
func (h *userHandler) New(c *gin.Context) {

	c.HTML(http.StatusOK, "user_new.html", nil)

}
func (h *userHandler) Create(c *gin.Context) {
	var input domain.FormUserInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusInternalServerError, "user_new.html", input)
		return
	}
	param := domain.UserParam{
		Name:       input.Name,
		Occupation: input.Occupation,
		Email:      input.Email,
		Password:   input.Password,
	}

	_, err = h.userService.Register(param)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}
func (h *userHandler) Edit(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	data := domain.FormUserUpdate{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
	}
	c.HTML(http.StatusOK, "user_edit.html", data)
}
func (h *userHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var input domain.FormUserUpdate
	err := c.ShouldBind(&input)

	if err != nil {
		input.Error = err
		c.HTML(http.StatusInternalServerError, "user_edit.html", input)
		return
	}
	input.ID = id

	_, err = h.userService.UpdateUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	}

	c.Redirect(http.StatusFound, "/users")

}
