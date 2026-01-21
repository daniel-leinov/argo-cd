// Command grafana-dashboard-gen generates Grafana dashboards from user-friendly YAML definitions.
//
// Usage:
//
//	grafana-dashboard-gen generate -i <input.yaml> -o <output-dir>
//	grafana-dashboard-gen generate -d <dashboards-dir> -o <output-dir>
//	grafana-dashboard-gen list-templates
//	grafana-dashboard-gen validate -i <input.yaml>
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/generator"
	"github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/schema"
	"github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/templates"
)

var (
	inputFile     string
	inputDir      string
	outputDir     string
	templatesDir  string
	verbose       bool
	outputFormat  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "grafana-dashboard-gen",
		Short: "Generate Grafana dashboards from user-friendly YAML definitions",
		Long: `grafana-dashboard-gen is a tool for creating Grafana dashboards using a 
human-readable YAML format instead of complex JSON.

Features:
- User-friendly YAML DSL for dashboard definitions
- Reusable panel templates (timeseries, gauge, stat, table, etc.)
- Automatic conversion to Grafana JSON format
- Generation of GrafanaDashboard CRs for the Grafana Operator
- GitOps-friendly versioning with Kustomize integration`,
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate Grafana dashboard JSON and CRs from YAML definitions",
		Long: `Generate Grafana dashboard JSON files and GrafanaDashboard Custom Resources 
from user-friendly YAML definitions.

Examples:
  # Generate from a single file
  grafana-dashboard-gen generate -i dashboard.yaml -o output/

  # Generate from a directory of dashboards
  grafana-dashboard-gen generate -d dashboards/ -o output/

  # Generate with verbose output
  grafana-dashboard-gen generate -i dashboard.yaml -o output/ -v`,
		RunE: runGenerate,
	}

	generateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input YAML dashboard file")
	generateCmd.Flags().StringVarP(&inputDir, "directory", "d", "", "Directory containing YAML dashboard files")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "Output directory for generated files")
	generateCmd.Flags().StringVarP(&templatesDir, "templates", "t", "", "Directory containing custom panel templates")
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	listTemplatesCmd := &cobra.Command{
		Use:   "list-templates",
		Short: "List available panel templates",
		Long: `List all available built-in panel templates with their descriptions.

Templates can be referenced in dashboard definitions using the 'template' field:
  panels:
    - title: CPU Usage
      template: timeseries-percent
      queries:
        - expr: rate(cpu_usage[5m])`,
		RunE: runListTemplates,
	}

	listTemplatesCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (table, json, yaml)")

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a YAML dashboard definition",
		Long: `Validate a YAML dashboard definition without generating output.

Examples:
  grafana-dashboard-gen validate -i dashboard.yaml`,
		RunE: runValidate,
	}

	validateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input YAML dashboard file to validate")
	validateCmd.MarkFlagRequired("input")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new dashboard project",
		Long: `Initialize a new Grafana dashboard project with example files and directory structure.

This creates:
  - dashboards/        - Directory for dashboard YAML definitions
  - panels/            - Directory for custom panel templates  
  - output/            - Directory for generated files
  - kustomize/         - Kustomize configuration for GitOps
  - example-dashboard.yaml - Example dashboard to get started`,
		RunE: runInit,
	}

	initCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for project initialization")

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert existing Grafana JSON to YAML DSL",
		Long: `Convert an existing Grafana dashboard JSON file to the user-friendly YAML DSL format.

This is useful for migrating existing dashboards to the new format.

Examples:
  grafana-dashboard-gen convert -i existing-dashboard.json -o converted.yaml`,
		RunE: runConvert,
	}

	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input Grafana JSON file to convert")
	convertCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output YAML file path")
	convertCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(generateCmd, listTemplatesCmd, validateCmd, initCmd, convertCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if inputFile == "" && inputDir == "" {
		return fmt.Errorf("either --input or --directory must be specified")
	}

	gen := generator.NewGenerator(templatesDir, outputDir, verbose)

	if inputFile != "" {
		return gen.GenerateFromFile(inputFile)
	}

	// Process directory
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		inputPath := filepath.Join(inputDir, entry.Name())
		if err := gen.GenerateFromFile(inputPath); err != nil {
			return fmt.Errorf("failed to process %s: %w", entry.Name(), err)
		}

		if verbose {
			fmt.Printf("Processed: %s\n", entry.Name())
		}
	}

	fmt.Printf("Generated files written to: %s\n", outputDir)
	return nil
}

