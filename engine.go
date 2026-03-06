package cache

type Engine interface {
	Set(id any, v []byte)
	SetForever(id any, v []byte)
	SetJson(id, v any)
	Get(id any) (res []byte)
	Has(id any) (has bool)
	Delete(id any)
	Reset()
	SetInt(id any, v int)
	GetInt(id any) (has bool, v int)
	SetInts(id any, v []int)
	GetInts(id any) (has bool, v []int)
	Batch(ids ...any) (res []*KV, err error)
}
