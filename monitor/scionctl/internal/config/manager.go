package config

import (
	"os"
	"scionctl/internal/api"

	"go.yaml.in/yaml/v3"
)

// Define config type for our node list
// We will parse from the config.yaml file
type NodeConfig struct {
	Hosts []struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
		ISD     int16  `yaml:"isd"`
		AS      int16  `yaml:"as"`
	} `yaml:"hosts"`
	DefaultPort          int16  `yaml:"default_port"`
	DefaultSciondAddress string `yaml:"default_sciond_address"`
}

type NodeManager struct {
	nodes          map[string]*api.ScionNode
	port           int16
	sciond_address string
}

// To be used globally
var CmdNodeManager *NodeManager

func InitializeManager(configPath string) error {
	nodeConfig, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	CmdNodeManager, err = NewNodeManager(nodeConfig)
	return err
}

func NewNodeManager(config *NodeConfig) (*NodeManager, error) {
	nm := &NodeManager{
		nodes: make(map[string]*api.ScionNode),
		// Share default port and sciond address across all nodes
		port:           config.DefaultPort,
		sciond_address: config.DefaultSciondAddress,
	}

	// Populate nodes map
	for _, host := range config.Hosts {
		nm.nodes[host.Name] = &api.ScionNode{
			Addr:       host.Address,
			Name:       host.Name,
			Port:       nm.port,
			ISD:        host.ISD,
			AS:         host.AS,
			ScionDAddr: nm.sciond_address,
		}
	}
	return nm, nil
}

func (nm *NodeManager) GetNode(name string) (*api.ScionNode, bool) {
	node, exists := nm.nodes[name]

	if !exists {
		return nil, false
	}
	// Return a copy to avoid modification
	return &api.ScionNode{
		Addr:       node.Addr,
		Name:       node.Name,
		Port:       node.Port,
		ISD:        node.ISD,
		AS:         node.AS,
		ScionDAddr: node.ScionDAddr,
	}, exists
}

func loadConfig(path string) (*NodeConfig, error) {
	var config NodeConfig

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
