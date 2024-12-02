package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Date time.Time

func (m *Date) UnmarshalJSON(data []byte) error {
	var strVal string
	err := json.Unmarshal(data, &strVal)
	date, err := time.Parse("2006-01-02", strVal)
	if err != nil {
		return err
	}
	*m = Date(date)
	return nil
}

func (m *Date) MarshalJSON() ([]byte, error) {
	if m != nil {
		return []byte(time.Time(*m).Format("2006-01-02")), nil
	}
	return nil, nil
}

func (m *Date) GetGormDate() *datatypes.Date {
	var date datatypes.Date
	date = datatypes.Date(time.Time(*m))
	return &date
}

func UtilDateValidateValuer(field reflect.Value) interface{} {
	// log.Println("UtilDateValidateValuer")
	// For future use

	// if valuer, ok := field.Interface().(driver.Valuer); ok {

	//   val, err := valuer.Value()
	//   if err == nil {
	//     return val
	//   }
	//   // handle the error how you want
	// }

	return nil
}

type Iso8601 time.Time

func (m *Iso8601) UnmarshalJSON(data []byte) error {
	var strVal string
	err := json.Unmarshal(data, &strVal)
	date, err := time.Parse("2006-01-02T15:04:05Z0700", strVal)
	if err != nil {
		return err
	}
	*m = Iso8601(date)
	return nil
}

func (m *Iso8601) MarshalJSON() ([]byte, error) {
	if m != nil {
		return []byte(time.Time(*m).Format("2006-01-02T15:04:05Z0700")), nil
	}
	return nil, nil
}

func (m *Iso8601) GetGormDate() *datatypes.Date {
	var date datatypes.Date
	date = datatypes.Date(time.Time(*m))
	return &date
}

func (m *Iso8601) GetTime() *time.Time {
	var date time.Time
	date = time.Time(*m)
	return &date
}

func UtilIsoValidateValuer(field reflect.Value) interface{} {
	return nil
}

type RequestModel struct {
	pageNumber int                          `json:"page_number,omitempty"`
	pageSize   int                          `json:"page_size,omitempty"`
	request    *http.Request                `json:"-"`
	filter     *interface{}                 `json:"-"`
	filterKVal map[string]map[string]string `json:"-"`
	orderKVal  map[string]string            `json:"-"`
}

func NewRequestModel(r *http.Request) *RequestModel {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNumber = 0
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 0
	}
	return &RequestModel{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		request:    r,
	}
}

func (r *RequestModel) PageSize() int {
	if r.pageSize > 200 {
		return 200
	} else if r.pageSize > 0 {
		return r.pageSize
	}
	return 20
}

func (r *RequestModel) PageNumber() int {
	if r.pageNumber > 0 {
		return r.pageNumber
	}
	return 1
}

func (r *RequestModel) PaginateScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (r.PageNumber() - 1) * r.PageSize()
		return db.Offset(offset).Limit(r.PageSize())
	}
}

func (r *RequestModel) parseTagMap(field reflect.StructField, mapKey string) (map[string]string, error) {
	tagString := field.Tag.Get(mapKey)
	kval := make(map[string]string)
	if len(tagString) == 0 {
		return kval, errors.New("tag not found")
	}
	var key string
	var val string
	parsingKey := true
	for _, char := range tagString {
		if char == ';' {
			kval[key] = val
			parsingKey = true
			key = ""
			val = ""
			continue
		}
		if char == ':' {
			parsingKey = false
			continue
		}
		if parsingKey {
			key += string(char)
		} else {
			val += string(char)
		}
	}
	if len(key) > 0 {
		kval[key] = val
	}
	return kval, nil
}

