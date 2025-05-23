package raymond

import (
	"strings"
)

type contextMember struct {
	path      string
	asMapping string
}

type handlebarsContext struct {
	contextMembers []contextMember
}

func newHandlebarsContext() *handlebarsContext {
	var cm []contextMember
	return &handlebarsContext{contextMembers: cm}
}

func (h *handlebarsContext) AddMemberContext(path, asMapping string) {
	cmp := contextMember{path: path, asMapping: asMapping}
	h.contextMembers = append(h.contextMembers, cmp)
}

func (h *handlebarsContext) GetCurrentContext() []string {
	return h.GetParentContext(0)
}

func (h *handlebarsContext) GetCurrentContextString() string {
	return h.GetParentContextString(0)
}

func (h *handlebarsContext) GetParentContext(numAncestors int) []string {
	if len(h.contextMembers) == 0 {
		return []string{}
	}
	return strings.Split(h.GetParentContextString(numAncestors), ".")
}

func (h *handlebarsContext) GetParentContextString(numAncestors int) string {
	if len(h.contextMembers) == 0 {
		return ""
	}
	if numAncestors > len(h.contextMembers) {
		numAncestors = 0
	}
	var res string
	for _, val := range h.contextMembers[:len(h.contextMembers)-numAncestors] {
		if len(res) == 0 {
			res = val.path
		} else {
			res = res + "." + val.path
		}
	}
	return res
}

func (h *handlebarsContext) MoveUpContext() {
	if len(h.contextMembers) > 0 {
		h.contextMembers = h.contextMembers[:len(h.contextMembers)-1]
	}
}

func (h *handlebarsContext) HaveAsContexts(numAncestors int) bool {
	if numAncestors > len(h.contextMembers) {
		numAncestors = 0
	}
	for val := range h.contextMembers[:len(h.contextMembers)-numAncestors] {
		if h.contextMembers[val].asMapping != "" {
			return true
		}
	}
	return false
}

func (h *handlebarsContext) GetMappedContext(path []string, numAncestors int) []string {
	if len(path) == 0 {
		return []string{}
	}
	return strings.Split(h.GetMappedContextString(path, numAncestors), ".")
}

func (h *handlebarsContext) GetMappedContextString(path []string, numAncestors int) string {
	if len(h.contextMembers) == 0 {
		return strings.Join(path, ".")
	}
	if numAncestors > len(h.contextMembers) {
		numAncestors = 0
	}
	if !h.HaveAsContexts(numAncestors) {
		var res string
		if path[0] == "" {
			res = h.GetParentContextString(numAncestors)
		} else {
			res = h.GetParentContextString(numAncestors) + "." + strings.Join(path, ".")
		}
		return strings.Trim(res, ".")
	}
	var res string
	copiedMembers := make([]contextMember, 0)
	if numAncestors > 0 {
		copiedMembers = append(copiedMembers, h.contextMembers[:len(h.contextMembers)-numAncestors]...)
	} else {
		copiedMembers = append(copiedMembers, h.contextMembers...)
	}
	for p := len(path) - 1; p >= 0; p-- {
		var val contextMember
		var found string
		if len(copiedMembers) == 0 {
			found = path[p]
		} else {
			val = copiedMembers[len(copiedMembers)-1]
			if val.asMapping == path[p] {
				found = val.path
				if len(copiedMembers) > 1 {
					val2 := copiedMembers[len(copiedMembers)-2]
					tmp := strings.Split(val.path, ".")
					if tmp[0] == val2.asMapping {
						found = strings.Join(tmp[1:], ".")
					}
				}
				copiedMembers = copiedMembers[:len(copiedMembers)-1]
			} else {
				if len(val.asMapping) == 0 {
					found = val.path + "." + path[p]
					copiedMembers = copiedMembers[:len(copiedMembers)-1]
				} else {
					if len(copiedMembers) == 0 {
						ss := strings.Split(val.asMapping, ".")
						if ss[0] == path[p] {
							found = val.path
						}
					} else {
						if len(copiedMembers) > 1 {
							cv := copiedMembers[len(copiedMembers)-2]
							mappedPath := strings.Split(cv.path, ".")
							if len(mappedPath) > 1 {
								tmp := strings.Join(mappedPath[1:], ".")
								if tmp == val.asMapping {
									found = val.path
									copiedMembers = copiedMembers[:len(copiedMembers)-1]
								} else {
									found = path[p]
								}
							} else {
								found = path[p]
							}
						} else {
							found = path[p]
						}
					}
				}
			}
		}
		res = found + "." + res
	}
	if len(copiedMembers) > 0 {
		for p := len(copiedMembers) - 1; p >= 0; p-- {
			res = copiedMembers[p].path + "." + res
		}
	}
	return strings.Trim(res, ".")
}
