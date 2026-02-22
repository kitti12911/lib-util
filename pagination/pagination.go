package pagination

const DefaultPageSize = 20

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

func ParseInput(page, pageSize int32) PageInput {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}

	return PageInput{
		Limit:  int(pageSize),
		Offset: int((page - 1) * pageSize),
	}
}

func CalcOutput(page, pageSize int32, total int64) PageOutput {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}

	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	return PageOutput{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalSize:  int32(total),
	}
}
