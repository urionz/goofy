package db

import (
	"github.com/urionz/goofy/pagination"
	"gorm.io/gorm"
)

type SqlCnd struct {
	SelectCols []string                // 要查询的字段，如果为空，表示查询所有字段
	Params     []pagination.ParamPair  // 参数
	Orders     []pagination.OrderByCol // 排序
	Paging     *pagination.Paging      // 分页
}

// var _ contracts.SqlCondition = (*SqlCnd)(nil)

func NewSqlCnd() *SqlCnd {
	return &SqlCnd{}
}

// func (s *SqlCnd) AppendOrders(orders ...pagination.OrderByCol) {
// 	s.Orders = append(s.Orders, orders...)
// }

func (s *SqlCnd) Cols(selectCols ...string) *SqlCnd {
	if len(selectCols) > 0 {
		s.SelectCols = append(s.SelectCols, selectCols...)
	}
	return s
}

func (s *SqlCnd) WhereDate(column, operator, date string) *SqlCnd {
	s.Where("cast("+column+" as date) "+operator+" ?", date)
	return s
}

func (s *SqlCnd) Eq(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" = ?", args)
	return s
}

func (s *SqlCnd) IsNull(column string) *SqlCnd {
	s.Where(column + " is null")
	return s
}

func (s *SqlCnd) NotEq(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" <> ?", args)
	return s
}

func (s *SqlCnd) Gt(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" > ?", args)
	return s
}

func (s *SqlCnd) Gte(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" >= ?", args)
	return s
}

func (s *SqlCnd) Lt(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" < ?", args)
	return s
}

func (s *SqlCnd) Lte(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" <= ?", args)
	return s
}

func (s *SqlCnd) Like(column string, str string) *SqlCnd {
	s.Where(column+" LIKE ?", "%"+str+"%")
	return s
}

func (s *SqlCnd) Starting(column string, str string) *SqlCnd {
	s.Where(column+" LIKE ?", str+"%")
	return s
}

func (s *SqlCnd) Ending(column string, str string) *SqlCnd {
	s.Where(column+" LIKE ?", "%"+str)
	return s
}

func (s *SqlCnd) In(column string, params interface{}) *SqlCnd {
	s.Where(column+" in (?) ", params)
	return s
}

func (s *SqlCnd) Where(query string, args ...interface{}) *SqlCnd {
	s.Params = append(s.Params, pagination.ParamPair{Query: query, Args: args})
	return s
}

func (s *SqlCnd) Asc(column string) *SqlCnd {
	s.Orders = append(s.Orders, pagination.OrderByCol{Column: column, Asc: true})
	return s
}

func (s *SqlCnd) Desc(column string) *SqlCnd {
	s.Orders = append(s.Orders, pagination.OrderByCol{Column: column, Asc: false})
	return s
}

func (s *SqlCnd) Limit(limit int) *SqlCnd {
	s.Page(1, limit)
	return s
}

func (s *SqlCnd) Page(page, limit int) *SqlCnd {
	if s.Paging == nil {
		s.Paging = &pagination.Paging{Page: page, Limit: limit}
	} else {
		s.Paging.Page = page
		s.Paging.Limit = limit
	}
	return s
}

func (s *SqlCnd) Build(db *gorm.DB) *gorm.DB {
	ret := db

	if len(s.SelectCols) > 0 {
		ret = ret.Select(s.SelectCols)
	}

	// where
	if len(s.Params) > 0 {
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	// order
	if len(s.Orders) > 0 {
		for _, order := range s.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	// limit
	if s.Paging != nil && s.Paging.Limit > 0 {
		ret = ret.Limit(s.Paging.Limit)
	}

	// offset
	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}
	return ret
}

func (s *SqlCnd) Find(db *gorm.DB, out interface{}) error {
	if err := s.Build(db).Find(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *SqlCnd) FindOne(db *gorm.DB, out interface{}) error {
	if err := s.Limit(1).Build(db).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *SqlCnd) Count(db *gorm.DB, model interface{}) (int64, error) {
	ret := db.Model(model)

	// where
	if len(s.Params) > 0 {
		for _, query := range s.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	var count int64
	if err := ret.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