func (r *RequestModel) ParseFilter(filter interface{}) {
	r.filterKVal = make(map[string]map[string]string)
	length := reflect.TypeOf(filter).Elem().NumField()
	for i := 0; i < length; i++ {
		field := reflect.TypeOf(filter).Elem().Field(i)
		// log.Println(field)
		if kval, err := r.parseTagMap(field, "filter"); err == nil {
			fieldVal := reflect.ValueOf(filter).Elem().Field(i)
			if fieldVal.IsValid() && fieldVal.CanSet() {
				kvalName := kval["name"]
				strVal := ""
				hasVal := false
				if _, condition := kval["condition"]; condition {
					if hasVal = r.request.URL.Query().Has(kvalName + ".lt"); hasVal {
						strVal = r.request.URL.Query().Get(kvalName + ".lt")
						kval["__condition"] = "lt"
					} else if hasVal = r.request.URL.Query().Has(kvalName + ".lte"); hasVal {
						strVal = r.request.URL.Query().Get(kvalName + ".lte")
						kval["__condition"] = "lte"
					} else if hasVal = r.request.URL.Query().Has(kvalName + ".gt"); hasVal {
						strVal = r.request.URL.Query().Get(kvalName + ".gt")
						kval["__condition"] = "gt"
					} else if hasVal = r.request.URL.Query().Has(kvalName + ".gte"); hasVal {
						strVal = r.request.URL.Query().Get(kvalName + ".gte")
						kval["__condition"] = "gte"
					} else if hasVal = r.request.URL.Query().Has(kvalName + ".not"); hasVal {
						strVal = r.request.URL.Query().Get(kvalName + ".not")
						kval["__condition"] = "not"
					} else {
						strVal = r.request.URL.Query().Get(kvalName)
						hasVal = r.request.URL.Query().Has(kvalName)
					}
				} else {
					strVal = r.request.URL.Query().Get(kvalName)
					hasVal = r.request.URL.Query().Has(kvalName)
				}
				kval["__value"] = strVal
				// log.Println(strVal)
				// log.Println(fieldVal.Type().String())
				if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().String() == "*util.Iso8601" {
					if hasVal {
						date, err := time.Parse("2006-01-02T15:04:05Z0700", strVal)
						if err != nil {
							log.Println(err.Error())
						}
						isoDate := Iso8601(date)
						fieldVal.Set(reflect.ValueOf(&isoDate))
					}
				} else if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().String() == "*string" {
					if hasVal {
						fieldVal.Set(reflect.ValueOf(&strVal))
					}
				} else if fieldVal.Kind() == reflect.String && fieldVal.Type().String() == "string" {
					fieldVal.SetString(strVal)
				} else if fieldVal.Kind() == reflect.Int {
					intVal, err := strconv.Atoi(strVal)
					if err != nil {
						intVal = 0
					}
					fieldVal.SetInt(int64(intVal))
				} else if fieldVal.Kind() == reflect.Float32 {
					float32Val, err := strconv.ParseFloat(strVal, 32)
					if err != nil {
						float32Val = 0
					}
					fieldVal.SetFloat(float32Val)
				} else if fieldVal.Kind() == reflect.Float64 {
					float64Val, err := strconv.ParseFloat(strVal, 64)
					if err != nil {
						float64Val = 0
					}
					fieldVal.SetFloat(float64Val)
				} else if fieldVal.Kind() == reflect.Bool {
					bVal := false
					if strVal == "true" || strVal == "1" {
						bVal = true
					}
					fieldVal.SetBool(bVal)
				}
			}
			kval["__field_name"] = field.Name
			// log.Println(tagString)
			// log.Println(kval)
			r.filterKVal[field.Name] = kval
		}

		if kval, err := r.parseTagMap(field, "order_by"); err == nil {
			fieldVal := reflect.ValueOf(filter).Elem().Field(i)
			r.orderKVal = kval
			strVal := r.request.URL.Query().Get("order")
			hasVal := r.request.URL.Query().Has("order")
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().String() == "*string" {
				if hasVal {
					fieldVal.Set(reflect.ValueOf(&strVal))
				}
			} else if fieldVal.Kind() == reflect.String && fieldVal.Type().String() == "string" {
				fieldVal.SetString(strVal)
			}
		}
	}
	r.filter = &filter
}

func (r *RequestModel) GenerateScopeOrder() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.filter != nil {
			collate := "mn-MN-x-icu"
			isDesc := "asc"
			strVal := r.request.URL.Query().Get("order")
			orderKey := ""
			if strings.HasPrefix(strVal, "title") {
				isDesc = fmt.Sprintf("collate "+"%q "+"asc", collate)
			}
			if strings.HasSuffix(strVal, ".desc") {
				isDesc = "desc NULLS LAST"
				orderKey = strings.TrimSuffix(strVal, ".desc")
			} else if strings.HasSuffix(strVal, ".asc") {
				orderKey = strings.TrimSuffix(strVal, ".asc")
			}
			for k, v := range r.orderKVal {
				if orderKey == k {
					if v == "contents.title" && isDesc == "asc" {
						db.Order(fmt.Sprintf("replace(contents.title collate "+"%q "+",'\"','')", collate) + " " + isDesc)
					} else if v == "contents.title" && isDesc != "asc" {
						db.Order("replace(contents.title,'\"','')" + " " + isDesc)
					} else {
						db.Order(v + " " + isDesc)
					}
				}
			}
		}
		return db
	}
}

