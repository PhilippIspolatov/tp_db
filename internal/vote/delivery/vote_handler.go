package delivery

import (
	"net/http"
	"strconv"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/thread"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/vote"
	"github.com/labstack/echo"
)

type VoteHandler struct {
	voteUcase vote.Usecase
	threadUcase thread.Usecase
}

func NewVoteHandler(router *echo.Echo, vUC vote.Usecase, tUC thread.Usecase) *VoteHandler {
	vh := &VoteHandler{
		voteUcase: vUC,
		threadUcase :tUC,
	}

	router.POST("/thread/:slug_or_id/vote", vh.CreateVote())

	return vh
}

func (vh *VoteHandler) CreateVote() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug_or_id")
		v := &models.Vote{}
		if err := c.Bind(v); err != nil {
			return c.JSON(http.StatusBadRequest, tools.Message{
				Message:"bad request",
			})
		}

		id, err := strconv.ParseUint(slug, 10, 64)
		if err != nil {
			t, err := vh.threadUcase.SelectThreadBySlug(slug)
			if err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			}
			v.Thread = t.Id
			if err := vh.voteUcase.SelectVote(v); err == nil {
				if err = vh.voteUcase.UpdateVote(v); err != nil {
					return c.JSON(http.StatusBadRequest, tools.Message{
						Message:"bad request",
					})
				}
				t, _ := vh.threadUcase.SelectThreadById(v.Thread)
				return c.JSON(http.StatusOK, t)
			}
			if err := vh.voteUcase.InsertVote(v); err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			}
			t, err = vh.threadUcase.SelectThreadBySlug(slug)
			return c.JSON(http.StatusOK, t)
		}
		t, err := vh.threadUcase.SelectThreadById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		v.Thread = t.Id
		if err := vh.voteUcase.SelectVote(v); err == nil {
			if err = vh.voteUcase.UpdateVote(v); err != nil {
				return c.JSON(http.StatusBadRequest, tools.Message{
					Message:"bad request",
				})
			}
			t, _ := vh.threadUcase.SelectThreadById(v.Thread)
			return c.JSON(http.StatusOK, t)
		}
		if err := vh.voteUcase.InsertVote(v); err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		t, err = vh.threadUcase.SelectThreadById(id)
		return c.JSON(http.StatusOK, t)

	}
}
