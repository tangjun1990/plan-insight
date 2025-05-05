package response

type Page struct {
	PageNum    int   `json:"pageNum" form:"pageNum"`       // 分页-页码
	PageSize   int   `json:"pageSize" form:"pageSize"`     // 分页-页数
	TotalCount int64 `json:"totalCount" form:"totalCount"` // 总数
}
