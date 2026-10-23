package models

import "time"

// 公安局动环。表名走 GORM 单数策略（bjga_asset 等）。

type BjgaBlob struct {
	ID   uint   `gorm:"primaryKey"`
	K    string `gorm:"column:k;uniqueIndex;size:64"`
	Body string `gorm:"type:text"`
}

type BjgaRoom struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"index;size:64"`
	Floor    string `gorm:"size:16"`
	RoomType string `gorm:"column:room_type;size:32"`
	Zone     string `gorm:"size:64"`
	CabCount int
	Payload  string `gorm:"type:text"`
}

type BjgaAsset struct {
	ID        uint   `gorm:"primaryKey"`
	Sn        string `gorm:"index;size:64"`
	Room      string `gorm:"index;size:64"`
	Floor     string `gorm:"size:16"`
	Cabinet   string `gorm:"index;size:80"`
	Status    string `gorm:"size:32"`
	UStart    int
	UEnd      int
	UH        int    `gorm:"column:u_h"`
	Cat       string `gorm:"size:32"`
	DevType   string `gorm:"column:dev_type;size:64"`
	Brand     string `gorm:"size:64"`
	ModelName string `gorm:"column:model_name;size:64"`
	Kw        float64
	Org       string `gorm:"size:64"`
	Biz       string `gorm:"size:64"`
	Origin    string `gorm:"size:16"`
	Payload   string `gorm:"type:text"`
}

type BjgaAlarm struct {
	ID      uint   `gorm:"primaryKey"`
	Level   string `gorm:"size:16"`
	Room    string `gorm:"size:64"`
	Floor   string `gorm:"size:16"`
	Obj     string `gorm:"size:128"`
	Time    string `gorm:"size:32"`
	State   string `gorm:"size:32"`
	Payload string `gorm:"type:text"`
}

type BjgaSeries struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"uniqueIndex;size:32"`
	Payload string `gorm:"type:text"`
}

type BjgaThDevice struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:128"`
	Room    string `gorm:"size:64"`
	Payload string `gorm:"type:text"`
}

type BjgaCapacityInput struct {
	ID           uint `gorm:"primaryKey"`
	CoefServer   float64
	CoefStorage  float64
	CoefNet      float64
	CoefFiber    float64
	CoefBattery  float64
	PowerRedline float64
}

type BjgaFlow struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	No      string `gorm:"uniqueIndex;size:32" json:"no"`
	Kind    string `gorm:"column:kind;size:32" json:"type"`
	Dev     string `gorm:"size:64" json:"dev"`
	Name    string `gorm:"size:64" json:"name"`
	FromLoc string `gorm:"column:from_loc;size:64" json:"from"`
	ToLoc   string `gorm:"column:to_loc;size:64" json:"to"`
	ByName  string `gorm:"column:by_name;size:64" json:"by"`
	At      string `gorm:"column:at;size:32" json:"t"`
	Stage   int    `json:"stage"`
	St      string `gorm:"size:16" json:"st"`
	Round   int    `json:"round"`
	Why     string `gorm:"size:128" json:"why"`
}

type BjgaStockDiff struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Code  string `gorm:"size:64" json:"code"`
	Cab   string `gorm:"size:32" json:"cab"`
	Model string `gorm:"size:64" json:"model"`
	Kind  string `gorm:"column:kind;size:32" json:"type"`
	Book  string `gorm:"size:64" json:"book"`
	Real  string `gorm:"size:64" json:"real"`
	At    string `gorm:"column:at;size:32" json:"t"`
	St    string `gorm:"size:16" json:"st"`
}

type BjgaStockBill struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	No     string `gorm:"size:32" json:"no"`
	Kind   string `gorm:"size:16" json:"kind"`
	Sn     string `gorm:"size:64" json:"sn"`
	Status string `gorm:"size:16" json:"status"`
	Note   string `gorm:"size:128" json:"note"`
	At     string `gorm:"size:32" json:"at"`
}

type BjgaMaintPlan struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Device   string `gorm:"size:64" json:"device"`
	Cat      string `gorm:"size:32" json:"cat"`
	Item     string `gorm:"size:64" json:"item"`
	Cycle    string `gorm:"size:32" json:"cycle"`
	NextDue  string `gorm:"size:32" json:"nextDue"`
	Status   string `gorm:"size:16" json:"status"`
	SpecNote string `gorm:"size:64" json:"spec"`
}

