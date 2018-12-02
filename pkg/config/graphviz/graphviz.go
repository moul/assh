package configviz // import "moul.io/assh/pkg/config/graphviz"

import (
	"fmt"

	"github.com/awalterschulze/gographviz"

	"moul.io/assh/pkg/config"
)

func nodename(input string) string {
	return fmt.Sprintf(`"%s"`, input)
}

// GraphSettings are used to change the Graph() function behavior.
type GraphSettings struct {
	ShowIsolatedHosts bool
	NoResolveWildcard bool
	NoInherits        bool
}

// Graph computes and returns a dot-compatible graph representation of the config.
func Graph(cfg *config.Config, settings *GraphSettings) (string, error) {
	graph := gographviz.NewGraph()
	if err := graph.SetName("G"); err != nil {
		return "", err
	}
	if err := graph.SetDir(true); err != nil {
		return "", err
	}

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
				if settings.NoResolveWildcard {
					continue
				}
				gw := cfg.GetGatewaySafe(gateway)
				if gw == nil {
					continue
				}
				if err := graph.AddEdge(nodename(host.Name()), nodename(gw.RawName()), true, map[string]string{"color": "red", "label": nodename(gateway)}); err != nil {
					return "", err
				}
				hostsToShow[nodename(gw.RawName())] = true
				continue
			}
			idx++
			hostsToShow[nodename(gateway)] = true
			if err := graph.AddEdge(nodename(host.Name()), nodename(gateway), true, map[string]string{"color": "red", "label": fmt.Sprintf("%d", idx)}); err != nil {
				return "", err
			}
		}

		if !settings.NoInherits {
			for _, inherit := range host.Inherits {
				hostsToShow[nodename(inherit)] = true
				if err := graph.AddEdge(nodename(host.Name()), nodename(inherit), true, map[string]string{"color": "black", "style": "dashed"}); err != nil {
					return "", err
				}
			}
		}
	}

	for hostname := range hostsToShow {
		if err := graph.AddNode("G", hostname, map[string]string{"color": "blue"}); err != nil {
			return "", err
		}
	}

	return graph.String(), nil
}
