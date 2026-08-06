package mechanic

import (
	"strconv"
	"strings"
	"time"
)

type GridMeta struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ActionMeta struct {
	PK        string `form:"pk"`
	FK        string `form:"fk"`
	UK        string `form:"uk"`
	DateNow   string `form:"date_now"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	Search    string `form:"search"`
	Page      int    `form:"page"`
	Size      int    `form:"size"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
	Permanent string `form:"permanent"`
}

func NullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}

func NullableFloat(value float64) any {
	if value == 0 {
		return nil
	}
	return value
}

func CheckMeta(page, size int) (int, int, error) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	if page < 1 {
		return 0, 0, ValidationError("Page must be greater than zero")
	}
	if size < 1 || size > 100 {
		return 0, 0, ValidationError("Page size must be between 1 and 100")
	}
	return page, size, nil
}

func BuildMeta(page, size, total int) GridMeta {
	totalPages := 1
	if total > 0 {
		totalPages = (total + size - 1) / size
	}
	return GridMeta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
}

func ParseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func ParseDuration(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	if strings.HasSuffix(s, "ms") {
		v, _ := strconv.ParseInt(strings.TrimSuffix(s, "ms"), 10, 64)
		return v
	}
	if strings.HasSuffix(s, "s") {
		v, _ := strconv.ParseInt(strings.TrimSuffix(s, "s"), 10, 64)
		return v
	}
	if strings.HasSuffix(s, "m") {
		v, _ := strconv.ParseInt(strings.TrimSuffix(s, "m"), 10, 64)
		return v * 60
	}
	if strings.HasSuffix(s, "h") {
		v, _ := strconv.ParseInt(strings.TrimSuffix(s, "h"), 10, 64)
		return v * 3600
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func ParseDateTime(s string) time.Time {
	s = strings.TrimSpace(s)
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05 MST",
		"02 Jan 2006 15:04:05 MST",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
