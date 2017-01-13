package configviz

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
	"github.com/moul/advanced-ssh-config/pkg/config"
)

func nodename(input string) string {
	return fmt.Sprintf(`"%s"`, input)
}

type GraphSettings struct {
	ShowIsolatedHosts bool
}

func Graph(cfg *config.Config, settings *GraphSettings) (string, error) {
	graph := gographviz.NewGraph()
	graph.SetName("G")
	graph.SetDir(true)

	hostsToShow := map[string]bool{}

	for _, host := range cfg.Hosts {
		if len(host.Gateways) == 0 && !settings.ShowIsolatedHosts {
			continue
		}

		hostsToShow[nodename(host.Name())] = true
		idx := 0
		for _, gateway := range host.Gateways {
			if gateway == "direct" {
				continue
			}
			if _, found := cfg.Hosts[gateway]; !found {
				continue
			}
			idx++
			hostsToShow[nodename(gateway)] = true
			graph.AddEdge(nodename(host.Name()), nodename(gateway), true, map[string]string{"color": "red", "label": fmt.Sprintf("%d", idx)})
		}
	}

	for hostname := range hostsToShow {
		graph.AddNode("G", hostname, map[string]string{"color": "blue"})
	}

	return graph.String(), nil
}
