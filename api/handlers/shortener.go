package handlers

import (
	"fmt"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/pkg/utils"
	"net/http"

	http_status "go_auth_api_gateway/api/http"

	"github.com/gin-gonic/gin"
)

// CreateShortUrl godoc
// @ID create_short_url
// @Router /v1/short-url [POST]
// @Summary Create ShortUrl
// @Description Create ShortUrl
// @Tags urls
// @Accept json
// @Produce json
// @Param body body auth_service.CreateShortUrlRequest true "Request body"
// @Success 201 {object} http.Response{data=auth_service.CreateShortUrlResponse} "Response Body"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) CreateShortUrl(c *gin.Context) {
	var req pb.CreateShortUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		h.handleResponse(c, http_status.BadRequest, err.Error())
		return
	}

	if !utils.IsLongCorrect(string(req.GetLongUrl())) {
		err := fmt.Errorf(utils.InvalidURLError, req.GetLongUrl())
		h.handleResponse(c, http_status.BadRequest, err.Error())
		return
	}

	resp, err := h.services.ShortenerService().CreateShortUrl(c, &req)
	if err != nil {
		h.handleResponse(c, http_status.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http_status.OK, resp)
}

// GetShortUrl godoc
// @ID get_short_url
// @Router /v1/short-url/{hash} [GET]
// @Summary Get ShortUrl
// @Description Get ShortUrl
// @Tags urls
// @Accept json
// @Produce json
// @Param body body auth_service.GetShortUrlRequest true "Request body"
// @Success 201 {object} http.Response{data=auth_service.GetShortUrlResponse} "Response Body"
func (h *Handler) GetShortUrl(c *gin.Context) {

	var req pb.GetShortUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http_status.BadRequest, err.Error())
		return
	}

	hash := c.Param("hash")
	if !utils.IsShortCorrect(hash) {
		err := fmt.Errorf(utils.InvalidHashError, hash)
		h.handleResponse(c, http_status.BadRequest, err.Error())
		return
	}

	req.ShortUrl = hash
	resp, err := h.services.ShortenerService().GetShortUrl(c, &req)
	if err != nil {
		h.handleResponse(c, http_status.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http_status.OK, resp)
}

// HandleLonger godoc
// @ID handle_longer
// @Router /sigma/{hash} [GET]
// @Summary Handle Longer
// @Description Handle Longer
// @Tags urls
// @Param hash path string true "short url hash"
// @Success 201 {object} http.Response{data=string} "Response Body"
func (h *Handler) HandleLonger(c *gin.Context) {

	url := c.Param("hash")
	if !utils.IsShortCorrect(url) {
		err := fmt.Errorf(utils.InvalidHashError, url)
		h.handleResponse(c, http_status.BadRequest, err.Error())
		return
	}

	resp, err := h.services.ShortenerService().GetShortUrl(
		c.Request.Context(),
		&pb.GetShortUrlRequest{
			ShortUrl: url,
		},
	)
	if err != nil {
		h.handleResponse(c, http_status.InternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, resp.GetLongUrl())
}
