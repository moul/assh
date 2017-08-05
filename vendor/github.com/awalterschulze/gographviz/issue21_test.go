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

func TestIssue21Subgraph(t *testing.T) {
	inputString := `
	digraph G {
    Ga->Gb;
    sA->sB;
    ssA->ssB;
    
     subgraph clusterone {
        fillcolor=red;
        style=filled;
        sA;
        sB;
        
        subgraph clustertwo {
            fillcolor=blue;
            style=filled;
            ssA;
        	ssB;
       }
    }
    
    Ga;
    Gb;
}
`
	parsedGraph, err := Read([]byte(inputString))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("parsedGraph: %s", parsedGraph.String())

	_, c1ok := parsedGraph.Relations.ParentToChildren["G"]["clusterone"]
	_, c2ok := parsedGraph.Relations.ParentToChildren["clusterone"]["clustertwo"]
	if !c1ok || !c2ok {
		t.Fatalf("parsed: expected parent to child relation G-(%v)>clusterone-(%v)>clustertwo", c1ok, c2ok)
	}

	_, c1ok = parsedGraph.Relations.ChildToParents["clusterone"]["G"]
	_, c2ok = parsedGraph.Relations.ChildToParents["clustertwo"]["clusterone"]
	if !c1ok || !c2ok {
		t.Fatalf("parsed: expected child to parent relation G-(%v)>clusterone-(%v)>clustertwo", c1ok, c2ok)
	}

	g := NewGraph()
	if err := g.SetName("G"); err != nil {
		t.Fatal(err)
	}
	if err := g.SetDir(true); err != nil {
		t.Fatal(err)
	}

	if err := g.AddNode("G", "Ga", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("G", "Gb", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge("Ga", "Gb", true, nil); err != nil {
		t.Fatal(err)
	}

	if err := g.AddSubGraph("G", "clusterone", map[string]string{
		"style":     "filled",
		"fillcolor": "red",
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("clusterone", "sA", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("clusterone", "sB", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge("sA", "sB", true, nil); err != nil {
		t.Fatal(err)
	}

	if err := g.AddSubGraph("clusterone", "clustertwo", map[string]string{
		"style":     "filled",
		"fillcolor": "blue",
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("clustertwo", "ssA", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("clustertwo", "ssB", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge("ssA", "ssB", true, nil); err != nil {
		t.Fatal(err)
	}

	t.Logf("apiGraph: %s", g.String())

	_, c1ok = g.Relations.ParentToChildren["G"]["clusterone"]
	_, c2ok = g.Relations.ParentToChildren["clusterone"]["clustertwo"]
	if !c1ok || !c2ok {
		t.Fatalf("api: expected parent to child relation G-(%v)>clusterone-(%v)>clustertwo", c1ok, c2ok)
	}

	_, c1ok = g.Relations.ChildToParents["clusterone"]["G"]
	_, c2ok = g.Relations.ChildToParents["clustertwo"]["clusterone"]
	if !c1ok || !c2ok {
		t.Fatalf("api: expected child to parent relation G-(%v)>clusterone-(%v)>clustertwo", c1ok, c2ok)
	}

}
