package HJ212Protocol

import (
	"strconv"
	"strings"
	"time"
)

type Property2011 struct {
	Rtd  float64
	Flag string
}

type Data2021 map[string]Property2011

type HJ212_2011 struct {
	DataTime time.Time
	Data     Data2021
}

// TODO 数据编码
func (entity *HJ212_2011) Encode() ([]byte, error) {
	return nil, nil
}

// 数据解码
func (entity *HJ212_2011) Decode(data string) error {
	// 初始化entity.body
	entity.Data = make(Data2021)

	// 提取CP报文体并重新赋值给data
	start := strings.Index(data, "CP=&&") + len("CP=&&")
	end := strings.LastIndex(data, "&&")
	if start < 0 || end < 0 || end <= start {
		return nil
	}
	data = data[start:end]

	var items []string
	if len(data) > 0 {
		items = strings.Split(data, ";")
	}

	for _, item := range items {
		parts := strings.Split(item, ",")

		// 解析DataTime
		if len(parts) == 1 && strings.Contains(parts[0], "DataTime") {
			kv := strings.Split(parts[0], "=")
			t, err := time.ParseInLocation("20060102150405", kv[1], time.Local)
			if err != nil {
				return err
			}
			entity.DataTime = t
			continue
		}

		// 解析污染物数据
		var key string
		var points Property2011
		for _, part := range parts {

			kv := strings.Split(part, "=")
			if len(kv) == 2 {
				// 提取污染物检测因子名称和度量标记
				keys := strings.Split(kv[0], "-")
				if len(keys) == 2 {
					key = keys[0]

					switch keys[1] {
					case Field_Rtd:
						num, err := strconv.ParseFloat(kv[1], 64)
						if err != nil {
							return err
						}
						points.Rtd = num
					case Field_Flag:
						points.Flag = kv[1]
					}
				}
			}
		}
		entity.Data[key] = points
	}
	return nil
}