func runListTemplates(cmd *cobra.Command, args []string) error {
	allTemplates := templates.GetAllTemplates()

	// Sort template names
	names := make([]string, 0, len(allTemplates))
	for name := range allTemplates {
		names = append(names, name)
	}
	sort.Strings(names)

	switch outputFormat {
	case "json":
		output := make([]map[string]string, 0, len(names))
		for _, name := range names {
			t := allTemplates[name]
			output = append(output, map[string]string{
				"name":        t.Name,
				"type":        t.Type,
				"description": t.Description,
			})
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))

	case "yaml":
		output := make([]map[string]string, 0, len(names))
		for _, name := range names {
			t := allTemplates[name]
			output = append(output, map[string]string{
				"name":        t.Name,
				"type":        t.Type,
				"description": t.Description,
			})
		}
		data, _ := yaml.Marshal(output)
		fmt.Println(string(data))

	default:
		fmt.Println("Available Panel Templates:")
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("%-25s %-15s %s\n", "TEMPLATE", "TYPE", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 80))
		for _, name := range names {
			t := allTemplates[name]
			fmt.Printf("%-25s %-15s %s\n", t.Name, t.Type, t.Description)
		}
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("\nTotal: %d templates\n", len(names))
	}

	return nil
}

func runValidate(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	dashboard := &schema.Dashboard{}
	if err := yaml.Unmarshal(data, dashboard); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	// Validate required fields
	var errors []string

	if dashboard.Metadata.Name == "" {
		errors = append(errors, "metadata.name is required")
	}

	if dashboard.Spec.Title == "" {
		errors = append(errors, "spec.title is required")
	}

	// Validate panel templates
	allPanels := dashboard.Spec.Panels
	for _, row := range dashboard.Spec.Rows {
		allPanels = append(allPanels, row.Panels...)
	}

	for i, p := range allPanels {
		if p.Title == "" {
			errors = append(errors, fmt.Sprintf("panel[%d].title is required", i))
		}
		if p.Type == "" && p.Template == "" {
			errors = append(errors, fmt.Sprintf("panel[%d] requires either 'type' or 'template'", i))
		}
		if p.Template != "" {
			if templates.GetTemplate(p.Template) == nil {
				errors = append(errors, fmt.Sprintf("panel[%d].template '%s' is not a valid template", i, p.Template))
			}
		}
	}

	// Validate variables
	for i, v := range dashboard.Spec.Variables {
		if v.Name == "" {
			errors = append(errors, fmt.Sprintf("variable[%d].name is required", i))
		}
		if v.Type == "" {
			errors = append(errors, fmt.Sprintf("variable[%d].type is required", i))
		}
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(errors))
	}

	fmt.Printf("✓ %s is valid\n", inputFile)
	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	// Create directory structure
	dirs := []string{
		"dashboards",
		"panels",
		"output",
		"kustomize/base",
		"kustomize/overlays/dev",
		"kustomize/overlays/staging",
		"kustomize/overlays/prod",
	}

	for _, dir := range dirs {
		path := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create example dashboard
	exampleDashboard := `# Example Grafana Dashboard
# This file demonstrates the user-friendly YAML DSL for creating dashboards.
apiVersion: grafana-dashboards/v1
kind: GrafanaDashboard
metadata:
  name: example-dashboard
  namespace: monitoring
  labels:
    app: my-application
    team: platform
spec:
  title: "Example Application Dashboard"
  description: "Dashboard demonstrating the YAML DSL features"
  tags:
    - example
    - application
  folder: "Applications"
  editable: true
  refreshInterval: "30s"
  
  timeRange:
    from: "now-1h"
    to: "now"
  
  # Template variables for dynamic filtering
  variables:
    - name: datasource
      label: "Data Source"
      type: datasource
      datasource: prometheus
    
    - name: namespace
      label: "Namespace"
      type: query
      datasource: "$datasource"
      query: "label_values(up, namespace)"
      multi: true
      includeAll: true
      allValue: ".*"
      refresh: 2
    
    - name: interval
      label: "Interval"
      type: interval
      options:
        values: ["1m", "5m", "15m", "30m", "1h"]
        auto: true
        autoCount: 30
        autoMin: "10s"
  
  # Instance selector for Grafana Operator
  instanceSelector:
    matchLabels:
      dashboards: grafana
  
  resyncPeriod: "10m"
  
  rows:
    - title: "Overview"
      collapsed: false
      panels:
        - title: "Total Requests"
          template: stat-count
          gridPos:
            w: 6
            h: 4
          datasource: "$datasource"
          queries:
            - expr: 'sum(increase(http_requests_total{namespace=~"$namespace"}[$interval]))'
              legendFormat: "Requests"
        
        - title: "Error Rate"
          template: stat-percent
          gridPos:
            w: 6
            h: 4
          datasource: "$datasource"
          queries:
            - expr: 'sum(rate(http_requests_total{status=~"5.."}[$interval])) / sum(rate(http_requests_total[$interval])) * 100'
              legendFormat: "Error %"
        
        - title: "Uptime"
          template: stat-uptime
          gridPos:
            w: 6
            h: 4
          datasource: "$datasource"
          queries:
            - expr: 'avg_over_time(up{namespace=~"$namespace"}[$interval]) * 100'
              legendFormat: "Uptime"
        
        - title: "Active Pods"
          template: stat-value
          gridPos:
            w: 6
            h: 4
          datasource: "$datasource"
          queries:
            - expr: 'count(up{namespace=~"$namespace"} == 1)'
              legendFormat: "Pods"
    
    - title: "Request Metrics"
      collapsed: false
      panels:
        - title: "Request Rate"
          template: timeseries-rate
          gridPos:
            w: 12
            h: 8
          datasource: "$datasource"
          queries:
            - expr: 'sum(rate(http_requests_total{namespace=~"$namespace"}[$interval])) by (service)'
              legendFormat: "{{service}}"
        
        - title: "Request Latency"
          template: timeseries-duration
          gridPos:
            w: 12
            h: 8
          datasource: "$datasource"
          queries:
            - refId: A
              expr: 'histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket{namespace=~"$namespace"}[$interval])) by (le))'
              legendFormat: "p50"
            - refId: B
              expr: 'histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{namespace=~"$namespace"}[$interval])) by (le))'
              legendFormat: "p95"
            - refId: C
              expr: 'histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{namespace=~"$namespace"}[$interval])) by (le))'
              legendFormat: "p99"
    
    - title: "Resource Usage"
      collapsed: true
      panels:
        - title: "CPU Usage"
          template: timeseries-percent
          gridPos:
            w: 12
            h: 8
          datasource: "$datasource"
          queries:
            - expr: 'sum(rate(container_cpu_usage_seconds_total{namespace=~"$namespace"}[$interval])) by (pod) * 100'
              legendFormat: "{{pod}}"
        
        - title: "Memory Usage"
          template: timeseries-bytes
          gridPos:
            w: 12
            h: 8
          datasource: "$datasource"
          queries:
            - expr: 'sum(container_memory_usage_bytes{namespace=~"$namespace"}) by (pod)'
              legendFormat: "{{pod}}"
`

	if err := os.WriteFile(filepath.Join(outputDir, "dashboards", "example-dashboard.yaml"), []byte(exampleDashboard), 0644); err != nil {
		return fmt.Errorf("failed to create example dashboard: %w", err)
	}

	// Create base kustomization
	baseKustomization := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources: []
# Add generated GrafanaDashboard CRs here after running:
#   grafana-dashboard-gen generate -d ../dashboards -o ../output
# Then add:
#   - ../output/example-dashboard-cr.yaml

commonLabels:
  managed-by: grafana-dashboard-gen

namespace: monitoring
`

	if err := os.WriteFile(filepath.Join(outputDir, "kustomize/base/kustomization.yaml"), []byte(baseKustomization), 0644); err != nil {
		return fmt.Errorf("failed to create base kustomization: %w", err)
	}

	// Create dev overlay
	devOverlay := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

namespace: monitoring-dev

patches:
  - patch: |-
      - op: replace
        path: /spec/instanceSelector/matchLabels
        value:
          dashboards: grafana-dev
    target:
      kind: GrafanaDashboard
`

	if err := os.WriteFile(filepath.Join(outputDir, "kustomize/overlays/dev/kustomization.yaml"), []byte(devOverlay), 0644); err != nil {
		return fmt.Errorf("failed to create dev overlay: %w", err)
	}

	// Create staging overlay
	stagingOverlay := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

namespace: monitoring-staging

patches:
  - patch: |-
      - op: replace
        path: /spec/instanceSelector/matchLabels
        value:
          dashboards: grafana-staging
    target:
      kind: GrafanaDashboard
`

	if err := os.WriteFile(filepath.Join(outputDir, "kustomize/overlays/staging/kustomization.yaml"), []byte(stagingOverlay), 0644); err != nil {
		return fmt.Errorf("failed to create staging overlay: %w", err)
	}

	// Create prod overlay
	prodOverlay := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

namespace: monitoring-prod

patches:
  - patch: |-
      - op: replace
        path: /spec/instanceSelector/matchLabels
        value:
          dashboards: grafana-prod
    target:
      kind: GrafanaDashboard
`

	if err := os.WriteFile(filepath.Join(outputDir, "kustomize/overlays/prod/kustomization.yaml"), []byte(prodOverlay), 0644); err != nil {
		return fmt.Errorf("failed to create prod overlay: %w", err)
	}

	fmt.Printf("Project initialized in %s\n", outputDir)
	fmt.Println("\nCreated:")
	fmt.Println("  dashboards/example-dashboard.yaml  - Example dashboard definition")
	fmt.Println("  kustomize/base/kustomization.yaml  - Base Kustomize configuration")
	fmt.Println("  kustomize/overlays/*/              - Environment overlays")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit dashboards/example-dashboard.yaml or create new dashboards")
	fmt.Println("  2. Generate: grafana-dashboard-gen generate -d dashboards/ -o output/")
	fmt.Println("  3. Apply with Kustomize: kubectl apply -k kustomize/overlays/dev")
	return nil
}

