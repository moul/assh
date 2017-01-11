package configviz

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
	"github.com/moul/advanced-ssh-config/pkg/config"
)

func nodename(input string) string {
	return fmt.Sprintf(`"%s"`, input)
}

func Graph(cfg *config.Config) (string, error) {
	graph := gographviz.NewGraph()
	graph.SetName("G")
	graph.SetDir(true)

	for _, host := range cfg.Hosts {
		graph.AddNode("G", nodename(host.Name()), map[string]string{"color": "blue"})
		idx := 0
		for _, gateway := range host.Gateways {
			if gateway == "direct" {
				continue
			}
			if _, found := cfg.Hosts[gateway]; !found {
				continue
			}
			idx++
			graph.AddEdge(nodename(host.Name()), nodename(gateway), true, map[string]string{"color": "red", "label": fmt.Sprintf("%d", idx)})
		}
	}

	return graph.String(), nil
}
