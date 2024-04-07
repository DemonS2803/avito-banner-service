package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func Send200Success(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	//render.JSON(w, r, Response{Error: "very bad"})
	//w.Write([]byte("don't miss required params)"))
}

func Send400Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, Response{Error: "very bad"})
	//w.Write([]byte("don't miss required params)"))
}

func Send401Error(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusUnauthorized)
	render.Status(r, http.StatusUnauthorized)
	render.JSON(w, r, Response{Error: "invalid token("})
	//render.JSON(w, r, Response{Error: "invalid token("})
}

func Send404Error(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusNotFound)
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, Response{Error: "resource not found"})
	//render.JSON(w, r, Response{Error: "resource not found"})
	//w.Write([]byte("resource not found"))
}

func Send403Error(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusForbidden)
	render.Status(r, http.StatusForbidden)
	render.JSON(w, r, Response{Error: "no access"})
	//render.JSON(w, r, Response{Error: "no access"})
	//w.Write([]byte("no access"))
}

func Send500Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, Response{Error: "technical chocolates..."})

}