func runConvert(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var grafanaJSON map[string]interface{}
	if err := json.Unmarshal(data, &grafanaJSON); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Convert to YAML DSL
	dashboard := convertJSONToDSL(grafanaJSON)

	yamlData, err := yaml.Marshal(dashboard)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write output
	outputPath := outputDir
	if outputPath == "" {
		baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputPath = baseName + ".yaml"
	}

	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Converted dashboard written to: %s\n", outputPath)
	return nil
}

func convertJSONToDSL(grafanaJSON map[string]interface{}) *schema.Dashboard {
	dashboard := &schema.Dashboard{
		APIVersion: "grafana-dashboards/v1",
		Kind:       "GrafanaDashboard",
		Metadata: schema.Metadata{
			Name: slugify(getString(grafanaJSON, "title")),
		},
		Spec: schema.DashboardSpec{
			Title:       getString(grafanaJSON, "title"),
			Description: getString(grafanaJSON, "description"),
			Editable:    getBool(grafanaJSON, "editable"),
		},
	}

	// Convert tags
	if tags, ok := grafanaJSON["tags"].([]interface{}); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				dashboard.Spec.Tags = append(dashboard.Spec.Tags, s)
			}
		}
	}

	// Convert time range
	if timeMap, ok := grafanaJSON["time"].(map[string]interface{}); ok {
		dashboard.Spec.TimeRange = &schema.TimeRange{
			From: getString(timeMap, "from"),
			To:   getString(timeMap, "to"),
		}
	}

	// Convert refresh
	if refresh := getString(grafanaJSON, "refresh"); refresh != "" {
		dashboard.Spec.RefreshInterval = refresh
	}

	// Convert variables (simplified)
	if templating, ok := grafanaJSON["templating"].(map[string]interface{}); ok {
		if list, ok := templating["list"].([]interface{}); ok {
			for _, v := range list {
				if varMap, ok := v.(map[string]interface{}); ok {
					variable := schema.Variable{
						Name:  getString(varMap, "name"),
						Label: getString(varMap, "label"),
						Type:  getString(varMap, "type"),
						Query: getString(varMap, "query"),
						Multi: getBool(varMap, "multi"),
						Hide:  getInt(varMap, "hide"),
					}
					dashboard.Spec.Variables = append(dashboard.Spec.Variables, variable)
				}
			}
		}
	}

	// Convert panels (simplified - handles flat panel list)
	if panels, ok := grafanaJSON["panels"].([]interface{}); ok {
		for _, p := range panels {
			if panelMap, ok := p.(map[string]interface{}); ok {
				panel := convertPanelFromJSON(panelMap)
				if panel != nil {
					if getString(panelMap, "type") == "row" {
						// Handle row
						row := schema.Row{
							Title:     panel.Title,
							Collapsed: getBool(panelMap, "collapsed"),
						}
						// Convert nested panels if collapsed
						if nestedPanels, ok := panelMap["panels"].([]interface{}); ok {
							for _, np := range nestedPanels {
								if npMap, ok := np.(map[string]interface{}); ok {
									if nestedPanel := convertPanelFromJSON(npMap); nestedPanel != nil {
										row.Panels = append(row.Panels, *nestedPanel)
									}
								}
							}
						}
						dashboard.Spec.Rows = append(dashboard.Spec.Rows, row)
					} else {
						dashboard.Spec.Panels = append(dashboard.Spec.Panels, *panel)
					}
				}
			}
		}
	}

	return dashboard
}

