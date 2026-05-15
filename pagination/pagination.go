package pagination

const DefaultPageSize = 20
const maxInt32 = 1<<31 - 1

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

func normalize(page, pageSize int32) (int32, int32) {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}
	return page, pageSize
}

func ParseInput(page, pageSize int32) PageInput {
	page, pageSize = normalize(page, pageSize)
	return PageInput{
		Limit:  int(pageSize),
		Offset: int((page - 1) * pageSize),
	}
}

func CalcOutput(page, pageSize int32, total int64) PageOutput {
	page, pageSize = normalize(page, pageSize)

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

func clampInt32(value int64) int32 {
	if value <= 0 {
		return 0
	}
	if value > maxInt32 {
		return maxInt32
	}
	return int32(value)
}
