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

	router.POST("/api/thread/:slug_or_id/vote", vh.CreateVote())

	return vh
}

func (vh *VoteHandler) CreateVote() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug_or_id")
		v := &models.Vote{}
		oldV := &models.Vote{}
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
			oldV.Thread = t.Id
			oldV.Nickname = v.Nickname
			if err := vh.voteUcase.SelectVote(oldV); err == nil {
				t, _ := vh.threadUcase.SelectThreadById(v.Thread)
				t.Votes = t.Votes - oldV.Voice + v.Voice
				if err = vh.voteUcase.UpdateVote(v); err != nil {
					return c.JSON(http.StatusBadRequest, tools.Message{
						Message:"bad request1",
					})
				}

				return c.JSON(http.StatusOK, t)
			}
			t, err = vh.threadUcase.SelectThreadBySlug(slug)
			t.Votes += v.Voice
			if err := vh.voteUcase.InsertVote(v); err != nil {
				return c.JSON(http.StatusNotFound, tools.Message{
					Message:"not found",
				})
			}
			return c.JSON(http.StatusOK, t)
		}
		t, err := vh.threadUcase.SelectThreadById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		v.Thread = t.Id
		oldV.Thread = t.Id
		oldV.Nickname = v.Nickname
		if err := vh.voteUcase.SelectVote(oldV); err == nil {
			t, _ := vh.threadUcase.SelectThreadById(v.Thread)
			t.Votes = t.Votes - oldV.Voice + v.Voice
			if err = vh.voteUcase.UpdateVote(v); err != nil {
				return c.JSON(http.StatusBadRequest, tools.Message{
					Message:"bad request2",
				})
			}

			return c.JSON(http.StatusOK, t)
		}
		t, err = vh.threadUcase.SelectThreadById(id)
		t.Votes += v.Voice
		if err := vh.voteUcase.InsertVote(v); err != nil {
			return c.JSON(http.StatusNotFound, tools.Message{
				Message:"not found",
			})
		}
		return c.JSON(http.StatusOK, t)

	}
}
