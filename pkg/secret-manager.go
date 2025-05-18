package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/pigen-plugins/secret-manager/helpers"
	"github.com/pigen-plugins/secret-manager/pkg/terraform"
	shared "github.com/pigen-dev/shared"
	tfengine "github.com/pigen-dev/shared/tfengine"
)


type SecretManager struct {
	Label string `yaml:"label" json:"label"`
	Config Config `yaml:"config" json:"config"`
	Output Output `yaml:"output" json:"output"`
}



type Config struct {
	ProjectId string `yaml:"project_id" json:"project_id"`
	Prefix		string `yaml:"prefix" json:"prefix"`
	Secrets map[string]string `yaml:"secrets" json:"secrets"`
}

type Output struct {
	SecretsList []string `yaml:"secrets_list" json:"secrets_list"`
}


func (s *SecretManager) Initializer(plugin shared.Plugin) (*tfengine.Terraform ,error) {
	config := Config{}
	err:= helpers.YamlConfigParser(plugin.Config, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse YAML config: %v", err)
	}
	s.Config = config
	s.Label = plugin.Label
	s.Config.Secrets = helpers.PrefixMapKeys(s.Config.Secrets, s.Config.Prefix)
	fmt.Println("Parsed config:", s)
	// Initialize Terraform
	files := terraform.LoadTFFiles()
	tfVars, err := helpers.StructToMap(s.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert struct to map: %v", err)
	}
	fmt.Println("Terraform variables:", tfVars)
	t, err := tfengine.NewTF(tfVars, files, s.Label)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup Terraform executor: %v", err)
	}
	
	return t, nil
}



func (s *SecretManager) SetupPlugin(plugin shared.Plugin) error {
	tf, err := s.Initializer(plugin)
	ctx := context.Background()
	if err != nil {
		return fmt.Errorf("Failed to initialize plugin: %v", err)
	}

	// 1. Initialize Terraform
	fmt.Println(s.Label)
	if err := tf.TerraformInit(ctx, s.Config.ProjectId, s.Label); err != nil {
		return fmt.Errorf("Error during Terraform init: %v", err)
	}

	// 2. Plan Terraform changes
	if err := tf.TerraformPlan(ctx); err != nil {
		return fmt.Errorf("Error during Terraform plan: %v", err)
	}

	
	if err := tf.TerraformApply(ctx); err != nil {
		return fmt.Errorf("Error during Terraform apply: %v", err)
	}
	log.Println("Terraform apply completed.")
	return nil
}


func (s *SecretManager) GetOutput(plugin shared.Plugin) shared.GetOutputResponse {
	_, err := s.Initializer(plugin)
	if err != nil {
		return shared.GetOutputResponse{Output: nil, Error: fmt.Errorf("Failed to initialize plugin: %v", err)}
	}
	var secretsList []string
	for k, _ := range s.Config.Secrets {
		secretsList = append(secretsList, k)
	}
	output, err := helpers.StructToMap(s.Config)
	if err != nil {
		return shared.GetOutputResponse{Output: nil, Error: fmt.Errorf("Failed to convert struct to map: %v", err)}
	}
	return shared.GetOutputResponse{Output: output, Error: nil}
}


func (s *SecretManager) Destroy(plugin shared.Plugin) error {
	tf, err := s.Initializer(plugin)
	if err != nil {
		return fmt.Errorf("Failed to initialize plugin: %v", err)
	}
	ctx := context.Background()
	// 1. Initialize Terraform
	if err := tf.TerraformInit(ctx, s.Config.ProjectId, s.Label); err != nil {
		return fmt.Errorf("Error during Terraform init: %v", err)
	}

	if err := tf.TerraformDestroy(ctx); err != nil {
		return fmt.Errorf("Error during Terraform destroy: %v", err)
	}
	log.Println("Terraform destroy completed.")
	return nil
}