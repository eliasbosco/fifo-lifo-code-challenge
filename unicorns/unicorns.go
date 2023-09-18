package unicorns

type CapabilitiesList = []string

var Capabilities = CapabilitiesList{}

// LIFO stack
type UnicornElement struct {
	Name         string           `json:"name"`
	Capabilities CapabilitiesList `json:"capabilities"`
}

type UnicornList []UnicornElement

func (u *UnicornList) Push(unicorn *UnicornElement) {
	*u = append(*u, *unicorn)
}

func (u *UnicornList) Pop() *UnicornElement {
	uni := *u
	if len(*u) > 0 {
		res := uni[len(uni)-1]
		*u = uni[:len(uni)-1]
		return &res
	}
	return nil
}
