package ontology

func FilterView(view string, nodes []OntNode, edges []OntEdge) ([]OntNode, []OntEdge) {
	if view == "" || view == "all" {
		return nodes, edges
	}
	if view == "bridge" {
		var outE []OntEdge
		need := map[string]bool{}
		for _, e := range edges {
			if e.Deleted || e.Layer != "bridge" {
				continue
			}
			outE = append(outE, e)
			need[e.FromKey] = true
			need[e.ToKey] = true
		}
		var outN []OntNode
		for _, n := range nodes {
			if need[n.NodeKey] {
				outN = append(outN, n)
			}
		}
		if outN == nil {
			outN = []OntNode{}
		}
		if outE == nil {
			outE = []OntEdge{}
		}
		return outN, outE
	}
	keep := map[string]bool{}
	var outN []OntNode
	for _, n := range nodes {
		if viewMatch(view, n) {
			keep[n.NodeKey] = true
			outN = append(outN, n)
		}
	}
	var outE []OntEdge
	for _, e := range edges {
		if e.Deleted {
			continue
		}
		switch view {
		case "quarantine":
			if e.Quarantined {
				outE = append(outE, e)
			}
		case "gap":
			if keep[e.FromKey] || keep[e.ToKey] {
				outE = append(outE, e)
			}
		default:
			if keep[e.FromKey] && keep[e.ToKey] {
				outE = append(outE, e)
			}
		}
	}
	if outN == nil {
		outN = []OntNode{}
	}
	if outE == nil {
		outE = []OntEdge{}
	}
	return outN, outE
}

func viewMatch(view string, n OntNode) bool {
	switch view {
	case "physical":
		return n.Layer == "physical"
	case "business":
		return n.Layer == "business" && !n.Gap
	case "system":
		return n.Layer == "system" && !n.Gap
	case "gap":
		return n.Gap || n.NodeType == "Gap"
	case "quarantine":
		return n.Quarantined
	case "bridge":
		return n.Layer == "bridge"
	default:
		return true
	}
}
