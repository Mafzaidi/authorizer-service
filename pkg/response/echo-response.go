package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Response struct {
		Status  *StatusInfo `json:"status"`
		Message string      `json:"message,omitempty"`
		Meta    interface{} `json:"meta,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}

	StatusInfo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Pagination struct {
		CurrentPage  int `json:"current_page"`
		PageLimit    int `json:"page_limit"`
		TotalRecords int `json:"total_records"`
		TotalPages   int `json:"total_pages"`
		PrevPages    int `json:"prev_pages"`
		NextPages    int `json:"next_pages"`
	}
)

func SuccesHandler(c echo.Context, r *Response) error {
	return c.JSON(http.StatusOK, &Response{
		Status: &StatusInfo{
			Code:    int(http.StatusOK),
			Message: r.Message,
		},
		Data: r.Data,
		Meta: r.Meta,
	})
}

func ErrorHandler(c echo.Context, statusCode int, statusStr, message string) error {
	return c.JSON(int(statusCode), &Response{
		Status: &StatusInfo{
			Code:    int(statusCode),
			Message: statusStr,
		},
		Message: "error occurred",
		Error:   message,
	})
}
