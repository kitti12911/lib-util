package pagination

const DefaultPageSize = 20
const maxInt32 = 1<<31 - 1

type Request struct {
	Page     int32
	PageSize int32
}

type PageInput struct {
	Limit  int
	Offset int
}

type PageOutput struct {
	Page       int32
	PageSize   int32
	TotalPages int32
	TotalSize  int32
}

func (r Request) Input() PageInput {
	if r.Page <= 0 && r.PageSize <= 0 {
		return PageInput{}
	}
	page, pageSize := normalize(r.Page, r.PageSize)
	return PageInput{
		Limit:  int(pageSize),
		Offset: int((page - 1) * pageSize),
	}
}

func (r Request) Output(total int64) PageOutput {
	if r.Page <= 0 && r.PageSize <= 0 {
		totalPages := int32(0)
		if total > 0 {
			totalPages = 1
		}
		return PageOutput{
			Page:       1,
			PageSize:   clampInt32(total),
			TotalPages: totalPages,
			TotalSize:  clampInt32(total),
		}
	}

	page, pageSize := normalize(r.Page, r.PageSize)

	totalPages := int64(0)
	if total > 0 {
		totalPages = (total-1)/int64(pageSize) + 1
	}

	return PageOutput{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: clampInt32(totalPages),
		TotalSize:  clampInt32(total),
	}
}

func normalize(page, pageSize int32) (normalizedPage, normalizedPageSize int32) {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}
	return page, pageSize
}

func ParseInput(page, pageSize int32) PageInput {
	return Request{Page: page, PageSize: pageSize}.Input()
}

func CalcOutput(page, pageSize int32, total int64) PageOutput {
	return Request{Page: page, PageSize: pageSize}.Output(total)
}

func clampInt32(value int64) int32 {
	if value <= 0 {
		return 0
	}
	if value > maxInt32 {
		return maxInt32
	}
	return int32(value)
}
