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
	"strconv"
)

func ExampleRead() {
	g, err := Read([]byte(`digraph G {Hello->World}`))
	if err != nil {
		panic(err)
	}
	s := g.String()
	fmt.Println(s)
	// Output: digraph G {
	//	Hello->World;
	//	Hello;
	//	World;
	//
	//}
}

func ExampleNewGraph() {
	g := NewGraph()
	if err := g.SetName("G"); err != nil {
		panic(err)
	}
	if err := g.SetDir(true); err != nil {
		panic(err)
	}
	if err := g.AddNode("G", "Hello", nil); err != nil {
		panic(err)
	}
	if err := g.AddNode("G", "World", nil); err != nil {
		panic(err)
	}
	if err := g.AddEdge("Hello", "World", true, nil); err != nil {
		panic(err)
	}
	s := g.String()
	fmt.Println(s)
	// Output: digraph G {
	//	Hello->World;
	//	Hello;
	//	World;
	//
	//}
}

type MyOwnGraphStructure struct {
	weights map[int]map[int]int
	max     int
}

func NewMyOwnGraphStructure() *MyOwnGraphStructure {
	return &MyOwnGraphStructure{
		make(map[int]map[int]int),
		0,
	}
}

func (myown *MyOwnGraphStructure) SetStrict(strict bool) error { return nil }
func (myown *MyOwnGraphStructure) SetDir(directed bool) error  { return nil }
func (myown *MyOwnGraphStructure) SetName(name string) error   { return nil }
func (myown *MyOwnGraphStructure) AddPortEdge(src, srcPort, dst, dstPort string, directed bool, attrs map[string]string) error {
	srci, err := strconv.Atoi(src)
	if err != nil {
		return err
	}
	dsti, err := strconv.Atoi(dst)
	if err != nil {
		return err
	}
	ai, err := strconv.Atoi(attrs["label"])
	if err != nil {
		return err
	}
	if _, ok := myown.weights[srci]; !ok {
		myown.weights[srci] = make(map[int]int)
	}
	myown.weights[srci][dsti] = ai
	if srci > myown.max {
		myown.max = srci
	}
	if dsti > myown.max {
		myown.max = dsti
	}
	return nil
}
func (myown *MyOwnGraphStructure) AddEdge(src, dst string, directed bool, attrs map[string]string) error {
	return myown.AddPortEdge(src, "", dst, "", directed, attrs)
}
func (myown *MyOwnGraphStructure) AddNode(parentGraph string, name string, attrs map[string]string) error {
	return nil
}
func (myown *MyOwnGraphStructure) AddAttr(parentGraph string, field, value string) error {
	return nil
}
func (myown *MyOwnGraphStructure) AddSubGraph(parentGraph string, name string, attrs map[string]string) error {
	return nil
}
func (myown *MyOwnGraphStructure) String() string { return "" }

//An Example of how to parse into your own simpler graph structure and output it back to graphviz.
//This example reads in only numbers and outputs a matrix graph.
func ExampleMyOwnGraphStructure() {
	name := "matrix"
	parsed, err := Parse([]byte(`
		digraph G {
			1 -> 2 [ label = 5 ];
			4 -> 2 [ label = 1 ];
			4 -> 1 [ label = 2 ];
			1 -> 1 [ label = 0 ];
		}

	`))
	if err != nil {
		panic(err)
	}
	mine := NewMyOwnGraphStructure()
	if err := Analyse(parsed, mine); err != nil {
		panic(err)
	}
	output := NewGraph()
	if err := output.SetName(name); err != nil {
		panic(err)
	}
	if err := output.SetDir(true); err != nil {
		panic(err)
	}
	for i := 1; i <= mine.max; i++ {
		if err := output.AddNode(name, fmt.Sprintf("%v", i), nil); err != nil {
			panic(err)
		}
		if _, ok := mine.weights[i]; !ok {
			mine.weights[i] = make(map[int]int)
		}
	}
	for i := 1; i <= mine.max; i++ {
		for j := 1; j <= mine.max; j++ {
			if err := output.AddEdge(fmt.Sprintf("%v", i), fmt.Sprintf("%v", j), true, map[string]string{"label": fmt.Sprintf("%v", mine.weights[i][j])}); err != nil {
				panic(err)
			}
		}
	}
	s := output.String()
	fmt.Println(s)
	// Output: digraph matrix {
	//	1->1[ label=0 ];
	//	1->2[ label=5 ];
	//	1->3[ label=0 ];
	//	1->4[ label=0 ];
	//	2->1[ label=0 ];
	//	2->2[ label=0 ];
	//	2->3[ label=0 ];
	//	2->4[ label=0 ];
	//	3->1[ label=0 ];
	//	3->2[ label=0 ];
	//	3->3[ label=0 ];
	//	3->4[ label=0 ];
	//	4->1[ label=2 ];
	//	4->2[ label=1 ];
	//	4->3[ label=0 ];
	//	4->4[ label=0 ];
	//	1;
	//	2;
	//	3;
	//	4;
	//
	//}
}
