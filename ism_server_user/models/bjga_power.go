package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func BjgaFlowAct(no, action, why, actor string) (BjgaFlow, string) {
	var f BjgaFlow
	if err := Db.Where("no = ?", no).First(&f).Error; err != nil {
		return f, "单据不存在"
	}
	before := f.St
	switch action {
	case "pass":
		if f.St == "待初审" {
			f.Stage = 2
			f.St = "待终审"
		} else if f.St == "待终审" {
			f.St = "已通过"
		} else {
			return f, "当前状态不能通过"
		}
	case "reject":
		if f.St != "待初审" && f.St != "待终审" {
			return f, "当前状态不能驳回"
		}
		f.St = "已驳回"
		if why == "" {
			why = "照片与工单未附齐"
		}
		f.Why = why
	case "redo":
		if f.St != "已驳回" {
			return f, "只有驳回单可以重提"
		}
		f.St = "待初审"
		f.Stage = 1
		f.Round++
	default:
		return f, "未知动作"
	}
	Db.Save(&f)
	BjgaAudit(actor, "资产审批"+action, no, before, f.St)
	return f, ""
}

func BjgaVerifyDiff(id uint, actor string) (BjgaStockDiff, string) {
	var d BjgaStockDiff
	if err := Db.First(&d, id).Error; err != nil {
		return d, "差异不存在"
	}
	before := d.St
	d.St = "已核实"
	Db.Save(&d)
	BjgaAudit(actor, "盘点核实", d.Code, before, "已核实（不直接改台账颜色）")
	return d, ""
}

func BjgaAddBill(kind, sn, note, actor string) BjgaStockBill {
	var n int64
	Db.Model(&BjgaStockBill{}).Count(&n)
	b := BjgaStockBill{
		No: fmt.Sprintf("IO-2026-%03d", n+1), Kind: kind, Sn: sn, Status: "已登记", Note: note, At: bjgaNow(),
	}
	Db.Create(&b)
	BjgaAudit(actor, "出入库", b.No, "", kind+" "+sn)
	return b
}

func BjgaSettings() []map[string]interface{} {
	var rows []BjgaProtectSetting
	Db.Order("id").Find(&rows)
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		mismatch := false
		if r.Sheet != 0 && math.Abs(r.Current-r.Sheet)/math.Abs(r.Sheet) > 0.05 {
			mismatch = true
		}
		out = append(out, map[string]interface{}{
			"id": r.ID, "device": r.Device, "item": r.Item,
			"current": r.Current, "sheet": r.Sheet, "last": r.Last, "unit": r.Unit,
			"mismatch": mismatch,
		})
	}
	return out
}

func BjgaUpdateSetting(id uint, current float64, actor string) string {
	var r BjgaProtectSetting
	if err := Db.First(&r, id).Error; err != nil {
		return "定值不存在"
	}
	before := fmt.Sprint(r.Current)
	r.Last = r.Current
	r.Current = current
	Db.Save(&r)
	BjgaAudit(actor, "定值修改", r.Device+" "+r.Item, before, fmt.Sprint(current))
	return ""
}

func BjgaQuality() []map[string]interface{} {
	var rows []BjgaPqPoint
	Db.Order("id").Find(&rows)
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		over := false
		if r.Name == "功率因数" {
			over = r.Value < r.Limit
		} else if r.Limit > 0 {
			over = r.Value > r.Limit
		}
		out = append(out, map[string]interface{}{
			"id": r.ID, "name": r.Name, "value": r.Value, "limit": r.Limit,
			"unit": r.Unit, "std": r.Std, "over": over,
		})
	}
	return out
}

type bjgaChan struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

func bjgaDFT(samples []float64) (mag, ang float64) {
	n := len(samples)
	if n == 0 {
		return 0, 0
	}
	var re, im float64
	for i, v := range samples {
		th := 2 * math.Pi * float64(i) / float64(n)
		re += v * math.Cos(th)
		im -= v * math.Sin(th)
	}
	re /= float64(n)
	im /= float64(n)
	mag = 2 * math.Hypot(re, im)
	ang = math.Atan2(im, re) * 180 / math.Pi
	return mag, ang
}

func bjgaConclude(phasors []map[string]interface{}) string {
	if len(phasors) < 3 {
		return "通道不足，无法判定"
	}
	mags := make([]float64, len(phasors))
	max := 0.0
	for i, p := range phasors {
		mags[i], _ = p["mag"].(float64)
		if mags[i] > max {
			max = mags[i]
		}
	}
	if max == 0 {
		return "无有效波形"
	}
	low := 0
	for _, m := range mags {
		if m < max*0.4 {
			low++
		}
	}
	switch low {
	case 1:
		return "单相接地（一相幅值明显偏低）"
	case 2:
		return "两相故障或断线"
	case 0:
		return "三相基本平衡"
	default:
		return "三相均跌落，疑似失压"
	}
}

