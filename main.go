package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	honeycombApiKey    string
	semanticModelPath  = "model"
	dryRun             = false
	semanticAttributes map[string]string
)

func main() {
	err := validateOptions()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Starting Honeycomb OpenTelemetry Semantic Model Updater...")
	if dryRun {
		fmt.Println("Running in dry run mode")
	}

	semanticAttributes = make(map[string]string)
	err = parseSemanticModels()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Found", len(semanticAttributes), "semantic attributes")

	err = updateHoneycombDatasets()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	fmt.Println("Done!")
}

func updateHoneycombDatasets() error {
	fmt.Println("Updating Honeycomb datasets...")

	hnyClient := NewHoneycombClient(honeycombApiKey)
	datasets, err := hnyClient.ListAllDatasets()
	if err != nil {
		fmt.Println("Error while listing all datasets")
		return err
	}

	fmt.Println("Found", len(datasets), "datasets")

	columnsUpdated := 0

	for _, dataset := range datasets {
		columns, err := hnyClient.ListAllColumns(dataset)
		if err != nil {
			fmt.Println("Error while listing all columns for dataset:", dataset.Name)
			return err
		}

		for _, column := range columns {
			if column.Description == "" {
				if description, ok := semanticAttributes[column.KeyName]; ok {
					fmt.Println("Updating column:", column.KeyName, "in dataset:", dataset.Name)
					column.Description = truncate(description, 255)
					if !dryRun {
						err = hnyClient.UpdateColumn(dataset, column)
						if err != nil {
							fmt.Println("Error while updating column:", column.KeyName)
							return err
						}
					}
					columnsUpdated++
				}
			}
		}
	}

	if dryRun {
		fmt.Println("Dry run mode enabled. Would have updated", columnsUpdated, "dataset columns")
	} else {
		fmt.Println("Updated", columnsUpdated, "dataset columns")
	}

	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for i := range s {
		if n == 0 {
			return s[:i]
		}
		n--
	}
	return s
}
func parseSemanticModels() error {
	fmt.Println("Parsing OpenTelemetry semantic models...")

	err := filepath.WalkDir(semanticModelPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println("ERROR Parsing:", path)
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".yaml") {
			fmt.Println("Parsing: ", path)
			err = parseModel(path)
			if err != nil {
				fmt.Println("ERROR Parsing:", path)
				return err
			}
		}
		return nil
	})

	return err
}

func parseModel(path string) error {

	yamlFile, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	var model OtelSemanticModel
	err = yaml.Unmarshal(yamlFile, &model)

	if err != nil {
		return err
	}

	for _, group := range model.Groups {
		prefix := group.Prefix

		for _, attribute := range group.Attributes {
			description := strings.TrimSpace(attribute.Brief)
			if description != "" {
				semanticAttributes[prefix+"."+attribute.ID] = strings.TrimSpace(description)
			}
		}
	}
	return nil
}

func validateOptions() error {
	flag.StringVar(&honeycombApiKey, "honeycomb-api-key", LookupEnvOrString("HONEYCOMB_API_KEY", honeycombApiKey), "Honeycomb API Key")
	flag.StringVar(&semanticModelPath, "model-path", LookupEnvOrString("SEMANTIC_MODEL_PATH", semanticModelPath), "Path for OpenTelemetry semantic models")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run Mode")
	flag.Parse()

	if honeycombApiKey == "" {
		printUsage()
		return fmt.Errorf("missing: Honeycomb API Key")
	}
	return nil
}

func printUsage() {
	fmt.Println("Usage: hny-otel-semantic [options]")
	flag.PrintDefaults()
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
