package delivery

import (
	"net/http"

	"github.com/ifo16u375/tp_db/internal/forum"
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/user"
	"github.com/labstack/echo"
)

type ForumHandler struct {
	forumUcase forum.Usecase
	userUcase  user.Usecase
}

func NewForumHandler(router *echo.Echo, fUC forum.Usecase, uUC user.Usecase) *ForumHandler {
	fh := &ForumHandler{
		forumUcase: fUC,
		userUcase:  uUC,
	}

	router.GET("/api/forum/:slug/details", fh.GetForum())
	router.POST("/api/forum/create", fh.CreateForum())

	return fh
}

func (fh *ForumHandler) GetForum() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")
		f, err := fh.forumUcase.GetForum(slug)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{Message: "not found"})
		}
		return c.JSON(http.StatusOK, f)
	}
}

func (fh *ForumHandler) CreateForum() echo.HandlerFunc {
	return func(c echo.Context) error {
		f := &models.Forum{}
		if err := c.Bind(f); err != nil {
			return err
		}

		if conflict, err := fh.forumUcase.GetForum(f.Slug); err == nil {
			return c.JSON(http.StatusConflict, conflict)
		}
		u, err := fh.userUcase.GetUser(f.User)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}
		f.User = u.Nickname
		if err := fh.forumUcase.CreateForum(f); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message: "bad request",
			})
		}
		return c.JSON(http.StatusCreated, f)
	}
}
