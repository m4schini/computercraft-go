package computer

import "fmt"

type Block map[string]interface{}

func (b Block) String() string {
	return fmt.Sprintf("%s[age=%v,tags=%v]", b.Name(), b.Age(), len(b.Tags()))
}

func (b Block) Name() string {
	name, ok := b["name"]
	if !ok {
		return ""
	}

	nameStr, ok := name.(string)
	if !ok {
		return ""
	}

	return nameStr
}

func (b Block) Age() int {
	state, ok := b["state"]
	if !ok {
		return -1
	}

	stateMap, ok := state.(map[string]interface{})
	if !ok {
		return -1
	}

	age, ok := stateMap["age"]
	if !ok {
		return -1
	}

	ageFloat, ok := age.(float64)
	if !ok {
		return -1
	}

	return int(ageFloat)
}

func (b Block) Tags() []string {
	tagsList := make([]string, 0)

	tags, ok := b["tags"]
	if !ok {
		return tagsList
	}

	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return tagsList
	}

	for k, _ := range tagsMap {
		tagsList = append(tagsList, k)
	}

	return tagsList
}

func (b Block) ContainsTag(tag string) bool {
	tags, ok := b["tags"]
	if !ok {
		return false
	}

	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return false
	}

	_, contains := tagsMap[tag]
	return contains
}
