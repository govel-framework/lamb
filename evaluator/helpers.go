package evaluator

import "strings"

func lookForConfigKeys(m map[interface{}]interface{}, key string) (exists bool, value interface{}) {
	split := strings.Split(key, ".")

	if len(split) == 0 {
		return false, ""
	}

	if len(split) == 1 {
		value, ok := m[split[0]]

		return ok, value
	}

	value, ok := m[split[0]]

	if !ok {
		return false, split[0]
	}

	submap, ok := value.(map[interface{}]interface{})

	if !ok {
		return false, split[0]
	}

	return lookForConfigKeys(submap, strings.Join(split[1:], "."))
}
