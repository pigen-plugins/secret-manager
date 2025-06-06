package main

import (
	"github.com/hashicorp/go-plugin"
	shared "github.com/pigen-dev/shared"
	"github.com/pigen-plugins/secret-manager/pkg"
	//"fmt"
)

func main() {
	// This is the entry point for the secret-manager plugin.
	// The plugin will be initialized and run here.
	// You can add any necessary initialization code here.
	// For example, you can set up logging, configuration, etc.
	
	// Initialize the plugin
	// plugin := shared.Plugin{
	// 	Label: "ARTIFACT_REGISTRY_DEMO",
	// 	Config: map[string]interface{}{
	// 		"project_id": "aidodev",
	// 		"prefix": "PIGEN",
	// 		"secrets": map[string]string{
	// 			"SECRET_TEST": "my-secret-value",
	// 			"ANOTHER_SECRET": "another-secret-value",
	// 		},
	// 	},
	// }
	
	// Create a new SecretManager instance
	// sm := &pkg.SecretManager{}
	
	// Set up the plugin
	// err := sm.SetupPlugin(plugin)
	// if err != nil {
	// 	fmt.Printf("Failed to set up plugin: %v\n", err)
	// }
	// output := sm.GetOutput(plugin)
	// if output.Error != nil {
	// 	fmt.Printf("Failed to set up plugin: %v\n", output.Error)
	// }
	// fmt.Println("output: ", output.Output)
	// err = sm.Destroy(plugin)
	// if err != nil {
	// 	fmt.Printf("Failed to set up plugin: %v\n", err)
	// }
	sm := &pkg.SecretManager{}
	pluginMap := map[string]plugin.Plugin{"pigenPlugin": &shared.PigenPlugin{Impl: sm}}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
	})
}