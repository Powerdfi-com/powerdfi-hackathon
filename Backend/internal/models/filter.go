package models

const (
	DefaultPage        = 1
	PageLimit8         = 8
	PageLimit9         = 9
	DefaultSearchLimit = 10
	Range30            = 720
	// Range30            = 30
	MaxPageLimit = 20
)

type Filter struct {
	Page       int
	Limit      int
	Search     string
	Range      uint
	Blockchain string
}

// Offset calculates and returns the position of the data to start selecting from.
// It depends on the Page and Limit.
func (f Filter) Offset() int {

	if f.Page == 0 || f.Page == 1 {
		return 0
	}
	return ((f.Page - 1) * f.Limit)
}
