package models

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	gormlog "github.com/beego/beego/v2/core/logs"
)

func bjgaDCIMPath() string {
	for _, p := range []string{"data/bjga/dcim.json", "ism_server_user/data/bjga/dcim.json"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/bjga/dcim.json"
}

func bjgaStr(m map[string]interface{}, k string) string {
	if m == nil || m[k] == nil {
		return ""
	}
	return fmt.Sprint(m[k])
}

func bjgaInt(m map[string]interface{}, k string) int {
	switch v := m[k].(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func bjgaFloat(m map[string]interface{}, k string) float64 {
	switch v := m[k].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		return 0
	}
}

func bjgaRaw(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func bjgaWipe() {
	models := []interface{}{
		&BjgaAsset{}, &BjgaAlarm{}, &BjgaRoom{}, &BjgaSeries{}, &BjgaThDevice{}, &BjgaBlob{},
		&BjgaFlow{}, &BjgaStockDiff{}, &BjgaStockBill{}, &BjgaMaintPlan{}, &BjgaMeter{},
		&BjgaMeterRevision{}, &BjgaHotspot{}, &BjgaCoolingUnit{}, &BjgaCarbonFactor{},
		&BjgaProtectSetting{}, &BjgaPqPoint{}, &BjgaComtrade{}, &BjgaTicket{}, &BjgaTicketStep{},
		&BjgaSimCase{}, &BjgaCapacityInput{}, &BjgaSecPolicy{}, &BjgaCompliance{}, &BjgaLicense{},
		&BjgaQualityRule{}, &BjgaInterface{}, &BjgaAlarmChannel{}, &BjgaHaNode{}, &BjgaDuty{},
		&BjgaPeace{}, &BjgaKnowledge{}, &BjgaVideo{}, &BjgaDoor{}, &BjgaComm{}, &BjgaChange{},
	}
	for _, m := range models {
		Db.Where("1 = 1").Delete(m)
	}
}

// BjgaEnsureSeed 表为空时导入原型 JSON，并写入招标缺口的种子数据。
func BjgaEnsureSeed() error {
	var ready BjgaBlob
	Db.Where("k = ?", "seed_ready").First(&ready)
	if ready.ID > 0 {
		return nil
	}
	var n int64
	Db.Model(&BjgaAsset{}).Count(&n)
	if n > 0 {
		bjgaWipe()
	}
	path := bjgaDCIMPath()
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err = json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	if err = bjgaSeedScreen(doc); err != nil {
		return err
	}
	if err = bjgaSeedExtra(); err != nil {
		return err
	}
	if err = Db.Create(&BjgaBlob{K: "seed_ready", Body: bjgaNow()}).Error; err != nil {
		return err
	}
	gormlog.Info("bjga 种子数据已导入")
	return nil
}

func bjgaSeedScreen(doc map[string]interface{}) error {
	for _, key := range []string{"meta", "capacity", "assetStats", "alarmStats", "floors"} {
		if doc[key] == nil {
			continue
		}
		if err := Db.Create(&BjgaBlob{K: key, Body: bjgaRaw(doc[key])}).Error; err != nil {
			return err
		}
	}
	realN := 31
	if meta, ok := doc["meta"].(map[string]interface{}); ok && bjgaInt(meta, "realSeed") > 0 {
		realN = bjgaInt(meta, "realSeed")
	}
	rooms, _ := doc["rooms"].([]interface{})
	batchRooms := make([]BjgaRoom, 0, len(rooms))
	for _, item := range rooms {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		batchRooms = append(batchRooms, BjgaRoom{
			Name: bjgaStr(m, "name"), Floor: bjgaStr(m, "floor"), RoomType: bjgaStr(m, "type"),
			Zone: bjgaStr(m, "zone"), CabCount: bjgaInt(m, "cabCount"), Payload: bjgaRaw(m),
		})
	}
	if len(batchRooms) > 0 {
		if err := Db.CreateInBatches(batchRooms, 20).Error; err != nil {
			return err
		}
	}
	assets, _ := doc["assets"].([]interface{})
	batch := make([]BjgaAsset, 0, 200)
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		err := Db.CreateInBatches(batch, 100).Error
		batch = batch[:0]
		return err
	}
	for i, item := range assets {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		origin := "synthetic"
		if i < realN {
			origin = "real"
		}
		m["origin"] = origin
		batch = append(batch, BjgaAsset{
			Sn: bjgaStr(m, "sn"), Room: bjgaStr(m, "room"), Floor: bjgaStr(m, "floor"),
			Cabinet: bjgaStr(m, "cabinet"), Status: bjgaStr(m, "status"),
			UStart: bjgaInt(m, "uStart"), UEnd: bjgaInt(m, "uEnd"), UH: bjgaInt(m, "uH"),
			Cat: bjgaStr(m, "cat"), DevType: bjgaStr(m, "type"), Brand: bjgaStr(m, "brand"),
			ModelName: bjgaStr(m, "model"), Kw: bjgaFloat(m, "kw"), Org: bjgaStr(m, "org"),
			Biz: bjgaStr(m, "biz"), Origin: origin, Payload: bjgaRaw(m),
		})
		if len(batch) >= 200 {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	if err := flush(); err != nil {
		return err
	}
	alarms, _ := doc["alarms"].([]interface{})
	ab := make([]BjgaAlarm, 0, 200)
	for _, item := range alarms {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		ab = append(ab, BjgaAlarm{
			Level: bjgaStr(m, "level"), Room: bjgaStr(m, "room"), Floor: bjgaStr(m, "floor"),
			Obj: bjgaStr(m, "obj"), Time: bjgaStr(m, "time"), State: bjgaStr(m, "state"),
			Payload: bjgaRaw(m),
		})
		if len(ab) >= 200 {
			if err := Db.CreateInBatches(ab, 100).Error; err != nil {
				return err
			}
			ab = ab[:0]
		}
	}
	if len(ab) > 0 {
		if err := Db.CreateInBatches(ab, 100).Error; err != nil {
			return err
		}
	}
	if series, ok := doc["series"].(map[string]interface{}); ok {
		for name, body := range series {
			if err := Db.Create(&BjgaSeries{Name: name, Payload: bjgaRaw(body)}).Error; err != nil {
				return err
			}
		}
	}
	ths, _ := doc["thDevices"].([]interface{})
	tb := make([]BjgaThDevice, 0, len(ths))
	for _, item := range ths {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		tb = append(tb, BjgaThDevice{Name: bjgaStr(m, "name"), Room: bjgaStr(m, "room"), Payload: bjgaRaw(m)})
	}
	if len(tb) > 0 {
		if err := Db.CreateInBatches(tb, 50).Error; err != nil {
			return err
		}
	}
	return Db.Create(&BjgaCapacityInput{
		CoefServer: 1, CoefStorage: 1, CoefNet: 0.5, CoefFiber: 0, CoefBattery: 0, PowerRedline: 85,
	}).Error
}

func bjgaRows(payloads []string) []interface{} {
	out := make([]interface{}, 0, len(payloads))
	for _, p := range payloads {
		var v interface{}
		if json.Unmarshal([]byte(p), &v) == nil {
			out = append(out, v)
		}
	}
	return out
}

func bjgaRound(v float64, n int) float64 {
	p := math.Pow(10, float64(n))
	return math.Round(v*p) / p
}
