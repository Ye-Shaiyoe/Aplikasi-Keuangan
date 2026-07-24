package util

// PageResult holds metadata for a paginated response.
type PageResult struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// Paginate normalises page / limit values and returns offset + a PageResult.
//
//   page        – requested page number (1-indexed); clamped to >= 1
//   limit       – requested page size; clamped to [1, maxLimit]
//   defaultLim  – limit used when the caller passes 0 or a negative value
//   maxLimit    – hard upper bound on page size
//
// Returns (normalisedPage, normalisedLimit, offset).
func Paginate(page, limit, defaultLim, maxLimit int) (int, int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultLim
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

// BuildPageResult constructs a PageResult from normalised pagination values.
func BuildPageResult(page, limit, totalItems int) PageResult {
	totalPages := 0
	if limit > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}
	return PageResult{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// HasNextPage returns true when there are more pages after the current one.
func HasNextPage(page, limit, totalItems int) bool {
	return page*limit < totalItems
}

// HasPrevPage returns true when the caller is past the first page.
func HasPrevPage(page int) bool {
	return page > 1
}