type BjgaMeter struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"size:64" json:"name"`
	Category string  `gorm:"size:16" json:"category"`
	Source   string  `gorm:"size:16" json:"source"`
	Value    float64 `json:"value"`
	Unit     string  `gorm:"size:16" json:"unit"`
}

type BjgaMeterRevision struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	MeterID uint    `gorm:"index" json:"meterId"`
	Actor   string  `gorm:"size:64" json:"actor"`
	Before  float64 `json:"before"`
	After   float64 `json:"after"`
	Reason  string  `gorm:"size:128" json:"reason"`
	At      string  `gorm:"size:32" json:"at"`
}

type BjgaHotspot struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Room   string  `gorm:"size:64" json:"room"`
	Cab    string  `gorm:"size:64" json:"cab"`
	Tin    float64 `json:"tin"`
	LoadKw float64 `json:"loadKw"`
	Advice string  `gorm:"size:128" json:"advice"`
}

type BjgaCoolingUnit struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `gorm:"size:64" json:"name"`
	Room       string  `gorm:"size:64" json:"room"`
	ReturnC    float64 `json:"returnC"`
	SupplyC    float64 `json:"supplyC"`
	PowerKw    float64 `json:"powerKw"`
	AirflowM3h float64 `json:"airflowM3h"`
	DesignCop  float64 `json:"designCop"`
}

type BjgaCarbonFactor struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Name   string  `gorm:"size:64" json:"name"`
	Scope  string  `gorm:"size:16" json:"scope"`
	Factor float64 `json:"factor"`
	Unit   string  `gorm:"size:32" json:"unit"`
	Source string  `gorm:"size:64" json:"source"`
	Kwh    float64 `json:"kwh"`
}

type BjgaProtectSetting struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Device  string  `gorm:"size:64" json:"device"`
	Item    string  `gorm:"size:64" json:"item"`
	Current float64 `json:"current"`
	Sheet   float64 `json:"sheet"`
	Last    float64 `json:"last"`
	Unit    string  `gorm:"size:16" json:"unit"`
}

type BjgaPqPoint struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `gorm:"size:64" json:"name"`
	Value float64 `json:"value"`
	Limit float64 `json:"limit"`
	Unit  string  `gorm:"size:16" json:"unit"`
	Std   string  `gorm:"size:32" json:"std"`
}

type BjgaComtrade struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"size:64" json:"name"`
	Conclusion string `gorm:"size:64" json:"conclusion"`
	Payload    string `gorm:"type:text" json:"-"`
}

type BjgaTicket struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	No          string `gorm:"uniqueIndex;size:32" json:"no"`
	Title       string `gorm:"size:128" json:"title"`
	Status      string `gorm:"size:16" json:"status"`
	Device      string `gorm:"size:64" json:"device"`
	DeviceState string `gorm:"size:16" json:"deviceState"`
	Author      string `gorm:"size:32" json:"author"`
}

type BjgaTicketStep struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	TicketNo string `gorm:"index;size:32" json:"ticketNo"`
	Seq      int    `json:"seq"`
	Text     string `gorm:"size:128" json:"text"`
	Done     int    `json:"done"`
}

type BjgaSimCase struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Kind    string  `gorm:"uniqueIndex;size:16" json:"kind"`
	Name    string  `gorm:"size:64" json:"name"`
	Kva     float64 `json:"kva"`
	UkPct   float64 `json:"ukPct"`
	IsetA   float64 `json:"isetA"`
	Note    string  `gorm:"size:128" json:"note"`
	Payload string  `gorm:"type:text" json:"-"`
}

type BjgaAuditLog struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Actor  string `gorm:"size:64" json:"actor"`
	Action string `gorm:"size:64" json:"action"`
	Object string `gorm:"size:128" json:"object"`
	Before string `gorm:"type:text" json:"before"`
	After  string `gorm:"type:text" json:"after"`
	At     string `gorm:"size:32" json:"at"`
}

func (BjgaAuditLog) TableName() string { return "bjga_audit" }

type BjgaSecPolicy struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	MinLen     int  `json:"minLen"`
	MinClasses int  `json:"minClasses"`
	MaxAgeDays int  `json:"maxAgeDays"`
	HistoryN   int  `json:"historyN"`
	FailLimit  int  `json:"failLimit"`
	LockMin    int  `json:"lockMin"`
}

type BjgaCompliance struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Clause string `gorm:"size:64" json:"clause"`
	Item   string `gorm:"size:128" json:"item"`
	Where  string `gorm:"size:128" json:"where"`
	Status string `gorm:"size:16" json:"status"`
}

type BjgaLicense struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"size:64" json:"name"`
	Kind    string `gorm:"size:32" json:"kind"`
	Version string `gorm:"size:32" json:"version"`
	Count   int    `json:"count"`
	CertNo  string `gorm:"size:64" json:"certNo"`
	Expire  string `gorm:"size:32" json:"expire"`
	Place   string `gorm:"size:64" json:"place"`
}

type BjgaQualityRule struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:64" json:"name"`
	Kind   string `gorm:"size:32" json:"kind"`
	Scope  string `gorm:"size:64" json:"scope"`
	Status string `gorm:"size:16" json:"status"`
	Hits   int    `json:"hits"`
}

