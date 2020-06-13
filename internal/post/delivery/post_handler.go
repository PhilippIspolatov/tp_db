package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ifo16u375/tp_db/internal/forum"
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/post"
	"github.com/ifo16u375/tp_db/internal/thread"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/user"
	"github.com/labstack/echo"
)

type PostHandler struct {
	postUcase post.Usecase
	threadUcase thread.Usecase
	userUcase user.Usecase
	forumUcase forum.Usecase
}

func NewPostHandler(router *echo.Echo, pUC post.Usecase, tUC thread.Usecase, uUC user.Usecase, fUC forum.Usecase) *PostHandler {
	ph := &PostHandler{
		postUcase: pUC,
		threadUcase: tUC,
		userUcase: uUC,
		forumUcase: fUC,
	}

	router.POST("/api/thread/:slug_or_id/create", ph.CreatePosts())
	router.GET("/api/post/:id/details", ph.GetPost())
	router.POST("/api/post/:id/details", ph.UpdatePost())
	router.GET("/api/thread/:slug_or_id/posts", ph.GetSortedPosts())

	return ph
}

func (ph *PostHandler) UpdatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}
		p := &models.Post{
			Id:id,
		}
		if err := c.Bind(p); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}
		if err := ph.postUcase.UpdatePost(p); err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		return c.JSON(http.StatusOK, p)
	}
}

func (ph *PostHandler) CreatePosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		posts := []*models.Post{}
		if err := c.Bind(&posts); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}

		slug := c.Param("slug_or_id")
		id, err := strconv.ParseUint(slug, 10, 64)
		if err != nil {
			t, err := ph.threadUcase.SelectThreadBySlug(slug)
			if err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			}
			if len(posts) == 0 {
				return c.JSON(http.StatusCreated, posts)
			}
			if err := ph.postUcase.CheckPosts(posts, t.Id); err != nil {
				if err == tools.ErrNotFound {
					return c.JSON(http.StatusNotFound, tools.Message{
						Message:"not found",
					})
				} else if err == tools.ErrConflict {
					return c.JSON(http.StatusConflict, tools.Message{
						Message:"conflict",
					})
				}
			}

			if err := ph.postUcase.InsertPosts(posts, t.Id, t.Forum); err != nil {
				return c.JSON(http.StatusConflict, tools.Message{
					Message:"conflict",
				})
			}
			return c.JSON(http.StatusCreated, posts)
		}
		t, err := ph.threadUcase.SelectThreadById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		if len(posts) == 0 {
			return c.JSON(http.StatusCreated, posts)
		}
		if err := ph.postUcase.CheckPosts(posts, t.Id); err != nil {
			if err == tools.ErrNotFound {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			} else if err == tools.ErrConflict {
				return c.JSON(http.StatusConflict, tools.Message{
					Message:"conflict",
				})
			}
		}
		if err := ph.postUcase.InsertPosts(posts, t.Id, t.Forum); err != nil {
			return c.JSON(http.StatusConflict, tools.Message{
				Message:"conflict",
			})
		}
		return c.JSON(http.StatusCreated, posts)
	}
}

func (ph *PostHandler) GetPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		type PostInfo struct {
			Author *models.User   `json:"author"`
			Forum  *models.Forum  `json:"forum"`
			Post   *models.Post   `json:"post"`
			Thread *models.Thread `json:"thread"`
		}
		res := &PostInfo{}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}
		related := strings.Split(c.QueryParam("related"), ",")

		p, err := ph.postUcase.SelectPost(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		res.Post = p

		for _, val := range related {
			switch val {
			case "forum": {
				f, _ := ph.forumUcase.GetForum(p.Forum)
				res.Forum = f
			}
			case "thread": {
				t, _ := ph.threadUcase.SelectThreadById(p.Thread)
				res.Thread = t
			}
			case "user": {
				u, _ := ph.userUcase.GetUser(p.Author)
				res.Author = u
			}
			}
		}
		return c.JSON(http.StatusOK, res)
	}
}

func (ph *PostHandler) GetSortedPosts() echo.HandlerFunc {
	type Request struct {
		Sort string `json:"sort"`
		Desc bool `json:"desc"`
		Since uint64 `json:"since"`
		Limit uint64 `json:"limit"`
	}
	return func(c echo.Context) error {
		req := &Request{}

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message: "bad request",
			})
		}

		slug := c.Param("slug_or_id")

		id, err := strconv.ParseUint(slug, 10, 64)

		if err != nil {
			t, err := ph.threadUcase.SelectThreadBySlug(slug)
			if err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message: "not found",
				})
			}

			posts, err := ph.postUcase.SelectSortesPosts(t.Id, req.Sort, req.Desc, req.Since, req.Limit)

			if err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message: "not found",
				})
			}

			return c.JSON(http.StatusOK, posts)
		}

		_, err = ph.threadUcase.SelectThreadById(id)

		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}

		posts, err := ph.postUcase.SelectSortesPosts(id, req.Sort, req.Desc, req.Since, req.Limit)

		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message: "not found",
			})
		}

		return c.JSON(http.StatusOK, posts)
	}
}
