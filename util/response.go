package util

import (
	"encoding/json"
	"math"
	"net/http"

	"log"
	"reflect"

	"gorm.io/gorm"
)

type JsonData struct {
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Code       int         `json:"code,omitempty"`
	PageNumber int         `json:"page_number,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	PageTotal  int         `json:"total_page,omitempty"`
	Count      int         `json:"count"`
}

func JsonErrorResponse(message string) *JsonData {
	return &JsonData{
		Error:  message,
		Status: "fail",
	}
}

func JsonResponse(data interface{}) *JsonData {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		slice := reflect.ValueOf(data)
		sliceLen := slice.Len()
		// log.Println(slice)
		array := make([]interface{}, 0)
		// log.Println(array)
		for i := 0; i < sliceLen; i++ {
			array = append(array, slice.Index(i).Interface())
		}
		return &JsonData{
			Status: "success",
			Data:   array,
		}
	} else {
		return &JsonData{
			Status: "success",
			Data:   data,
		}
	}
}

func (res *JsonData) WithMessage(message string) *JsonData {
	res.Message = message
	return res
}

func (res *JsonData) WithPaging(pageNumber int, pageSize int, pageTotal int, count int) *JsonData {
	res.PageNumber = pageNumber
	res.PageSize = pageSize
	res.PageTotal = pageTotal
	res.Count = count
	return res
}

func (res *JsonData) WithPagingEmpty(rModel *RequestModel) *JsonData {
	res.PageNumber = rModel.PageNumber()
	res.PageSize = rModel.PageSize()
	res.PageTotal = 0
	res.Count = 0
	return res
}

func (res *JsonData) WithPagingScope(scope *gorm.DB, rModel *RequestModel) *JsonData {
	var count int64
	delete(scope.Statement.Clauses, "ORDER BY")
	scope.Offset(-1).Limit(-1).Count(&count)
	res.Count = int(count)
	res.PageNumber = rModel.PageNumber()
	res.PageSize = rModel.PageSize()
	res.PageTotal = int(math.Ceil(float64(res.Count) / float64(rModel.PageSize())))
	return res
}

func (res *JsonData) WithPagingScopeDistinct(scope *gorm.DB, rModel *RequestModel, distinct string) *JsonData {
	var count int64
	delete(scope.Statement.Clauses, "ORDER BY")
	scope.Offset(-1).Limit(-1).Distinct(distinct).Count(&count)
	res.Count = int(count)
	res.PageNumber = rModel.PageNumber()
	res.PageSize = rModel.PageSize()
	res.PageTotal = int(math.Ceil(float64(res.Count) / float64(rModel.PageSize())))
	return res
}

func (res *JsonData) WithErrorCode(code int) *JsonData {
	res.Code = code
	return res
}

func (res *JsonData) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err.Error())
	}
}

func (res *JsonData) Write500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	res.Write(w)
}

func (res *JsonData) Write403(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	res.Write(w)
}

func (res *JsonData) Write401(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	res.Write(w)
}
