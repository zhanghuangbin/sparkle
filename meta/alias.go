package meta

import "strings"

type EnvType int

const (
	OVERWRITE EnvType = iota
	APPEND    EnvType = iota
)

type Alias struct {
	Alias    string  `json:"alias"`
	Desc     string  `json:"desc"`
	LongDesc string  `json:"longDesc"`
	Type     EnvType `json:"type"`
	Key      string  `json:"key"`
	Value    string  `json:"value"`
}

type AliasList []Alias

func (list *AliasList) Add(alias *Alias) bool {
	idx := list.IndexOf(alias.Alias)
	if idx == -1 {
		// 添加alias到列表
		*list = append(*list, *alias)
		return true
	}
	// 修改alias
	(*list)[idx] = *alias
	return false
}

func (list *AliasList) IndexOf(alias string) int {
	for idx, v := range *list {
		if v.Alias == alias {
			return idx
		}
	}
	return -1
}

func (list *AliasList) Get(alias string) *Alias {
	for _, v := range *list {
		if v.Alias == alias {
			return &v
		}
	}
	return nil
}

func (list *AliasList) Query(alias string, key string) []Alias {
	var result = make([]Alias, 0, len(*list))
	for _, v := range *list {
		if (alias == "" || v.Alias == alias) && (key == "" || strings.Contains(v.Key, key)) {
			result = append(result, v)
		}
	}
	return result
}

func (list *AliasList) Remove(alias string) bool {
	idx := list.IndexOf(alias)
	if idx == -1 {
		return false
	}
	*list = append((*list)[:idx], (*list)[idx+1:]...)
	return true
}
