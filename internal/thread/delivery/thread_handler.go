package delivery

import (
	"net/http"
	"strconv"

	"github.com/PhilippIspolatov/tp_db/internal/forum"
	"github.com/PhilippIspolatov/tp_db/internal/models"
	"github.com/PhilippIspolatov/tp_db/internal/thread"
	"github.com/PhilippIspolatov/tp_db/internal/tools"
	"github.com/PhilippIspolatov/tp_db/internal/user"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ThreadHandler struct {
	threadUcase thread.Usecase
	userUcase   user.Usecase
	forumUcase  forum.Usecase
}

func NewThreadHandler(router *echo.Echo, tUC thread.Usecase, uUC user.Usecase, fUC forum.Usecase) *ThreadHandler {
	th := &ThreadHandler{
		threadUcase: tUC,
		userUcase:   uUC,
		forumUcase:  fUC,
	}

	router.POST("/api/forum/:slug/create", th.CreateThread())
	router.GET("/api/forum/:slug/threads", th.GetThreads())
	router.GET("/api/thread/:slug_or_id/details", th.GetThread())
	router.POST("/api/thread/:slug_or_id/details", th.UpdateThread())

	return th
}

func (th *ThreadHandler) UpdateThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		t := &models.Thread{}
		if err := c.Bind(t); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}
		slug := c.Param("slug_or_id")
		id, err := strconv.ParseUint(slug, 10, 64)
		if err != nil {
			t.Slug = slug
			if err := th.threadUcase.UpdateThreadBySlug(t); err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			}
			return c.JSON(http.StatusOK, t)
		}
		t.Id = id
		if err := th.threadUcase.UpdateThreadById(t); err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		return c.JSON(http.StatusOK, t)
	}
}

func (th *ThreadHandler) CreateThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		f := c.Param("slug")

		gf, err := th.forumUcase.GetForum(f)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}

		t := &models.Thread{}

		if err := c.Bind(t); err != nil {
			return err
		}

		t.Forum = gf.Slug

		if _, err := th.userUcase.GetUser(t.Author); err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}

		if conflict, err := th.threadUcase.SelectThreadBySlug(t.Slug); err == nil && conflict.Slug != "" {
			return c.JSON(http.StatusConflict, conflict)
		}

		if err := th.threadUcase.InsertThread(t); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message {

				Message: "bad request",
			})
		}

		return c.JSON(http.StatusCreated, t)
	}
}

func (th *ThreadHandler) GetThreads() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")
		limit, _ := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
		since := c.QueryParam("since")
		desc, _ := strconv.ParseBool(c.QueryParam("desc"))

		if _, err := th.forumUcase.GetForum(slug); err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}

		t, err := th.threadUcase.SelectThreads(slug, limit, since, desc)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message: "bad request",
			})
		}
		return c.JSON(http.StatusOK, t)
	}
}

func (th *ThreadHandler) GetThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug_or_id")
		id, err := strconv.ParseUint(slug, 10, 64)
		if err != nil {
			t, err := th.threadUcase.SelectThreadBySlug(slug)
			if err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message: "not found",
				})
			}
			return c.JSON(http.StatusOK, t)
		}
		t, err := th.threadUcase.SelectThreadById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		return c.JSON(http.StatusOK, t)
	}
}
