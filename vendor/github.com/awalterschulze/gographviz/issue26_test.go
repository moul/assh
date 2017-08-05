//Copyright 2017 GoGraphviz Authors
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package gographviz

import (
	"testing"
)

// https://github.com/awalterschulze/gographviz/issues/26
func TestIssue26DefaultAttrs(t *testing.T) {
	inputString := `
	digraph G {
		node [shape=record];
		edge [style=dashed];
		subgraph cluster_key {
			graph [label=KEY];
			user;
			node [style=dashed];
			edge [color=blue];
			dep;
			user -> dep;
		}
		A;
		node [color=red];
		B;
		A -> B;
		node [shape=diamond];
		edge [arrowhead=open];
		C;
		A -> C;
		edge [style=bold];
		B -> C;
	}
	`

	g, err := Read([]byte(inputString))
	if err != nil {
		t.Fatal(err)
	}

	_, ok := g.Relations.ParentToChildren["cluster_key"]["user"]
	if !ok {
		t.Fatal(`expected node "user" in "cluster_key"`)
	}
	_, ok = g.Relations.ParentToChildren["cluster_key"]["dep"]
	if !ok {
		t.Fatal(`expected node "dep" in "cluster_key"`)
	}

	type nodeCase struct {
		name  string
		attr  Attr
		value string
	}
	for _, c := range []nodeCase{
		// Top-level defaults apply within subgraph.
		{"user", "shape", "record"},
		{"dep", "shape", "record"},
		// Default in subgraph only applies to later nodes in that subgraph.
		{"user", "style", ""},
		{"dep", "style", "dashed"},
		{"A", "style", ""},
		{"B", "style", ""},
		{"C", "style", ""},
		// Default at top level only applies to later nodes.
		{"A", "color", ""},
		{"B", "color", "red"},
		{"C", "color", "red"},
		{"C", "shape", "diamond"},
		// New default attribute combines with previous default.
		{"A", "shape", "record"},
		{"B", "shape", "record"},
	} {
		node := g.Nodes.Lookup[c.name]
		if value := node.Attrs[c.attr]; value != c.value {
			t.Errorf("expected %s=%s in %s; got %s", c.attr, c.value, c.name, value)
		}
	}

	type edgeCase struct {
		src, dst string
		attr     Attr
		value    string
	}
	for _, c := range []edgeCase{
		// Top-level defaults apply within subgraph.
		{"user", "dep", "style", "dashed"},
		// Default in subgraph only applies to later edges in that subgraph.
		{"user", "dep", "color", "blue"},
		{"A", "B", "color", ""},
		{"A", "C", "color", ""},
		{"B", "C", "color", ""},
		// Default at top level only applies to later edges, and
		// new default attribute combines with previous default.
		{"A", "B", "style", "dashed"},
		{"A", "C", "style", "dashed"},
		{"B", "C", "style", "bold"},
		{"A", "B", "arrowhead", ""},
		{"user", "dep", "arrowhead", ""},
		{"A", "C", "arrowhead", "open"},
		{"B", "C", "arrowhead", "open"},
	} {
		edges := g.Edges.SrcToDsts[c.src][c.dst]
		if len(edges) == 0 {
			t.Fatalf("expected edge %s->%s", c.src, c.dst)
		}
		for _, e := range edges {
			if value := e.Attrs[c.attr]; value != c.value {
				t.Errorf("expected %s=%s in %s->%s; got %s", c.attr, c.value,
					c.src, c.dst, value)
			}
		}
	}
}
