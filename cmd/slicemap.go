package cmd

type kv struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

type SliceMap []kv

func NewSliceMap() SliceMap {
	sm := SliceMap{}
	return sm
}

func (sm *SliceMap) Set(key, val string) {
	for i := range *sm {
		if (*sm)[i].Key == key {
			(*sm)[i].Val = val
			return
		}
	}
	*sm = append(*sm, kv{Key: key, Val: val})
}

func (sm *SliceMap) Get(key string) (string, bool) {
	for _, v := range *sm {
		if v.Key == key {
			return v.Val, true
		}
	}
	return "", false
}
