package HJ212Protocol

import "sync"

// 命令编码
type Cn uint16

const (
	// 实时数据
	Cn_2011 Cn = 2011
	// 分钟数据
	Cn_2041 Cn = 2041
	// 分钟数据
	Cn_2051 Cn = 2051
	// 小时数据
	Cn_2061 Cn = 2061
	// 日数据
	Cn_2031 Cn = 2031
)

var entityMapper = map[uint16]Entity{
	uint16(Cn_2011): new(HJ212_2011),
	uint16(Cn_2041): new(HJ212_2011),
	uint16(Cn_2051): new(HJ212_2051),
	uint16(Cn_2061): new(HJ212_2061),
	uint16(Cn_2031): new(HJ212_2031),
}

func RegisterEntity(typ uint16, entity Entity) {
	entityMapper[typ] = entity
}

func RemoveEntity(typ uint16) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := entityMapper[typ]; ok {
		delete(entityMapper, typ)
	}
}
