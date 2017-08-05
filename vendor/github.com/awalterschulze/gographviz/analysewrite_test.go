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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/awalterschulze/gographviz/internal/parser"
)

func (nodes *Nodes) String() string {
	s := "Nodes:"
	for i := range nodes.Nodes {
		s += fmt.Sprintf("Node{%v}", nodes.Nodes[i])
	}
	return s + "\n"
}

func (edges *Edges) String() string {
	s := "Edges:"
	for i := range edges.Edges {
		s += fmt.Sprintf("Edge{%v}", edges.Edges[i])
	}
	return s + "\n"
}

func anal(t *testing.T, input string) Interface {
	t.Logf("Input: %v\n", input)
	g, err := parser.ParseString(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Parsed: %v\n", g)
	ag := NewGraph()
	if err := Analyse(g, ag); err != nil {
		t.Fatal(err)
	}
	t.Logf("Analysed: %v\n", ag)
	agstr := ag.String()
	t.Logf("Written: %v\n", agstr)
	g2, err := parser.ParseString(agstr)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Parsed %v\n", g2)
	ag2 := NewEscape()
	if err := Analyse(g2, ag2); err != nil {
		t.Fatal(err)
	}
	t.Logf("Analysed %v\n", ag2)
	ag2str := ag2.String()
	t.Logf("Written: %v\n", ag2str)
	if agstr != ag2str {
		t.Fatalf("analysed: want %s got %s", agstr, ag2str)
	}
	return ag2
}

func analfile(t *testing.T, filename string) Interface {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	all, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return anal(t, string(all))
}

func analtest(t *testing.T, testname string) Interface {
	return analfile(t, "./testdata/"+testname)
}

func TestHelloWorldString(t *testing.T) {
	input := `digraph G {Hello->World}`
	anal(t, input)
}

func TestHelloWorldFile(t *testing.T) {
	analfile(t, "./testdata/helloworld.gv.txt")
}

func TestAttr(t *testing.T) {
	anal(t,
		"digraph finite_state { rankdir = LR }")
}

func TestString(t *testing.T) {
	anal(t,
		`digraph finite_state { rankdir = "LR" }`)
}

func TestAttrList(t *testing.T) {
	anal(t, `
digraph { node [ shape = doublecircle ] }`)
}

func TestStringLit(t *testing.T) {
	anal(t, `digraph finite_state_machine {
	size= "8" ; }`)
}

func TestHashComments(t *testing.T) {
	anal(t, `## bla \n
  digraph G {Hello->World}`)
}

func TestIntLit(t *testing.T) {
	anal(t, `graph G {
	1 -- 30 [dim=1];}`)
}

func TestFloat1(t *testing.T) {
	anal(t, `digraph { Damping = 2.0 }`)
}

func TestFloat2(t *testing.T) {
	anal(t, `digraph { Damping = .1 }`)
}

func TestNegative(t *testing.T) {
	anal(t, `digraph { -2 -> -1 }`)
}

func TestUnderscore(t *testing.T) {
	anal(t, `digraph { dim = 1 }`)
}

func TestNonAscii(t *testing.T) {
	anal(t, `digraph {	label=Tóth }`)
}

func TestPorts(t *testing.T) {
	anal(t, `digraph { "node6":f0 -> "node9":f1 }`)
}

func TestHtml(t *testing.T) {
	anal(t, `digraph { tooltip = <<table></table>> }`)
}

func TestIdWithKeyword(t *testing.T) {
	anal(t, `digraph { edgeURL = "a" }`)
}

func TestSubGraph(t *testing.T) {
	anal(t, `digraph { subgraph { a -> b } }`)
}

func TestImplicitSubGraph(t *testing.T) {
	anal(t, `digraph { { a -> b } }`)
}

func TestEdges(t *testing.T) {
	anal(t, `digraph { a0 -> a1 -> a2 -> a3 }`)
}

func TestEasyFsm1(t *testing.T) {
	anal(t, `digraph finite_state_machine {
	rankdir=LR;
	size="8,5";
	node [shape = circle];
	LR_0 -> LR_2 [ label = "SS(B)" ];
	LR_0 -> LR_1 [ label = "SS(S)" ];
	LR_1 -> LR_3 [ label = "S($end)" ];
	LR_2 -> LR_6 [ label = "SS(b)" ];
	LR_2 -> LR_5 [ label = "SS(a)" ];
	LR_2 -> LR_4 [ label = "S(A)" ];
	LR_5 -> LR_7 [ label = "S(b)" ];
	LR_5 -> LR_5 [ label = "S(a)" ];
	LR_6 -> LR_6 [ label = "S(b)" ];
	LR_6 -> LR_5 [ label = "S(a)" ];
	LR_7 -> LR_8 [ label = "S(b)" ];
	LR_7 -> LR_5 [ label = "S(a)" ];
	LR_8 -> LR_6 [ label = "S(b)" ];
	LR_8 -> LR_5 [ label = "S(a)" ];
}`)
}

//node [shape = doublecircle]; LR_0 LR_3 LR_4 LR_8; should be applied to the nodes
func TestEasyFsm2(t *testing.T) {
	anal(t, `digraph finite_state_machine {
	rankdir=LR;
	size="8,5";
	node [shape = doublecircle]; LR_0 LR_3 LR_4 LR_8;
	node [shape = circle];
	LR_0 -> LR_2 [ label = "SS(B)" ];
	LR_0 -> LR_1 [ label = "SS(S)" ];
	LR_1 -> LR_3 [ label = "S($end)" ];
	LR_2 -> LR_6 [ label = "SS(b)" ];
	LR_2 -> LR_5 [ label = "SS(a)" ];
	LR_2 -> LR_4 [ label = "S(A)" ];
	LR_5 -> LR_7 [ label = "S(b)" ];
	LR_5 -> LR_5 [ label = "S(a)" ];
	LR_6 -> LR_6 [ label = "S(b)" ];
	LR_6 -> LR_5 [ label = "S(a)" ];
	LR_7 -> LR_8 [ label = "S(b)" ];
	LR_7 -> LR_5 [ label = "S(a)" ];
	LR_8 -> LR_6 [ label = "S(b)" ];
	LR_8 -> LR_5 [ label = "S(a)" ];
}`)
}

func TestSubSubGraph(t *testing.T) {
	anal(t, `
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
`)
}

func TestEmptyAttrList(t *testing.T) {
	anal(t, `digraph g { edge [ ] }`)
}

func TestHelloWorld(t *testing.T) {
	analtest(t, "helloworld.gv.txt")
}

func TestCluster(t *testing.T) {
	analtest(t, "cluster.gv.txt")
}

func TestPsg(t *testing.T) {
	analtest(t, "psg.gv.txt")
}

func TestTransparency(t *testing.T) {
	analtest(t, "transparency.gv.txt")
}

func TestCrazy(t *testing.T) {
	analtest(t, "crazy.gv.txt")
}

func TestKennedyanc(t *testing.T) {
	analtest(t, "kennedyanc.gv.txt")
}

func TestRoot(t *testing.T) {
	analtest(t, "root.gv.txt")
}

func TestTwpoi(t *testing.T) {
	analtest(t, "twopi.gv.txt")
}

func TestDataStruct(t *testing.T) {
	analtest(t, "datastruct.gv.txt")
}

func TestLionShare(t *testing.T) {
	analtest(t, "lion_share.gv.txt")
}

func TestSdh(t *testing.T) {
	analtest(t, "sdh.gv.txt")
}

func TestUnix(t *testing.T) {
	analtest(t, "unix.gv.txt")
}

func TestEr(t *testing.T) {
	analtest(t, "er.gv.txt")
}

func TestNerworkMapTwopi(t *testing.T) {
	analtest(t, "networkmap_twopi.gv.txt")
}

func TestSibling(t *testing.T) {
	analtest(t, "siblings.gv.txt")
}

func TestWorld(t *testing.T) {
	analtest(t, "world.gv.txt")
}

func TestFdpclust(t *testing.T) {
	analtest(t, "fdpclust.gv.txt")
}

func TestPhilo(t *testing.T) {
	analtest(t, "philo.gv.txt")
}

func TestSoftmaint(t *testing.T) {
	analtest(t, "softmaint.gv.txt")
}

func TestFsm(t *testing.T) {
	analtest(t, "fsm.gv.txt")
}

func TestProcess(t *testing.T) {
	analtest(t, "process.gv.txt")
}

func TestSwitchGv(t *testing.T) {
	analtest(t, "switch.gv.txt")
}

func TestGd19942007(t *testing.T) {
	analtest(t, "gd_1994_2007.gv.txt")
}

func TestProfile(t *testing.T) {
	analtest(t, "profile.gv.txt")
}

func TestTrafficLights(t *testing.T) {
	analtest(t, "traffic_lights.gv.txt")
}