func BjgaAnalyzeChannels(chs []bjgaChan) map[string]interface{} {
	phasors := make([]map[string]interface{}, 0, len(chs))
	for _, c := range chs {
		mag, ang := bjgaDFT(c.Data)
		phasors = append(phasors, map[string]interface{}{
			"name": c.Name, "mag": bjgaRound(mag, 2), "ang": bjgaRound(ang, 1),
		})
	}
	return map[string]interface{}{
		"channels": chs, "phasors": phasors, "conclusion": bjgaConclude(phasors),
		"sequence": bjgaSequence(phasors),
		"demo":     true,
		"demoNote": "演示解析：采样不是现场保护装置导出的 COMTRADE 文件。",
	}
}

func bjgaPhasorC(p map[string]interface{}) (re, im float64) {
	mag, _ := p["mag"].(float64)
	ang, _ := p["ang"].(float64)
	rad := ang * math.Pi / 180
	return mag * math.Cos(rad), mag * math.Sin(rad)
}

func bjgaSequence(phasors []map[string]interface{}) map[string]interface{} {
	if len(phasors) < 3 {
		return map[string]interface{}{}
	}
	ar, ai := -0.5, math.Sqrt(3)/2
	mul := func(re, im, br, bi float64) (float64, float64) {
		return re*br - im*bi, re*bi + im*br
	}
	add3 := func(re, im, r2, i2, r3, i3 float64) (float64, float64) {
		return (re + r2 + r3) / 3, (im + i2 + i3) / 3
	}
	aRe, aIm := bjgaPhasorC(phasors[0])
	bRe, bIm := bjgaPhasorC(phasors[1])
	cRe, cIm := bjgaPhasorC(phasors[2])
	bPosR, bPosI := mul(bRe, bIm, ar, ai)
	cPosR, cPosI := mul(cRe, cIm, ar, -ai)
	bNegR, bNegI := mul(bRe, bIm, ar, -ai)
	cNegR, cNegI := mul(cRe, cIm, ar, ai)
	zRe, zIm := add3(aRe, aIm, bRe, bIm, cRe, cIm)
	pRe, pIm := add3(aRe, aIm, bPosR, bPosI, cPosR, cPosI)
	nRe, nIm := add3(aRe, aIm, bNegR, bNegI, cNegR, cNegI)
	zm := bjgaRound(math.Hypot(zRe, zIm), 2)
	pm := bjgaRound(math.Hypot(pRe, pIm), 2)
	nm := bjgaRound(math.Hypot(nRe, nIm), 2)
	hint := "正序为主，三相基本对称（演示判据）"
	if pm > 0 && nm > pm*0.2 {
		hint = "负序偏大，倾向不对称故障（演示判据）"
	}
	if pm > 0 && zm > pm*0.15 {
		hint = "零序偏大，倾向接地故障（演示判据）"
	}
	return map[string]interface{}{"zero": zm, "positive": pm, "negative": nm, "hint": hint}
}

func BjgaComtradeGet(id uint) map[string]interface{} {
	var row BjgaComtrade
	q := Db.Order("id")
	if id > 0 {
		q = Db.Where("id = ?", id)
	}
	if err := q.First(&row).Error; err != nil {
		return map[string]interface{}{"conclusion": "无录波"}
	}
	var body struct {
		Channels []bjgaChan `json:"channels"`
	}
	_ = json.Unmarshal([]byte(row.Payload), &body)
	out := BjgaAnalyzeChannels(body.Channels)
	out["id"] = row.ID
	out["name"] = row.Name
	return out
}

func BjgaParseComtradeText(name, text string) map[string]interface{} {
	lines := strings.Split(text, "\n")
	var a, b, c []float64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.FieldsFunc(line, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' })
		if len(parts) < 3 {
			continue
		}
		va, e1 := strconv.ParseFloat(parts[len(parts)-3], 64)
		vb, e2 := strconv.ParseFloat(parts[len(parts)-2], 64)
		vc, e3 := strconv.ParseFloat(parts[len(parts)-1], 64)
		if e1 != nil || e2 != nil || e3 != nil {
			continue
		}
		a = append(a, va)
		b = append(b, vb)
		c = append(c, vc)
	}
	if len(a) < 8 {
		return map[string]interface{}{"code": -3, "message": "没有解析到三相采样，请用每行 va,vb,vc"}
	}
	chs := []bjgaChan{{Name: "Va", Data: a}, {Name: "Vb", Data: b}, {Name: "Vc", Data: c}}
	payload := bjgaRaw(map[string]interface{}{"channels": chs})
	row := BjgaComtrade{Name: name, Payload: payload}
	if row.Name == "" {
		row.Name = "上传录波"
	}
	Db.Create(&row)
	out := BjgaAnalyzeChannels(chs)
	out["id"] = row.ID
	out["name"] = row.Name
	Db.Model(&row).Update("conclusion", out["conclusion"])
	return out
}

