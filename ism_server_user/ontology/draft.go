package ontology

import "encoding/json"

type DraftNode struct {
	Layer       string
	NodeType    string
	NodeKey     string
	Label       string
	Maturity    string
	SourceAgent string
	Props       map[string]interface{}
	Quarantined bool
	Gap         bool
	HumanLocked bool
}

type DraftEdge struct {
	Layer         string
	Relation      string
	From          string
	To            string
	SourceAgent   string
	Note          string
	HumanVerified bool
	Quarantined   bool
}

type Draft struct {
	Nodes map[string]*DraftNode
	Edges []*DraftEdge
	seen  map[string]struct{}
}

func newDraft() *Draft {
	return &Draft{Nodes: map[string]*DraftNode{}, seen: map[string]struct{}{}}
}

func (d *Draft) AddNode(n DraftNode) {
	if n.NodeKey == "" {
		return
	}
	if n.Props == nil {
		n.Props = map[string]interface{}{}
	}
	if prev, ok := d.Nodes[n.NodeKey]; ok {
		if prev.HumanLocked {
			return
		}
		if prev.Maturity == "live" && n.Maturity != "live" {
			n.Maturity = "live"
		}
		if prev.Label != "" && n.Label == "" {
			n.Label = prev.Label
		}
		for k, v := range prev.Props {
			if _, exists := n.Props[k]; !exists {
				n.Props[k] = v
			}
		}
	}
	cp := n
	d.Nodes[n.NodeKey] = &cp
}

func (d *Draft) AddEdge(e DraftEdge) {
	if e.From == "" || e.To == "" || e.Relation == "" {
		return
	}
	k := e.From + "|" + e.Relation + "|" + e.To
	if _, ok := d.seen[k]; ok {
		return
	}
	d.seen[k] = struct{}{}
	d.Edges = append(d.Edges, &e)
}

func (d *Draft) HasNode(key string) bool {
	_, ok := d.Nodes[key]
	return ok
}

func propsJSON(m map[string]interface{}) string {
	if m == nil {
		return "{}"
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func parseProps(s string) map[string]interface{} {
	out := map[string]interface{}{}
	if s == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}
