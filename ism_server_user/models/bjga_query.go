package models

import (
	"encoding/json"
	"math"
)

func BjgaAudit(actor, action, object, before, after string) {
	if actor == "" {
		actor = "系统"
	}
	_ = Db.Create(&BjgaAuditLog{
		Actor: actor, Action: action, Object: object, Before: before, After: after, At: bjgaNow(),
	}).Error
}

func bjgaBlob(key string) interface{} {
	var row BjgaBlob
	if err := Db.Where("k = ?", key).First(&row).Error; err != nil || row.Body == "" {
		return nil
	}
	var v interface{}
	if json.Unmarshal([]byte(row.Body), &v) != nil {
		return nil
	}
	return v
}

func bjgaPayloads(table string) []string {
	type row struct{ Payload string }
	var rows []row
	Db.Table(table).Order("id").Find(&rows)
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Payload)
	}
	return out
}

// BjgaBootstrap 拼出原型 window.DCIM 的形状，并附上 live / extra。
func BjgaBootstrap() map[string]interface{} {
	dcim := map[string]interface{}{}
	for _, key := range []string{"meta", "capacity", "assetStats", "alarmStats", "floors"} {
		if v := bjgaBlob(key); v != nil {
			dcim[key] = v
		}
	}
	dcim["rooms"] = bjgaRows(bjgaPayloads("bjga_room"))
	dcim["assets"] = bjgaRows(bjgaPayloads("bjga_asset"))
	dcim["alarms"] = bjgaRows(bjgaPayloads("bjga_alarm"))
	dcim["thDevices"] = bjgaRows(bjgaPayloads("bjga_th_device"))
	series := map[string]interface{}{}
	var ss []BjgaSeries
	Db.Order("id").Find(&ss)
	for _, s := range ss {
		var v interface{}
		if json.Unmarshal([]byte(s.Payload), &v) == nil {
			series[s.Name] = v
		}
	}
	dcim["series"] = series
	live := BjgaLive()
	devs, als := BjgaCollectLive(bjgaCollectProject)
	live["devices"] = devs
	live["alarms"] = als
	dcim["live"] = live
	dcim["extra"] = BjgaExtra()
	return dcim
}

func BjgaExtra() map[string]interface{} {
	var flows []BjgaFlow
	var diffs []BjgaStockDiff
	var bills []BjgaStockBill
	var maint []BjgaMaintPlan
	var ha []BjgaHaNode
	var videos []BjgaVideo
	var doors []BjgaDoor
	var comms []BjgaComm
	Db.Order("id").Find(&flows)
	Db.Order("id").Find(&diffs)
	Db.Order("id desc").Find(&bills)
	Db.Order("id").Find(&maint)
	Db.Order("id").Find(&ha)
	Db.Order("id").Find(&videos)
	Db.Order("id").Find(&doors)
	Db.Order("id").Find(&comms)
	return map[string]interface{}{
		"flows": flows, "diffs": diffs, "bills": bills, "maint": maint,
		"ha": ha, "haPlan": BjgaHaPlan(), "video": videos, "doors": doors, "comms": comms,
	}
}

func BjgaLive() map[string]interface{} {
	var meters []BjgaMeter
	Db.Order("id").Find(&meters)
	var it, cool, loss, other, ups, energy, manual, pueBase float64
	src := map[string]string{}
	for _, m := range meters {
		switch m.Category {
		case "it":
			it += m.Value
		case "cooling":
			cool += m.Value
		case "loss":
			loss += m.Value
		case "other":
			other += m.Value
		case "ups":
			ups += m.Value
		case "energy":
			energy += m.Value
		}
		if m.Category == "it" || m.Category == "cooling" || m.Category == "loss" || m.Category == "other" {
			if prev, ok := src[m.Category]; !ok || prev == "" {
				src[m.Category] = m.Source
			} else if prev != m.Source {
				src[m.Category] = "mixed"
			}
		}
		if m.Category != "energy" && m.Category != "ups" {
			pueBase += m.Value
			if m.Source == "manual" || m.Source == "estimate" {
				manual += m.Value
			}
		}
	}
	pue := 0.0
	if it > 0 {
		pue = bjgaRound(pueBase/it, 2)
	}
	ratio := 0.0
	if pueBase > 0 {
		ratio = manual / pueBase * 100
	}
	cabRate := 97.96
	if meta, ok := bjgaBlob("meta").(map[string]interface{}); ok {
		total := bjgaFloat(meta, "cabTotal")
		used := bjgaFloat(meta, "cabUsed")
		if total > 0 {
			cabRate = used / total * 100
		}
	}
	parts := []map[string]interface{}{
		{"name": "IT", "kw": it, "source": src["it"]},
		{"name": "制冷", "kw": cool, "source": src["cooling"]},
		{"name": "供配电损耗", "kw": loss, "source": src["loss"]},
		{"name": "照明及其他", "kw": other, "source": src["other"]},
	}
	return map[string]interface{}{
		"pue": pue, "upsPower": ups, "upsEnergy": energy, "itPower": it,
		"cooling": cool, "loss": loss, "other": other,
		"manualRatio": bjgaRound(ratio, 1), "cabRate": bjgaRound(cabRate, 2),
		"it": it, "total": pueBase,
		"parts": parts,
		// 种子分项不是总进线电表。有总表接入后再改为 true。
		"totalMeter": false,
	}
}

