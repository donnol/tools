package list

// StringList 列表
type StringList []string

// Filter 从列表中过滤指定元素
func (list StringList) Filter(s string) StringList {
	newKeys := make(StringList, 0, len(list))
	for _, key := range list {
		if key == s {
			continue
		}
		newKeys = append(newKeys, key)
	}
	return newKeys
}

// Filter 从列表中过滤指定元素
func Filter(keys []string, s string) []string {
	return StringList(keys).Filter(s)
}