func convertPanelFromJSON(panelMap map[string]interface{}) *schema.Panel {
	panel := &schema.Panel{
		Title: getString(panelMap, "title"),
		Type:  getString(panelMap, "type"),
	}

	if panel.Type == "row" {
		return panel
	}

	if desc := getString(panelMap, "description"); desc != "" {
		panel.Description = desc
	}

	// Convert grid position
	if gridPos, ok := panelMap["gridPos"].(map[string]interface{}); ok {
		panel.GridPos = schema.GridPos{
			X: getInt(gridPos, "x"),
			Y: getInt(gridPos, "y"),
			W: getInt(gridPos, "w"),
			H: getInt(gridPos, "h"),
		}
	}

	// Convert targets/queries
	if targets, ok := panelMap["targets"].([]interface{}); ok {
		for _, t := range targets {
			if targetMap, ok := t.(map[string]interface{}); ok {
				query := schema.Query{
					RefID:        getString(targetMap, "refId"),
					Expr:         getString(targetMap, "expr"),
					LegendFormat: getString(targetMap, "legendFormat"),
					Instant:      getBool(targetMap, "instant"),
					Range:        getBool(targetMap, "range"),
					Interval:     getString(targetMap, "interval"),
				}
				panel.Queries = append(panel.Queries, query)
			}
		}
	}

	return panel
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	if v, ok := m[key].(int); ok {
		return v
	}
	return 0
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}
