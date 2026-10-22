package ontology

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Scope struct {
	GraphKind   string
	ProjectUuid string
	FocusKey    string
}

func Persist(db *gorm.DB, scope Scope, d *Draft) (string, error) {
	runID := uuid.NewString()
	var oldNodes []OntNode
	if err := db.Where("graph_kind = ? AND project_uuid = ? AND focus_key = ?", scope.GraphKind, scope.ProjectUuid, scope.FocusKey).Find(&oldNodes).Error; err != nil {
		return "", err
	}
	oldByKey := map[string]OntNode{}
	for _, n := range oldNodes {
		oldByKey[n.NodeKey] = n
	}
	for key, n := range d.Nodes {
		props := propsJSON(n.Props)
		if old, ok := oldByKey[key]; ok {
			if old.HumanLocked {
				continue
			}
			oldProps := parseProps(old.Props)
			if fmt.Sprint(oldProps["state"]) == "acked" {
				n.Props["state"] = "acked"
				props = propsJSON(n.Props)
			}
			old.Layer = n.Layer
			old.NodeType = n.NodeType
			old.Label = n.Label
			old.Props = props
			old.Maturity = n.Maturity
			old.Gap = n.Gap
			old.Quarantined = n.Quarantined
			old.CompileRunID = runID
			if err := db.Save(&old).Error; err != nil {
				return "", err
			}
			continue
		}
		row := OntNode{
			GraphKind: scope.GraphKind, ProjectUuid: scope.ProjectUuid, FocusKey: scope.FocusKey,
			Layer: n.Layer, NodeType: n.NodeType, NodeKey: key, Label: n.Label,
			Props: props, Maturity: n.Maturity, SourceAgent: n.SourceAgent,
			HumanLocked: n.HumanLocked, Quarantined: n.Quarantined, Gap: n.Gap, CompileRunID: runID,
		}
		if err := db.Create(&row).Error; err != nil {
			return "", err
		}
	}
	for _, old := range oldNodes {
		if old.HumanLocked {
			continue
		}
		if _, ok := d.Nodes[old.NodeKey]; !ok {
			if err := db.Unscoped().Delete(&OntNode{}, old.ID).Error; err != nil {
				return "", err
			}
		}
	}

	var oldEdges []OntEdge
	if err := db.Where("graph_kind = ? AND project_uuid = ? AND focus_key = ? AND deleted = ?", scope.GraphKind, scope.ProjectUuid, scope.FocusKey, false).Find(&oldEdges).Error; err != nil {
		return "", err
	}
	edgeKey := func(from, rel, to string) string { return from + "|" + rel + "|" + to }
	oldEdge := map[string]OntEdge{}
	for _, e := range oldEdges {
		oldEdge[edgeKey(e.FromKey, e.Relation, e.ToKey)] = e
	}
	seen := map[string]struct{}{}
	for _, e := range d.Edges {
		k := edgeKey(e.From, e.Relation, e.To)
		seen[k] = struct{}{}
		if old, ok := oldEdge[k]; ok {
			if old.HumanVerified {
				continue
			}
			old.Layer = e.Layer
			old.Note = e.Note
			old.SourceAgent = e.SourceAgent
			old.Quarantined = e.Quarantined
			old.CompileRunID = runID
			if err := db.Save(&old).Error; err != nil {
				return "", err
			}
			continue
		}
		row := OntEdge{
			GraphKind: scope.GraphKind, ProjectUuid: scope.ProjectUuid, FocusKey: scope.FocusKey,
			Layer: e.Layer, Relation: e.Relation, FromKey: e.From, ToKey: e.To,
			Note: e.Note, SourceAgent: e.SourceAgent, HumanVerified: e.HumanVerified,
			Quarantined: e.Quarantined, CompileRunID: runID,
		}
		if err := db.Create(&row).Error; err != nil {
			return "", err
		}
	}
	for k, old := range oldEdge {
		if _, ok := seen[k]; ok || old.HumanVerified {
			continue
		}
		old.Deleted = true
		if err := db.Save(&old).Error; err != nil {
			return "", err
		}
	}
	return runID, nil
}

func LoadGraph(db *gorm.DB, scope Scope) ([]OntNode, []OntEdge, error) {
	var nodes []OntNode
	var edges []OntEdge
	q := db.Where("graph_kind = ? AND project_uuid = ? AND focus_key = ?", scope.GraphKind, scope.ProjectUuid, scope.FocusKey)
	if err := q.Find(&nodes).Error; err != nil {
		return nil, nil, err
	}
	if err := db.Where("graph_kind = ? AND project_uuid = ? AND focus_key = ? AND deleted = ?", scope.GraphKind, scope.ProjectUuid, scope.FocusKey, false).Find(&edges).Error; err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}