func bjgaCoef(cat string, in BjgaCapacityInput) float64 {
	switch cat {
	case "服务器":
		return in.CoefServer
	case "存储设备":
		return in.CoefStorage
	case "交换机", "网络安全设备", "网络":
		return in.CoefNet
	case "光纤柜":
		return in.CoefFiber
	case "电池柜":
		return in.CoefBattery
	default:
		return in.CoefServer
	}
}

func bjgaLv(comp, pow float64, dev int, redline float64) string {
	if dev == 0 {
		return "off"
	}
	if pow > redline || comp > 90 {
		return "red"
	}
	if comp >= 70 {
		return "yellow"
	}
	return "blue"
}

// BjgaRecalcCapacity 按台账重算机柜空间/电力/综合（取极高），并回写容量摘要。
func BjgaRecalcCapacity() (map[string]interface{}, error) {
	var in BjgaCapacityInput
	if err := Db.First(&in).Error; err != nil {
		in = BjgaCapacityInput{CoefServer: 1, CoefStorage: 1, CoefNet: 0.5, PowerRedline: 85}
	}
	if in.PowerRedline == 0 {
		in.PowerRedline = 85
	}
	type agg struct {
		u, n int
		kw   float64
		coef float64
	}
	by := map[string]*agg{}
	var assets []BjgaAsset
	Db.Find(&assets)
	for _, a := range assets {
		if a.Cabinet == "" || a.Status == "下架待处理" {
			continue
		}
		g := by[a.Cabinet]
		if g == nil {
			g = &agg{}
			by[a.Cabinet] = g
		}
		g.u += a.UH
		g.n++
		g.kw += a.Kw
		g.coef += bjgaCoef(a.Cat, in)
	}
	var rooms []BjgaRoom
	Db.Order("id").Find(&rooms)
	var spaceSum, powSum, ifSum, compMax float64
	var cabN int
	for i := range rooms {
		var room map[string]interface{}
		if json.Unmarshal([]byte(rooms[i].Payload), &room) != nil {
			continue
		}
		cabs, _ := room["cabinets"].([]interface{})
		if len(cabs) == 0 {
			continue
		}
		for _, item := range cabs {
			cab, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name := bjgaStr(cab, "name")
			g := by[name]
			rated := bjgaFloat(cab, "rated")
			if rated <= 0 {
				rated = 4
			}
			ifRate := bjgaFloat(cab, "ifRate")
			u, n := 0, 0
			kw, coef := 0.0, 0.0
			if g != nil {
				u, n, kw, coef = g.u, g.n, g.kw, g.coef
			}
			space := 0.0
			if n > 0 {
				space = (float64(u) + coef) / 42 * 100
			}
			pow := kw / rated * 100
			comp := math.Max(space, math.Max(pow, ifRate))
			cab["dev"] = n
			cab["usedU"] = u
			cab["kw"] = bjgaRound(kw, 2)
			cab["rtKw"] = bjgaRound(kw, 2)
			cab["spaceRate"] = bjgaRound(space, 1)
			cab["powRate"] = bjgaRound(pow, 1)
			cab["comp"] = bjgaRound(comp, 1)
			cab["lv"] = bjgaLv(comp, pow, n, in.PowerRedline)
			spaceSum += space
			powSum += pow
			ifSum += ifRate
			if comp > compMax {
				compMax = comp
			}
			cabN++
		}
		b, err := json.Marshal(room)
		if err != nil {
			continue
		}
		rooms[i].Payload = string(b)
		Db.Model(&BjgaRoom{}).Where("id = ?", rooms[i].ID).Update("payload", rooms[i].Payload)
	}
	cap := map[string]interface{}{}
	if old, ok := bjgaBlob("capacity").(map[string]interface{}); ok {
		cap = old
	}
	if cabN > 0 {
		cap["spaceRate"] = bjgaRound(spaceSum/float64(cabN), 2)
		cap["powRate"] = bjgaRound(powSum/float64(cabN), 2)
		cap["ifRate"] = bjgaRound(ifSum/float64(cabN), 2)
		cap["comp"] = bjgaRound(compMax, 2)
	}
	body := bjgaRaw(cap)
	var blob BjgaBlob
	if err := Db.Where("k = ?", "capacity").First(&blob).Error; err != nil {
		Db.Create(&BjgaBlob{K: "capacity", Body: body})
	} else {
		Db.Model(&blob).Update("body", body)
	}
	return cap, nil
}