type BjgaInterface struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:64" json:"name"`
	Direct string `gorm:"size:16" json:"direct"`
	Proto  string `gorm:"size:32" json:"proto"`
	Status string `gorm:"size:16" json:"status"`
	Sample string `gorm:"size:128" json:"sample"`
}

type BjgaAlarmChannel struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:32" json:"name"`
	Status string `gorm:"size:16" json:"status"`
	Target string `gorm:"size:64" json:"target"`
	Remark string `gorm:"size:256" json:"remark"`
}

type BjgaHaNode struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"size:32" json:"name"`
	Role     string  `gorm:"size:16" json:"role"`
	Peer     string  `gorm:"size:32" json:"peer"`
	DelayMs  float64 `json:"delayMs"`
	Status   string  `gorm:"size:32" json:"status"`
	Spec     string  `gorm:"size:128" json:"spec"`
	Duty     string  `gorm:"size:512" json:"duty"`
	Link     string  `gorm:"size:512" json:"link"`
	Kind     string  `gorm:"size:16" json:"kind"`
	UseWhen  string  `gorm:"size:128" json:"useWhen"`
	Online   string  `gorm:"size:256" json:"online"`
	SwitchAt string  `gorm:"size:256" json:"switchAt"`
	Keep     string  `gorm:"size:256" json:"keep"`
}

type BjgaDuty struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Day   string `gorm:"size:16" json:"day"`
	Shift string `gorm:"size:16" json:"shift"`
	Who   string `gorm:"size:64" json:"who"`
	Post  string `gorm:"size:32" json:"post"`
	Phone string `gorm:"size:32" json:"phone"`
	Scope string `gorm:"size:128" json:"scope"`
	Note  string `gorm:"size:256" json:"note"`
	State string `gorm:"size:16" json:"state"`
}

type BjgaPeace struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	At     string `gorm:"size:32" json:"at"`
	To     string `gorm:"size:64" json:"to"`
	Status string `gorm:"size:16" json:"status"`
	Note   string `gorm:"size:256" json:"note"`
	Who    string `gorm:"size:64" json:"who"`
	Items  string `gorm:"type:text" json:"-"`
}

type BjgaPeaceItem struct {
	Room   string `json:"room"`
	Result string `json:"result"`
	Note   string `json:"note"`
}

type BjgaKnowledge struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `gorm:"size:128" json:"title"`
	Kind  string `gorm:"size:32" json:"kind"`
	Scene string `gorm:"size:128" json:"scene"`
	Body  string `gorm:"type:text" json:"body"`
}

type BjgaVideo struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:64" json:"name"`
	Room   string `gorm:"size:64" json:"room"`
	Status string `gorm:"size:16" json:"status"`
}

type BjgaDoor struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"size:64" json:"name"`
	Room   string `gorm:"size:64" json:"room"`
	Status string `gorm:"size:16" json:"status"`
	Last   string `gorm:"size:64" json:"last"`
}

type BjgaComm struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Name    string  `gorm:"size:64" json:"name"`
	Proto   string  `gorm:"size:32" json:"proto"`
	DeviceN int     `json:"deviceN"`
	LossPct float64 `json:"lossPct"`
	DelayMs float64 `json:"delayMs"`
	Status  string  `gorm:"size:16" json:"status"`
}

type BjgaChange struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"size:128" json:"title"`
	Status string `gorm:"size:16" json:"status"`
	At     string `gorm:"size:32" json:"at"`
	Owner  string `gorm:"size:32" json:"owner"`
}

func BjgaAutoMigrate() error {
	return Db.AutoMigrate(
		&BjgaBlob{}, &BjgaRoom{}, &BjgaAsset{}, &BjgaAlarm{}, &BjgaSeries{}, &BjgaThDevice{},
		&BjgaCapacityInput{}, &BjgaFlow{}, &BjgaStockDiff{}, &BjgaStockBill{}, &BjgaMaintPlan{},
		&BjgaMeter{}, &BjgaMeterRevision{}, &BjgaHotspot{}, &BjgaCoolingUnit{}, &BjgaCarbonFactor{},
		&BjgaProtectSetting{}, &BjgaPqPoint{}, &BjgaComtrade{}, &BjgaTicket{}, &BjgaTicketStep{},
		&BjgaSimCase{}, &BjgaAuditLog{}, &BjgaSecPolicy{}, &BjgaCompliance{}, &BjgaLicense{},
		&BjgaQualityRule{}, &BjgaInterface{}, &BjgaAlarmChannel{}, &BjgaHaNode{}, &BjgaDuty{},
		&BjgaPeace{}, &BjgaKnowledge{}, &BjgaVideo{}, &BjgaDoor{}, &BjgaComm{}, &BjgaChange{},
	)
}

func bjgaNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func BjgaNow() string { return bjgaNow() }
