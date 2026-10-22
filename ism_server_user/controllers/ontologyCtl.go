package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/ontology"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

type OntologyController struct {
	beego.Controller
}

func ontDB() *gorm.DB { return models.Db }

func (c *OntologyController) Schema() {
	s, err := ontology.LoadSchema()
	if err != nil {
		c.json(500, err.Error())
		return
	}
	filtered := s.Filter(c.GetString("layer"))
	c.Data["json"] = map[string]interface{}{"code": 1000, "schema": filtered}
	c.ServeJSON()
}

func (c *OntologyController) Compile() {
	var body struct {
		GraphKind   string `json:"graph_kind"`
		ProjectUuid string `json:"project_uuid"`
		FocusKey    string `json:"focus_key"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if body.GraphKind == "" {
		body.GraphKind = "world"
	}
	var runID string
	var err error
	if body.GraphKind == "system" {
		runID, err = ontology.CompileSystem(ontDB())
	} else {
		if body.ProjectUuid == "" {
			c.json(422, "缺少 project_uuid")
			return
		}
		runID, err = ontology.CompileWorld(ontDB(), body.ProjectUuid, body.FocusKey)
	}
	if err != nil {
		c.json(500, err.Error())
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "compile_run_id": runID, "graph_kind": body.GraphKind}
	c.ServeJSON()
}

func (c *OntologyController) Graph() {
	kind := c.GetString("graph_kind")
	if kind == "" {
		kind = "world"
	}
	nodes, edges, err := ontology.LoadGraph(ontDB(), ontology.Scope{
		GraphKind: kind, ProjectUuid: c.GetString("project_uuid"), FocusKey: c.GetString("focus_key"),
	})
	if err != nil {
		c.json(500, err.Error())
		return
	}
	view := c.GetString("view")
	if view == "" {
		view = kind
	}
	nodes, edges = ontology.FilterView(view, nodes, edges)
	lint := ontology.Evaluate(nodes, edges)
	if view != "gap" && view != "quarantine" {
		allN, allE, _ := ontology.LoadGraph(ontDB(), ontology.Scope{
			GraphKind: kind, ProjectUuid: c.GetString("project_uuid"), FocusKey: c.GetString("focus_key"),
		})
		lint = ontology.Evaluate(allN, allE)
	}
	c.Data["json"] = map[string]interface{}{
		"code": 1000, "view": view, "nodes": nodes, "edges": edges, "lint": lint,
		"stats": map[string]int{"nodes": len(nodes), "edges": len(edges)},
	}
	c.ServeJSON()
}

func (c *OntologyController) Live() {
	view := ontology.LiveSnapshot(ontDB(), c.GetString("project_uuid"), c.GetString("node_type"), c.GetString("node_key"))
	c.Data["json"] = map[string]interface{}{"code": 1000, "live": view}
	c.ServeJSON()
}

func (c *OntologyController) Lint() {
	kind := c.GetString("graph_kind")
	if kind == "" {
		kind = "world"
	}
	nodes, edges, err := ontology.LoadGraph(ontDB(), ontology.Scope{
		GraphKind: kind, ProjectUuid: c.GetString("project_uuid"), FocusKey: c.GetString("focus_key"),
	})
	if err != nil {
		c.json(500, err.Error())
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "lint": ontology.Evaluate(nodes, edges)}
	c.ServeJSON()
}

func (c *OntologyController) Action() {
	actionID := c.Ctx.Input.Param(":actionId")
	params := map[string]interface{}{}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	res, err := ontology.ExecuteAction(ontDB(), actionID, params, c.actor())
	if err != nil {
		if ae, ok := err.(*ontology.ActionError); ok {
			c.Data["json"] = map[string]interface{}{"code": ae.Status, "detail": ae.Detail}
			c.ServeJSON()
			return
		}
		c.json(500, err.Error())
		return
	}
	if actionID != "VerifyEdge" && actionID != "AckAlarm" {
		_, _ = ontology.CompileWorld(ontDB(), asProject(params), asFocus(params))
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "result": res}
	c.ServeJSON()
}

func (c *OntologyController) actor() string {
	token := c.Ctx.Input.Header("Authorization")
	if token == "" {
		token = c.Ctx.GetCookie("Authorization")
	}
	claims, code := middleware.CheckToken(token)
	if claims != nil && code == 200 && claims.Username != "" {
		return claims.Username
	}
	return "operator"
}

func (c *OntologyController) json(code int, msg string) {
	c.Data["json"] = map[string]interface{}{"code": code, "message": msg}
	c.ServeJSON()
}

func asProject(p map[string]interface{}) string {
	s, _ := p["project_uuid"].(string)
	return s
}
func asFocus(p map[string]interface{}) string {
	s, _ := p["focus_key"].(string)
	return s
}

type AgentRunController struct {
	beego.Controller
}

func (c *AgentRunController) Start() {
	var body struct {
		WorkflowCode string `json:"workflow_code"`
		ProjectUuid  string `json:"project_uuid"`
		FocusKey     string `json:"focus_key"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if body.ProjectUuid == "" {
		c.Data["json"] = map[string]interface{}{"code": 422, "message": "缺少 project_uuid"}
		c.ServeJSON()
		return
	}
	llm := false
	if v, err := config.String("ontology_llm_base_url"); err == nil && strings.TrimSpace(v) != "" {
		llm = true
	}
	out, err := ontology.RunWorkflow(ontDB(), body.WorkflowCode, body.ProjectUuid, body.FocusKey, "operator", llm)
	if err != nil {
		if ae, ok := err.(*ontology.ActionError); ok {
			c.Data["json"] = map[string]interface{}{"code": ae.Status, "detail": ae.Detail}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{"code": 500, "message": err.Error()}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "run": out}
	c.ServeJSON()
}

func (c *AgentRunController) Catalog() {
	view := ontology.AgentCatalog(ontDB(), c.GetString("project_uuid"))
	c.Data["json"] = map[string]interface{}{"code": 1000, "catalog": view}
	c.ServeJSON()
}

func (c *AgentRunController) Talk() {
	var body struct {
		ProjectUuid string `json:"project_uuid"`
		Page        string `json:"page"`
		Text        string `json:"text"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	turn := ontology.XiaoAnTalk(ontDB(), body.ProjectUuid, body.Page, body.Text)
	c.Data["json"] = map[string]interface{}{"code": 1000, "turn": turn}
	c.ServeJSON()
}

func (c *AgentRunController) Get() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	out, err := ontology.GetRun(ontDB(), uint(id))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": 404, "message": "run not found"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 1000, "run": out}
	c.ServeJSON()
}
