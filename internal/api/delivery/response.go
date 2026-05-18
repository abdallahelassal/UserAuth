package delivery

import "github.com/gin-gonic/gin"

type Response struct {
	Code 	int 		`json:"code"`
	Message string 		`json:"message"`
	Data 	interface{} `json:"data"`
}

func Success(data interface{})Response{
	return Response{
		Code: 200,
		Message: "success",
		Data: data,
	}
}

func Error(code int,msg string)Response{
	return Response{
		Code: code,
		Message: msg,
		Data: nil,
	}
}

func Json(c *gin.Context, resp Response){
	c.JSON(resp.Code,resp)
}
