package pagination

import "math"

const PageSizeDefault = 20

type Pagination struct {
	page     int
	pageSize int
	total    int
}

type Base struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

type Rsp struct {
	Base
	List interface{} `json:"list"`
}

func New(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = PageSizeDefault
	}

	return &Pagination{
		page:     page,
		pageSize: pageSize,
	}
}

func (p *Pagination) SetTotal(total int64) {
	p.total = int(total)
}

func (p *Pagination) GetOffset() int {
	return (p.page - 1) * p.pageSize
}

func (p *Pagination) GetLimit() int {
	return p.pageSize
}

func (p *Pagination) GetPageTotal() int {
	return int(math.Ceil(float64(p.total) / float64(p.pageSize)))
}

// GetBase 返回分页基础信息
func (p *Pagination) GetBase() Base {
	return Base{
		Page:      p.page,
		PageSize:  p.pageSize,
		Total:     p.total,
		TotalPage: p.GetPageTotal(),
	}
}

// Response 分页整体响应
func (p *Pagination) Response(list interface{}) Rsp {
	return Rsp{
		Base: p.GetBase(),
		List: list,
	}
}
