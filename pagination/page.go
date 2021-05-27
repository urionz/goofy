package pagination

import "database/sql"

type Paging struct {
	Page  int   `json:"page"`   // 页码
	Limit int   `json:"limit"`  // 每页条数
	Total int64 `json:"total"`  // 总数据条数
	MaxId int   `json:"max_id"` // 最大位置
	MinId int   `json:"min_id"` // 最小位置
}

func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}

func (p *Paging) TotalPage() int {
	if p.Total == 0 || p.Limit == 0 {
		return 0
	}
	totalPage := int(p.Total) / p.Limit
	if int(p.Total)%p.Limit > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

type ParamPair struct {
	Query string        // 查询
	Args  []interface{} // 参数
}

type OrderByCol struct {
	Column string
	Asc    bool
}

type PageResult struct {
	*Paging
	Results interface{} `json:"results"`
}

type CursorResult struct {
	Results interface{} `json:"results"`
	Cursor  string      `json:"cursor"`
}

func SqlNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) > 0,
	}
}
