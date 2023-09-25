package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	honeycombApiKey    string
	semanticModelPath  = "model"
	dryRun             = false
	parseModelsOnly    = false
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

	if parseModelsOnly {
		fmt.Println()

		keys := make([]string, 0, len(semanticAttributes))
		for k := range semanticAttributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Println(k, ":", semanticAttributes[k])
		}

		fmt.Println()
		fmt.Println(len(semanticAttributes), "semantic attributes found")

	} else {
		err = updateHoneycombDatasets()
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
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
						time.Sleep(200 * time.Millisecond)
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

		if group.Type == "metric" {
			// only capture the metric name for metric groups
			// this will ignore locally defined metric attributes, which will not be namespaced
			description := strings.TrimSpace(group.Brief)
			if description != "" {
				// upserts the description into the semanticAttributes map
				semanticAttributes[group.MetricName] = strings.TrimSpace(description)
			}

		} else {
			// get the group prefix
			prefix := group.Prefix
			if prefix == "" {
				// if no prefix is specified
				// skip this group as it is likely a metrics attribute_group without a namespace
				continue
			}

			// loop through all attributes in the group
			for _, attribute := range group.Attributes {
				description := strings.TrimSpace(attribute.Brief)
				if description != "" {
					// upserts the description into the semanticAttributes map
					semanticAttributes[prefix+"."+attribute.ID] = strings.TrimSpace(description)
				}
			}
		}

	}
	return nil
}

func validateOptions() error {
	flag.StringVar(&honeycombApiKey, "honeycomb-api-key", LookupEnvOrString("HONEYCOMB_API_KEY", honeycombApiKey), "Honeycomb API Key")
	flag.StringVar(&semanticModelPath, "model-path", LookupEnvOrString("SEMANTIC_MODEL_PATH", semanticModelPath), "Path for OpenTelemetry semantic models")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run Mode")
	flag.BoolVar(&parseModelsOnly, "parse-models-only", false, "Parse Semantic Models only")
	flag.Parse()

	if honeycombApiKey == "" && parseModelsOnly == false {
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
