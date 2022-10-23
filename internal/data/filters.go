package data

import "richwynmorris.co.uk/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "maximum per page is 10,000,000.")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero.")
	v.Check(f.PageSize <= 10_000_000, "page_size", "maximum page size is 10,000,000")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
