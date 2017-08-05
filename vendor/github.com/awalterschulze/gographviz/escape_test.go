//Copyright 2013 GoGraphviz Authors
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
	"strings"
	"testing"
)

func TestEscape(t *testing.T) {
	g := NewEscape()
	if err := g.SetName("asdf adsf"); err != nil {
		t.Fatal(err)
	}
	if err := g.SetDir(true); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("asdf asdf", "kasdf99 99", map[string]string{
		"URL": "<a",
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("asdf asdf", "7", map[string]string{
		"URL": "<a",
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("asdf asdf", "a << b", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge("kasdf99 99", "7", true, nil); err != nil {
		t.Fatal(err)
	}
	s := g.String()
	if !strings.HasPrefix(s, `digraph "asdf adsf" {
	"kasdf99 99"->7;
	"a &lt;&lt; b";
	"kasdf99 99" [ URL="<a" ];
	7 [ URL="<a" ];

}`) {
		t.Fatalf("%s", s)
	}
	if !g.IsNode("a << b") {
		t.Fatalf("should be a node")
	}
}

func TestClusterSubgraphs(t *testing.T) {
	g := NewEscape()
	if err := g.SetName("G"); err != nil {
		t.Fatal(err)
	}
	if err := g.SetDir(false); err != nil {
		t.Fatal(err)
	}
	if err := g.AddSubGraph("G", "cluster0", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddSubGraph("cluster0", "cluster_1", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddSubGraph("cluster0", "cluster_2", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode("G", "Code deployment", nil); err != nil {
		t.Fatal(err)
	}
	if err := g.AddPortEdge("cluster_2", "", "cluster_1", "", false, nil); err != nil {
		t.Fatal(err)
	}
	s := g.String()
	graphStr := `graph G {
	cluster_2--cluster_1;
	subgraph cluster0 {
	subgraph cluster_1 {

}
;
	subgraph cluster_2 {

}
;

}
;
	"Code deployment";

}`
	if !strings.HasPrefix(s, graphStr) {
		t.Fatalf("%s", s)
	}
	g2, err := Parse([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	s2 := g2.String()
	if !strings.HasPrefix(s2, graphStr) {
		t.Fatalf("%s", s2)
	}
}
