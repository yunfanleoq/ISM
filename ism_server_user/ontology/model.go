package ontology

import (
	"time"

	"gorm.io/gorm"
)

type OntNode struct {
	gorm.Model
	GraphKind    string `gorm:"index;type:varchar(16)" json:"graph_kind"`
	ProjectUuid  string `gorm:"index;type:varchar(64)" json:"project_uuid"`
	FocusKey     string `gorm:"index;type:varchar(128)" json:"focus_key"`
	Layer        string `gorm:"type:varchar(16)" json:"layer"`
	NodeType     string `gorm:"index;type:varchar(64)" json:"node_type"`
	NodeKey      string `gorm:"index;type:varchar(256)" json:"node_key"`
	Label        string `gorm:"type:varchar(512)" json:"label"`
	Props        string `gorm:"type:text" json:"props"`
	Maturity     string `gorm:"type:varchar(16)" json:"maturity"`
	SourceAgent  string `gorm:"type:varchar(64)" json:"source_agent"`
	HumanLocked  bool   `json:"human_locked"`
	Quarantined  bool   `json:"quarantined"`
	Gap          bool   `json:"gap"`
	CompileRunID string `gorm:"type:varchar(64)" json:"compile_run_id"`
}

func (OntNode) TableName() string { return "ont_node" }

type OntEdge struct {
	gorm.Model
	GraphKind     string   `gorm:"index;type:varchar(16)" json:"graph_kind"`
	ProjectUuid   string   `gorm:"index;type:varchar(64)" json:"project_uuid"`
	FocusKey      string   `gorm:"index;type:varchar(128)" json:"focus_key"`
	Layer         string   `gorm:"type:varchar(16)" json:"layer"`
	Relation      string   `gorm:"index;type:varchar(64)" json:"relation"`
	FromKey       string   `gorm:"index;type:varchar(256)" json:"from"`
	ToKey         string   `gorm:"index;type:varchar(256)" json:"to"`
	Props         string   `gorm:"type:text" json:"props"`
	Note          string   `gorm:"type:varchar(512)" json:"note"`
	SourceAgent   string   `gorm:"type:varchar(64)" json:"source_agent"`
	Confidence    *float64 `json:"confidence"`
	HumanVerified bool     `json:"human_verified"`
	Quarantined   bool     `json:"quarantined"`
	Deleted       bool     `json:"deleted"`
	CompileRunID  string   `gorm:"type:varchar(64)" json:"compile_run_id"`
}

func (OntEdge) TableName() string { return "ont_edge" }

type OntDecision struct {
	gorm.Model
	ProjectUuid string `gorm:"index;type:varchar(64)" json:"project_uuid"`
	FocusKey    string `gorm:"index;type:varchar(128)" json:"focus_key"`
	ActionID    string `gorm:"type:varchar(64)" json:"action_id"`
	Actor       string `gorm:"type:varchar(128)" json:"actor"`
	Outcome     string `gorm:"type:varchar(128)" json:"outcome"`
	Reasoning   string `gorm:"type:text" json:"reasoning"`
	TargetKey   string `gorm:"type:varchar(256)" json:"target_key"`
	Props       string `gorm:"type:text" json:"props"`
}

func (OntDecision) TableName() string { return "ont_decision" }

type AgentSkill struct {
	gorm.Model
	Code            string `gorm:"index;type:varchar(64)" json:"code"`
	Name            string `gorm:"type:varchar(128)" json:"name"`
	Body            string `gorm:"type:text" json:"body"`
	TargetAgentKeys string `gorm:"type:varchar(256)" json:"target_agent_keys"`
	Status          string `gorm:"type:varchar(32)" json:"status"`
}

func (AgentSkill) TableName() string { return "agent_skill" }

type AgentRun struct {
	gorm.Model
	WorkflowCode string `gorm:"index;type:varchar(64)" json:"workflow_code"`
	ProjectUuid  string `gorm:"index;type:varchar(64)" json:"project_uuid"`
	FocusKey     string `gorm:"type:varchar(128)" json:"focus_key"`
	Status       string `gorm:"type:varchar(32)" json:"status"`
	TriggeredBy  string `gorm:"type:varchar(128)" json:"triggered_by"`
	Summary      string `gorm:"type:text" json:"summary"`
}

func (AgentRun) TableName() string { return "agent_run" }

type AgentRunStep struct {
	gorm.Model
	RunID    uint   `gorm:"index" json:"run_id"`
	StepKey  string `gorm:"type:varchar(64)" json:"step_key"`
	AgentKey string `gorm:"type:varchar(64)" json:"agent_key"`
	Status   string `gorm:"type:varchar(32)" json:"status"`
	Message  string `gorm:"type:text" json:"message"`
}

func (AgentRunStep) TableName() string { return "agent_run_step" }

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&OntNode{}, &OntEdge{}, &OntDecision{}, &AgentSkill{}, &AgentRun{}, &AgentRunStep{})
}

type monitorRow struct {
	gorm.Model
	Sid              int32
	Pid              int32
	Name             string
	Type             int
	ProjectUuid      string
	Uuid             string
	DeviceType       int
	Interval         int
	Timeout          int
	FailedTimes      int
	Status           int
	ConfigurationUid string
	PageUUID         string
}

func (monitorRow) TableName() string { return "monitor_list" }

type pointRow struct {
	gorm.Model
	Name        string
	Uuid        string
	ProjectUuid string
	Value       string
	DeviceUuid  string
	DeviceName  string
	DeviceType  int
	DataUnit    string
}

func (pointRow) TableName() string { return "device_real_data" }

type alarmRow struct {
	gorm.Model
	AlarmName    string
	DeviceUuid   string
	ProjectUuid  string
	DeviceName   string
	DataUuid     string
	HappenTime   time.Time
	ClearTime    time.Time
	AlarmLevel   int
	AlarmMessage string
}

func (alarmRow) TableName() string { return "devices_alarm_list" }

type triggerRow struct {
	gorm.Model
	Uuid              string
	TriggerName       string
	ProjectUuid       string
	TriggerDeviceUuid string
	TriggerDataUuid   string
	TriggerAlarmLevel int
}

func (triggerRow) TableName() string { return "alarm_trigger" }

type displayRow struct {
	gorm.Model
	Name            string
	ProjectUuid     string
	DisplayModelUid string
}

func (displayRow) TableName() string { return "display_models" }

type userRow struct {
	ID       uint
	Username string
	Name     string
	Role     string
	Uuid     string
}

func (userRow) TableName() string { return "user" }

type roleRow struct {
	gorm.Model
	RoleName   string
	RoleId     string
	Permission string
}

func (roleRow) TableName() string { return "roles_list" }
