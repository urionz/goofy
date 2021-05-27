package web

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/db"
	"github.com/urionz/goofy/pagination"
	"github.com/urionz/goutil/strutil"
)

type QueryParams struct {
	iris.Context
	db.SqlCnd
}

func NewQueryParams(ctx iris.Context) *QueryParams {
	return &QueryParams{
		Context: ctx,
	}
}

func (q *QueryParams) getValueByColumn(column string, defCol ...string) string {
	if q.Context == nil {
		return ""
	}
	def := ""
	if len(defCol) > 0 && defCol[0] != "" {
		def = defCol[0]
	}
	fieldName := strutil.ToLowerCamel(column)
	switch q.Context.Method() {
	case "GET":
		return q.Context.URLParamDefault(column, def)
	default:
		return q.Context.FormValueDefault(fieldName, def)
	}
}

func (q *QueryParams) In(column string, args interface{}) *QueryParams {
	q.SqlCnd.In(column, args)
	return q
}

func (q *QueryParams) Eq(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.Eq(column, args...)
	return q
}

func (q *QueryParams) EqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Eq(column, value)
	}
	return q
}

func (q *QueryParams) AscByReq(defCol ...string) *QueryParams {
	value := q.getValueByColumn("asc", defCol...)
	if len(value) > 0 {
		q.Asc(value)
	}
	return q
}

func (q *QueryParams) DescByReq(defCol ...string) *QueryParams {
	value := q.getValueByColumn("desc", defCol...)
	if len(value) > 0 {
		q.Desc(value)
	}
	return q
}

func (q *QueryParams) DateByReq(column, operator string, def ...string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.WhereDate(column, operator, value)
	}
	if len(def) > 0 {
		q.WhereDate(column, operator, def[0])
	}
	return q
}

func (q *QueryParams) EqByReqs(columns []string) *QueryParams {
	for _, column := range columns {
		value := q.getValueByColumn(column)
		if len(value) > 0 {
			q.Eq(column, value)
		}
	}
	return q
}

func (q *QueryParams) NotEq(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.NotEq(column, args...)
	return q
}

func (q *QueryParams) NotEqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.NotEq(column, value)
	}
	return q
}

func (q *QueryParams) GtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gt(column, value)
	}
	return q
}

func (q *QueryParams) Gt(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.Gt(column, args...)
	return q
}

func (q *QueryParams) GteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gte(column, value)
	}
	return q
}

func (q *QueryParams) Gte(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.Gte(column, args...)
	return q
}

func (q *QueryParams) LtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lt(column, value)
	}
	return q
}

func (q *QueryParams) Lt(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.Lt(column, args...)
	return q
}

func (q *QueryParams) LteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lte(column, value)
	}
	return q
}

func (q *QueryParams) Lte(column string, args ...interface{}) *QueryParams {
	q.SqlCnd.Lte(column, args...)
	return q
}

func (q *QueryParams) LikeByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Like(column, value)
	}
	return q
}

func (q *QueryParams) Like(column string, str string) *QueryParams {
	q.SqlCnd.Like(column, str)
	return q
}

func (q *QueryParams) LikeByReqLeft(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Starting(column, value)
	}
	return q
}

func (q *QueryParams) LikeByRight(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Ending(column, value)
	}
	return q
}

func (q *QueryParams) Cols(columns []string) *QueryParams {
	q.SelectCols = columns
	return q
}

func (q *QueryParams) PageByReq() *QueryParams {
	if q.Context == nil {
		return q
	}
	paging := GetPaging(q.Context)
	q.Page(paging.Page, paging.Limit, paging.Last)
	return q
}

func (q *QueryParams) Asc(column string) *QueryParams {
	q.Orders = append(q.Orders, pagination.OrderByCol{Column: column, Asc: true})
	return q
}

func (q *QueryParams) Desc(column string) *QueryParams {
	q.Orders = append(q.Orders, pagination.OrderByCol{Column: column, Asc: false})
	return q
}

func (q *QueryParams) Limit(limit int, start ...int) *QueryParams {
	q.Page(1, limit, start...)
	return q
}

func (q *QueryParams) Page(page, limit int, last ...int) *QueryParams {
	lastId := 0
	if len(last) > 0 && last[0] != 0 {
		lastId = last[0]
	}
	if q.Paging == nil {
		q.Paging = &pagination.Paging{Page: page, Limit: limit, Last: lastId}
	} else {
		q.Paging.Page = page
		q.Paging.Limit = limit
		q.Paging.Last = lastId
	}
	return q
}
