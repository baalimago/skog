package construct

import (
	"fmt"

	"github.com/baalimago/skog/internal/models"
)

type index struct {
	// each string represents one step further into the tree
	root []string
	// key of the currently selected index
	key string
}

type builder struct {
	data models.JSONLike
	// postion indicates the current location in the json tree
	position index
}

func (b builder) Set(k string, v any) {
	b.data[k] = v
}

func (b builder) Del(k string) {
	delete(b.data, k)
}

func (b builder) CurrentLevel() (models.JSONLike, error) {
	posLevel := b.position.root
	var ret models.JSONLike
	for i, p := range posLevel {
		// Prime the rest of the recursion
		if i == 0 {
			pos, ok := b.data[p].(models.JSONLike)
			if !ok {
				return nil, fmt.Errorf("failed to traverse tree to current level. Position: '%v', index root: '%v'", p, b.position.root)
			}
			ret = pos
			continue
		}
		pos, ok := ret[p].(models.JSONLike)
		if !ok {
			return nil, fmt.Errorf("failed to traverse tree to current level. Position: '%v', index root: '%v'", p, b.position.root)
		}
		ret = pos
	}
	return ret, nil
}

func (b builder) Traverse(path []string) (models.JSONLike, error) {
	ret := b.data
	for _, p := range path {
		pos, ok := ret[p].(models.JSONLike)
		if !ok {
			return nil, fmt.Errorf("failed to traverse tree. Path: '%v'", path)
		}
		ret = pos
	}
	return ret, nil
}

func NewBuilder() builder {
	return builder{
		data: make(models.JSONLike),
	}
}
