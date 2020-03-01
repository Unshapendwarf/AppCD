package config

import (
	structtype "github.com/rbxorkt12/applink/pkg/application"
)

//argo Rollout 기능은 좀더 고민해봐야 할것 같다.

type Appoconfig struct {
	Orders []Order `json:"orders",yaml:"orders"`
}

type Order struct {
	Destination string  `json:"destination",yaml:"destination"`
	Charts      []Chart `json:"charts",yaml:"charts"`
}

type Chart struct {
	Repository string    `json:"repository",yaml:"repository"`
	Revision   string    `json:"revision",yaml:"revision"`
	Subpaths   []Subpath `json:"subpaths",yaml:"subpaths"`
}

type Subpath struct {
	Name                string   `json:"name",yaml:"name"`
	Path                string   `json:"path",yaml:"path"`
	Namespace           string   `json:"namespace",yaml:"namespace"`
	Chartvalues          []string `json:"chartvalues",yaml:"chartvalues"`
	Chartdeploystrategy string   `json:"chartdeploystrategy,yaml:"chartdeploystrategy"`
	Identifier 	string  `json:"identifier",yaml:"identifier"`
}


func (app *Appoconfig) ConvertApp() []*structtype.Item {
	var list []*structtype.Item
	for _, order := range app.Orders {
		for _, chart := range order.Charts {
			for _, subpath := range chart.Subpaths {
				item := &structtype.Item{}
				item.Spec.Dest.Namespace = subpath.Namespace
				item.Spec.Dest.Server = order.Destination
				item.Spec.Source.Revision = chart.Revision
				item.Spec.Source.Path = subpath.Path
				item.Spec.Source.Url = chart.Repository
				item.Spec.Project = "default"
				item.Meta.Name = subpath.Identifier + "-" + subpath.Name
				item.Spec.Source.Helm.ValueFiles = subpath.Chartvalues
				item.Meta.Annotations.Identifier = subpath.Identifier
				list = append(list, item)
			}
		}
	}
	return list

}
