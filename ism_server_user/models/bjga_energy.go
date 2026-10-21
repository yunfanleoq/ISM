package models

import "fmt"

func BjgaMeters() []BjgaMeter {
	var rows []BjgaMeter
	Db.Order("id").Find(&rows)
	return rows
}

func BjgaMeterRevisions(limit int) []BjgaMeterRevision {
	if limit <= 0 {
		limit = 20
	}
	var rows []BjgaMeterRevision
	Db.Order("id desc").Limit(limit).Find(&rows)
	return rows
}

func BjgaUpdateMeter(id uint, value float64, reason, actor string) (BjgaMeter, string) {
	var m BjgaMeter
	if err := Db.First(&m, id).Error; err != nil {
		return m, "计量点不存在"
	}
	if reason == "" {
		reason = "人工修订"
	}
	before := m.Value
	Db.Create(&BjgaMeterRevision{
		MeterID: m.ID, Actor: actor, Before: before, After: value, Reason: reason, At: bjgaNow(),
	})
	m.Value = value
	if m.Source == "auto" {
		m.Source = "manual"
	}
	Db.Save(&m)
	BjgaAudit(actor, "计量点修订", m.Name, fmt.Sprint(before), fmt.Sprint(value)+" "+reason)
	return m, ""
}

func BjgaPUE() map[string]interface{} {
	live := BjgaLive()
	return map[string]interface{}{
		"live": live, "meters": BjgaMeters(), "revisions": BjgaMeterRevisions(20),
		"parts": []map[string]interface{}{
			{"name": "IT", "kw": live["it"]},
			{"name": "制冷", "kw": live["cooling"]},
			{"name": "供配电损耗", "kw": live["loss"]},
			{"name": "照明及其他", "kw": live["other"]},
		},
	}
}

func BjgaAirflow() map[string]interface{} {
	var rows []BjgaHotspot
	Db.Order("tin desc").Find(&rows)
	hot := make([]BjgaHotspot, 0)
	for _, r := range rows {
		if r.Tin > 27 {
			hot = append(hot, r)
		}
	}
	return map[string]interface{}{"all": rows, "hot": hot, "limitC": 27}
}

func BjgaCooling() map[string]interface{} {
	var u BjgaCoolingUnit
	if err := Db.First(&u).Error; err != nil {
		return map[string]interface{}{}
	}
	// 近似冷量：1.2 kg/m3 * 1.006 kJ/kg·K * 风量 * 温差
	q := 1.2 * 1.006 * (u.AirflowM3h / 3600) * (u.ReturnC - u.SupplyC)
	cop := 0.0
	if u.PowerKw > 0 {
		cop = q / u.PowerKw
	}
	return map[string]interface{}{
		"unit": u,
		"qKw":  bjgaRound(q, 2),
		"cop":  bjgaRound(cop, 2),
		"process": []map[string]interface{}{
			{"name": "回风", "t": u.ReturnC},
			{"name": "混风", "t": bjgaRound(u.ReturnC-1, 1)},
			{"name": "盘管后", "t": u.SupplyC},
			{"name": "送风", "t": u.SupplyC},
		},
		"note": "COP 由送回风温差与风量估算，不是现场冷量表。",
	}
}

func BjgaCarbon() map[string]interface{} {
	live := BjgaLive()
	totalKw, _ := live["total"].(float64)
	annualKwh := totalKw * 8760
	var factors []BjgaCarbonFactor
	Db.Order("id").Find(&factors)
	var scope2, scope1 float64
	items := make([]map[string]interface{}, 0, len(factors))
	for _, f := range factors {
		kg := 0.0
		if f.Scope == "范围二" {
			kg = annualKwh * f.Factor
			scope2 = kg
		} else {
			kg = f.Kwh * f.Factor
			scope1 = kg
		}
		items = append(items, map[string]interface{}{
			"name": f.Name, "scope": f.Scope, "factor": f.Factor, "unit": f.Unit,
			"source": f.Source, "kg": bjgaRound(kg, 1),
		})
	}
	return map[string]interface{}{
		"annualKwh": bjgaRound(annualKwh, 0),
		"scope1kg":  bjgaRound(scope1, 1),
		"scope2kg":  bjgaRound(scope2, 1),
		"totalT":    bjgaRound((scope1+scope2)/1000, 2),
		"items":     items,
		"advice":    "提高冷冻水供水温度、延长自然冷却，并淘汰低负载 UPS。",
	}
}
