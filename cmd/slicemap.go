package cmd

// OrderedMap の代替

type kv struct {
	key string `json:"key"`
	val string `json:"value"`
}

type SliceMap []kv

func NewSliceMap() SliceMap {
	sm := SliceMap{}
	return sm
}

func (sm *SliceMap) exists(key string) bool {
	for _, v := range *sm {
		if v.key == key {
			return true
		}
	}
	return false
}

func (sm *SliceMap) Set(key, val string) {
	for i := range *sm {
		if (*sm)[i].key == key {
			(*sm)[i].val = val
			return
		}
	}
	*sm = append(*sm, kv{key: key, val: val})
}

func (sm *SliceMap) Get(key string) (string, bool) {
	for _, v := range *sm {
		if v.key == key {
			return v.val, true
		}
	}
	return "", false
}
