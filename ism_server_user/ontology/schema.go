package ontology

import (
	"embed"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed physical_v1.yaml business_v1.yaml system_v1.yaml
var schemaFS embed.FS

type EntityDef struct {
	Type        string `yaml:"type" json:"type"`
	Label       string `yaml:"label" json:"label"`
	Layer       string `yaml:"layer" json:"layer"`
	Maturity    string `yaml:"maturity" json:"maturity"`
	Source      string `yaml:"source" json:"source"`
	Color       string `yaml:"color" json:"color"`
	Description string `yaml:"description" json:"description"`
}

type RelationDef struct {
	Relation string   `yaml:"relation" json:"relation"`
	Label    string   `yaml:"label" json:"label"`
	From     []string `yaml:"from" json:"from"`
	To       []string `yaml:"to" json:"to"`
	Layer    string   `yaml:"layer" json:"layer"`
	Color    string   `yaml:"color" json:"color"`
}

type AxiomDef struct {
	ID           string `yaml:"id" json:"id"`
	Severity     string `yaml:"severity" json:"severity"`
	Label        string `yaml:"label" json:"label"`
	Layer        string `yaml:"layer" json:"layer"`
	BlocksExport bool   `yaml:"blocks_export" json:"blocks_export"`
}

type ActionDef struct {
	ID                string   `yaml:"id" json:"id"`
	Label             string   `yaml:"label" json:"label"`
	Implemented       bool     `yaml:"implemented" json:"implemented"`
	Workflow          string   `yaml:"workflow" json:"workflow"`
	AllowedInterfaces []string `yaml:"allowed_interfaces" json:"allowed_interfaces"`
}

type FunctionDef struct {
	ID             string   `yaml:"id" json:"id"`
	Label          string   `yaml:"label" json:"label"`
	Kind           string   `yaml:"kind" json:"kind"`
	InputTypes     []string `yaml:"input_types" json:"input_types"`
	OutputProps    []string `yaml:"output_props" json:"output_props"`
	TriggersAction string   `yaml:"triggers_action" json:"triggers_action"`
}

type InterfaceDef struct {
	ID              string   `yaml:"id" json:"id"`
	Label           string   `yaml:"label" json:"label"`
	Surface         string   `yaml:"surface" json:"surface"`
	AllowedActions  []string `yaml:"allowed_actions" json:"allowed_actions"`
	WriteJSONDirect bool     `yaml:"write_json_direct" json:"write_json_direct"`
}

type fileSchema struct {
	Version        string         `yaml:"version"`
	Domain         string         `yaml:"domain"`
	Layer          string         `yaml:"layer"`
	Description    string         `yaml:"description"`
	EntityTypes    []EntityDef    `yaml:"entity_types"`
	RelationTypes  []RelationDef  `yaml:"relation_types"`
	Axioms         []AxiomDef     `yaml:"axioms"`
	ActionTypes    []ActionDef    `yaml:"action_types"`
	FunctionTypes  []FunctionDef  `yaml:"function_types"`
	InterfaceTypes []InterfaceDef `yaml:"interface_types"`
}

type Schema struct {
	Version        string         `json:"version"`
	Domain         string         `json:"domain"`
	Description    string         `json:"description"`
	EntityTypes    []EntityDef    `json:"entity_types"`
	RelationTypes  []RelationDef  `json:"relation_types"`
	Axioms         []AxiomDef     `json:"axioms"`
	ActionTypes    []ActionDef    `json:"action_types"`
	FunctionTypes  []FunctionDef  `json:"function_types"`
	InterfaceTypes []InterfaceDef `json:"interface_types"`
}

var (
	schemaOnce sync.Once
	schemaVal  *Schema
	schemaErr  error
)

func LoadSchema() (*Schema, error) {
	schemaOnce.Do(func() {
		schemaVal, schemaErr = loadSchema()
	})
	return schemaVal, schemaErr
}

func loadSchema() (*Schema, error) {
	out := &Schema{Version: "v1", Domain: "dcim", Description: "DCM 三层世界模型"}
	for _, name := range []string{"physical_v1.yaml", "business_v1.yaml", "system_v1.yaml"} {
		raw, err := schemaFS.ReadFile(name)
		if err != nil {
			return nil, err
		}
		var f fileSchema
		if err := yaml.Unmarshal(raw, &f); err != nil {
			return nil, err
		}
		out.EntityTypes = append(out.EntityTypes, f.EntityTypes...)
		out.RelationTypes = append(out.RelationTypes, f.RelationTypes...)
		out.Axioms = append(out.Axioms, f.Axioms...)
		out.ActionTypes = append(out.ActionTypes, f.ActionTypes...)
		out.FunctionTypes = append(out.FunctionTypes, f.FunctionTypes...)
		out.InterfaceTypes = append(out.InterfaceTypes, f.InterfaceTypes...)
	}
	return out, nil
}

func (s *Schema) Filter(layer string) Schema {
	if layer == "" || layer == "all" {
		return *s
	}
	out := Schema{Version: s.Version, Domain: s.Domain, Description: s.Description}
	for _, e := range s.EntityTypes {
		if e.Layer == layer {
			out.EntityTypes = append(out.EntityTypes, e)
		}
	}
	for _, r := range s.RelationTypes {
		if r.Layer == layer {
			out.RelationTypes = append(out.RelationTypes, r)
		}
	}
	for _, a := range s.Axioms {
		if a.Layer == layer {
			out.Axioms = append(out.Axioms, a)
		}
	}
	if layer == "business" || layer == "" {
		out.ActionTypes = s.ActionTypes
		out.FunctionTypes = s.FunctionTypes
		out.InterfaceTypes = s.InterfaceTypes
	}
	return out
}

func (s *Schema) Action(id string) *ActionDef {
	for i := range s.ActionTypes {
		if s.ActionTypes[i].ID == id {
			return &s.ActionTypes[i]
		}
	}
	return nil
}
