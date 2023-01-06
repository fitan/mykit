package mycrud

import "strings"

type CrudConfig struct {
	SortableColumns      []string
	PreloadEnableColumns []string
	FilterableColumns    map[string]EnableFilterColumnsOption
	SelectableColumns    []string
	MaxLimit             int
}

func (c CrudConfig) EnableLimit(limit int) bool {
	if c.MaxLimit == 0 {
		return true
	}
	if limit > c.MaxLimit {
		return false
	}
	return true
}

func (c CrudConfig) EnableSelect(column string) bool {
	for _, v := range c.SelectableColumns {
		if strings.HasPrefix(v, column) {
			return true
		}
	}
	return false
}

func (c CrudConfig) EnableSort(column string) bool {
	for _, v := range c.SortableColumns {
		if v == column {
			return true
		}
	}
	return false
}

func (c CrudConfig) EnablePreload(column string) bool {
	for _, v := range c.PreloadEnableColumns {
		if strings.HasPrefix(v, column) {
			return true
		}
	}
	return false
}

func (c CrudConfig) EnableFilter(column string, option string) bool {
	for k, v := range c.FilterableColumns {
		if k == column {
			switch option {
			case "=":
				return v.Equal
			case "!=":
				return v.NotEqual
			case ">":
				return v.GreaterThan
			case "<":
				return v.LessThan
			case ">=":
				return v.GreaterThanOrEqual
			case "<=":
				return v.LessThanOrEqual
			case "?=":
				return v.In
			case "!?=":
				return v.NotIn
			case "~=":
				return v.Like
			case "!~=":
				return v.NotLike
			case "><":
				return v.Between
			case "!><":
				return v.NotBetween
			}
		}
	}
	return false
}

type EnableFilterColumnsOption struct {
	// "="
	Equal bool
	// "!="
	NotEqual bool
	// ">"
	GreaterThan bool
	// ">="
	GreaterThanOrEqual bool
	// "<"
	LessThan bool
	// "<="
	LessThanOrEqual bool
	// "like"
	Like bool
	// "not like"
	NotLike bool
	// "in"
	In bool
	// "not in"
	NotIn bool
	// "between"
	Between bool
	// "not between"
	NotBetween bool
}