func BjgaTickets() []map[string]interface{} {
	var tickets []BjgaTicket
	Db.Order("id").Find(&tickets)
	out := make([]map[string]interface{}, 0, len(tickets))
	for _, t := range tickets {
		var steps []BjgaTicketStep
		Db.Where("ticket_no = ?", t.No).Order("seq").Find(&steps)
		out = append(out, map[string]interface{}{
			"id": t.ID, "no": t.No, "title": t.Title, "status": t.Status,
			"device": t.Device, "deviceState": t.DeviceState, "author": t.Author, "steps": steps,
		})
	}
	return out
}

func BjgaTicketAct(no, action, actor string, stepSeq int) (map[string]interface{}, string) {
	var t BjgaTicket
	if err := Db.Where("no = ?", no).First(&t).Error; err != nil {
		return nil, "操作票不存在"
	}
	before := t.Status
	switch action {
	case "approve":
		if t.Status != "review" {
			return nil, "只有待审核票可以批准"
		}
		if t.DeviceState != "运行" && strings.Contains(t.Title, "运行转检修") {
			return nil, "设备当前不是运行状态，不能转检修"
		}
		t.Status = "approved"
	case "execute":
		if t.Status != "approved" && t.Status != "executing" {
			return nil, "请先批准再执行"
		}
		t.Status = "executing"
		var step BjgaTicketStep
		q := Db.Where("ticket_no = ? AND done = 0", t.No).Order("seq")
		if stepSeq > 0 {
			q = Db.Where("ticket_no = ? AND seq = ?", t.No, stepSeq)
		}
		if err := q.First(&step).Error; err != nil {
			t.Status = "archived"
		} else {
			step.Done = 1
			Db.Save(&step)
			var left int64
			Db.Model(&BjgaTicketStep{}).Where("ticket_no = ? AND done = 0", t.No).Count(&left)
			if left == 0 {
				t.Status = "archived"
			}
		}
	default:
		return nil, "未知动作"
	}
	Db.Save(&t)
	BjgaAudit(actor, "操作票"+action, t.No, before, t.Status)
	list := BjgaTickets()
	for _, item := range list {
		if item["no"] == t.No {
			return item, ""
		}
	}
	return nil, ""
}

func BjgaSim() map[string]interface{} {
	var rows []BjgaSimCase
	Db.Order("id").Find(&rows)
	station := map[string]interface{}{}
	for _, r := range rows {
		if r.Kind == "station" {
			calc := bjgaStationCalc(r.Kva, r.UkPct, r.IsetA)
			station = map[string]interface{}{
				"name": r.Name, "kva": r.Kva, "ukPct": r.UkPct, "isetA": r.IsetA,
				"note": r.Note, "calc": calc,
			}
		}
	}
	nodes := make([]map[string]interface{}, 33)
	for i := 0; i < 33; i++ {
		v := 1 - float64(i)*0.0026
		nodes[i] = map[string]interface{}{"n": i + 1, "vPu": bjgaRound(v, 3)}
	}
	return map[string]interface{}{
		"station": station,
		"ieee33": map[string]interface{}{
			"name": "IEEE33 标准算例", "note": "标准算例，非本站拓扑。用于校验算法，不能当作本站单线图。",
			"nodes": nodes,
		},
	}
}

func bjgaStationCalc(kva, uk, iset float64) map[string]interface{} {
	if uk <= 0 {
		uk = 6
	}
	in := kva / (math.Sqrt(3) * 0.4)
	isc := in / (uk / 100)
	sens := 0.0
	if iset > 0 {
		sens = isc / iset
	}
	conclusion := "灵敏度满足（相对 1.5）"
	if sens < 1.5 {
		conclusion = "灵敏度不足（低于 1.5）"
	}
	const busA = 2500.0
	margin := (busA - in) / busA * 100
	ampConc := "相对演示母线仍有裕度"
	if in > busA {
		ampConc = "超过演示母线载流量"
	}
	return map[string]interface{}{
		"inA": bjgaRound(in, 1), "iscA": bjgaRound(isc, 1), "sensitivity": bjgaRound(sens, 2),
		"conclusion": conclusion,
		"model":      "演示模型：只计变压器 uk，未计入电缆阻抗。用于看短路电流和灵敏度怎么变，不能代替短路计算书，也不是现场故障电流。",
		"ampacityA":  busA, "marginPct": bjgaRound(margin, 1), "ampacityConclusion": ampConc,
		"ampacityNote": "母线载流量 2500A 是演示额定，不是电缆实测温升。",
		"continuity":   "模型只有一台变压器。该台退出后没有第二台承接，供电连续性不满足。这是简化模型的结论，不是现场开关位置。",
	}
}

func BjgaSimUpdate(kva float64, actor string) map[string]interface{} {
	var row BjgaSimCase
	if err := Db.Where("kind = ?", "station").First(&row).Error; err != nil {
		return BjgaSim()
	}
	before := fmt.Sprint(row.Kva)
	row.Kva = kva
	Db.Save(&row)
	BjgaAudit(actor, "仿真改容量", "本站变压器", before, fmt.Sprint(kva))
	return BjgaSim()
}
