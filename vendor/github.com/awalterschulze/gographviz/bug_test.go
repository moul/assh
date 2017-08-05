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

	"github.com/awalterschulze/gographviz/ast"
	"github.com/awalterschulze/gographviz/internal/parser"
)

type bugSubGraphWorldVisitor struct {
	t     *testing.T
	found bool
}

func (w *bugSubGraphWorldVisitor) Visit(v ast.Elem) ast.Visitor {
	edge, ok := v.(ast.EdgeStmt)
	if !ok {
		return w
	}
	if edge.Source.GetID().String() != "2" {
		return w
	}
	dst := edge.EdgeRHS[0].Destination
	if _, ok := dst.(*ast.SubGraph); !ok {
		w.t.Fatalf("2 -> Not SubGraph")
	} else {
		w.found = true
	}
	return w
}

func TestBugSubGraphWorld(t *testing.T) {
	g := analtest(t, "world.gv.txt")
	st, err := parser.ParseString(g.String())
	if err != nil {
		t.Fatal(err)
	}
	s := &bugSubGraphWorldVisitor{
		t: t,
	}
	st.Walk(s)
	if !s.found {
		t.Fatalf("2 -> SubGraph not found")
	}
}