func (r *RequestModel) GenerateScopeAndWhere() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.filter != nil {
			for k, v := range r.filterKVal {
				if _, ignore := v["ignore"]; ignore {
					continue
				}
				if len(v["__value"]) == 0 {
					continue
				}
				if _, condition := v["condition"]; condition {
					if field_val, has_field := v["field"]; has_field {
						f := reflect.ValueOf(*r.filter).Elem().FieldByName(k)
						log.Println(f.Kind())
						condition := "="
						log.Println(v["__condition"])
						if v["__condition"] == "lt" {
							condition = "<"
						} else if v["__condition"] == "lte" {
							condition = "<="
						} else if v["__condition"] == "gt" {
							condition = ">"
						} else if v["__condition"] == "gte" {
							condition = ">="
						} else if v["__condition"] == "not" {
							condition = "!="
						}

						if f.Kind() == reflect.Ptr && !f.IsNil() && f.Type().String() == "*util.Iso8601" {
							fVal := *f.Interface().(*Iso8601)
							db.Where(field_val+condition+" ?", fVal.GetGormDate())
						} else if f.Kind() == reflect.Ptr && !f.IsNil() && f.Type().String() == "*string" {
							fVal := *f.Interface().(*string)
							db.Where(field_val+condition+" ?", fVal)
						} else if f.Kind() == reflect.String {
							db.Where(field_val+condition+" ?", f.String())
						} else if f.Kind() == reflect.Int {
							db.Where(field_val+condition+" ?", f.Int())
						} else if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
							db.Where(field_val+condition+" ?", f.Float())
						} else if f.Kind() == reflect.Bool {
							db.Where(field_val+condition+" ?", f.Bool())
						}
					}
					continue
				}
				if _, contains := v["contains"]; contains {
					if field_val, has_field := v["field"]; has_field {
						f := reflect.ValueOf(*r.filter).Elem().FieldByName(k)
						log.Println(f.Kind())
						if f.Kind() == reflect.Ptr && !f.IsNil() && f.Type().String() == "*string" {
							fVal := *f.Interface().(*string)
							db.Where("lower("+field_val+") LIKE lower(?)", "%"+fVal+"%")
						} else if f.Kind() == reflect.String {
							db.Where("lower("+field_val+") LIKE lower(?)", "%"+f.String()+"%")
						} else if f.Kind() == reflect.Int {
							db.Where("lower("+field_val+") LIKE lower(?)", "%"+strconv.FormatInt(f.Int(), 10)+"%")
						} else if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
							db.Where("lower("+field_val+") LIKE lower(?)", "%"+fmt.Sprintf("%.2f", f.Float())+"%")
						}
					}
					continue
				}
				if _, exclude := v["exclude"]; exclude {
					if field_val, has_field := v["field"]; has_field {
						f := reflect.ValueOf(*r.filter).Elem().FieldByName(k)
						log.Println(f.Kind())
						if f.Kind() == reflect.Ptr && !f.IsNil() && f.Type().String() == "*string" {
							fVal := *f.Interface().(*string)
							db.Where("lower("+field_val+") NOT LIKE lower(?)", "%"+fVal+"%")
						} else if f.Kind() == reflect.String {
							db.Where("lower("+field_val+") NOT LIKE lower(?)", "%"+f.String()+"%")
						} else if f.Kind() == reflect.Int {
							db.Where("lower("+field_val+") NOT LIKE lower(?)", "%"+strconv.FormatInt(f.Int(), 10)+"%")
						} else if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
							db.Where("lower("+field_val+") NOT LIKE lower(?)", "%"+fmt.Sprintf("%.2f", f.Float())+"%")
						}
					}
					continue
				}
				// log.Println(k)
				// log.Println(v)
				if field_val, has_field := v["field"]; has_field {
					f := reflect.ValueOf(*r.filter).Elem().FieldByName(k)
					// log.Println(f.IsNil())
					// log.Println(f.Kind())
					if f.Kind() == reflect.Ptr && !f.IsNil() && f.Type().String() == "*string" {
						fVal := *f.Interface().(*string)
						db.Where(field_val+" = ?", fVal)
					} else if f.Kind() == reflect.String {
						db.Where(field_val+" = ?", f.String())
					} else if f.Kind() == reflect.Int {
						db.Where(field_val+" = ?", f.Int())
					} else if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
						db.Where(field_val+" = ?", f.Float())
					} else if f.Kind() == reflect.Bool {
						db.Where(field_val+" = ?", f.Bool())
					}
				}
			}
		}
		return db
	}
}

func (r *RequestModel) GenerateScopeOrWhere() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.filter != nil {
			for k, v := range r.filterKVal {
				if _, ignore := v["ignore"]; ignore {
					continue
				}
				if len(v["__value"]) == 0 {
					continue
				}
				if field_val, has_field := v["field"]; has_field {
					f := reflect.ValueOf(*r.filter).Elem().FieldByName(k)
					if f.Kind() == reflect.Ptr && !f.IsNil() {
						fVal := *f.Interface().(*string)
						db.Or("? = ?", field_val, fVal)
					} else if f.Kind() == reflect.String {
						db.Or("? = ?", field_val, f.String())
					} else if f.Kind() == reflect.Int {
						db.Or("? = ?", field_val, f.Int())
					} else if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
						db.Or("? = ?", field_val, f.Float())
					} else if f.Kind() == reflect.Bool {
						db.Or("? = ?", field_val, f.Bool())
					}
				}
			}
		}
		return db
	}
}
