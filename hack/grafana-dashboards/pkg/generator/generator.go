// Package generator provides functionality to convert user-friendly YAML dashboard
// definitions into Grafana JSON format and GrafanaDashboard CRs.
package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/schema"
	"github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/templates"
)

// Generator converts YAML dashboard definitions to Grafana JSON.
type Generator struct {
	// TemplatesDir is the directory containing custom panel templates
	TemplatesDir string
	// OutputDir is the directory for generated files
	OutputDir string
	// Verbose enables detailed logging
	Verbose bool
}

// NewGenerator creates a new dashboard generator.
func NewGenerator(templatesDir, outputDir string, verbose bool) *Generator {
	return &Generator{
		TemplatesDir: templatesDir,
		OutputDir:    outputDir,
		Verbose:      verbose,
	}
}

// GenerateFromFile reads a YAML dashboard file and generates the output.
func (g *Generator) GenerateFromFile(inputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	dashboard := &schema.Dashboard{}
	if err := yaml.Unmarshal(data, dashboard); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return g.Generate(dashboard, inputPath)
}

// Generate converts a Dashboard definition to Grafana JSON and CR.
func (g *Generator) Generate(dashboard *schema.Dashboard, sourcePath string) error {
	// Generate Grafana JSON
	grafanaJSON, err := g.ToGrafanaJSON(dashboard)
	if err != nil {
		return fmt.Errorf("failed to generate Grafana JSON: %w", err)
	}

	// Create output directory if needed
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))

	// Write Grafana JSON
	jsonPath := filepath.Join(g.OutputDir, baseName+".json")
	jsonData, err := json.MarshalIndent(grafanaJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	// Generate GrafanaDashboard CR
	cr, err := g.ToGrafanaDashboardCR(dashboard, grafanaJSON)
	if err != nil {
		return fmt.Errorf("failed to generate CR: %w", err)
	}

	crPath := filepath.Join(g.OutputDir, baseName+"-cr.yaml")
	crData, err := yaml.Marshal(cr)
	if err != nil {
		return fmt.Errorf("failed to marshal CR: %w", err)
	}
	if err := os.WriteFile(crPath, crData, 0644); err != nil {
		return fmt.Errorf("failed to write CR file: %w", err)
	}

	if g.Verbose {
		fmt.Printf("Generated: %s, %s\n", jsonPath, crPath)
	}

	return nil
}

// GrafanaJSON represents the Grafana dashboard JSON structure.
type GrafanaJSON struct {
	ID            interface{}            `json:"id"`
	UID           string                 `json:"uid,omitempty"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	Style         string                 `json:"style"`
	Timezone      string                 `json:"timezone"`
	Editable      bool                   `json:"editable"`
	GraphTooltip  int                    `json:"graphTooltip"`
	Time          map[string]string      `json:"time"`
	Timepicker    map[string]interface{} `json:"timepicker"`
	Refresh       string                 `json:"refresh,omitempty"`
	SchemaVersion int                    `json:"schemaVersion"`
	Version       int                    `json:"version"`
	Panels        []map[string]interface{} `json:"panels"`
	Templating    map[string]interface{} `json:"templating"`
	Annotations   map[string]interface{} `json:"annotations"`
	Links         []interface{}          `json:"links"`
}

// ToGrafanaJSON converts a Dashboard definition to Grafana JSON format.
func (g *Generator) ToGrafanaJSON(dashboard *schema.Dashboard) (*GrafanaJSON, error) {
	// Initialize with defaults
	grafana := &GrafanaJSON{
		ID:            nil,
		UID:           dashboard.Metadata.Name,
		Title:         dashboard.Spec.Title,
		Description:   dashboard.Spec.Description,
		Tags:          dashboard.Spec.Tags,
		Style:         "dark",
		Timezone:      "browser",
		Editable:      dashboard.Spec.Editable,
		GraphTooltip:  0,
		SchemaVersion: 39,
		Version:       1,
		Links:         []interface{}{},
	}

	// Set time range
	if dashboard.Spec.TimeRange != nil {
		grafana.Time = map[string]string{
			"from": dashboard.Spec.TimeRange.From,
			"to":   dashboard.Spec.TimeRange.To,
		}
	} else {
		grafana.Time = map[string]string{
			"from": "now-6h",
			"to":   "now",
		}
	}

	// Set refresh
	if dashboard.Spec.RefreshInterval != "" {
		grafana.Refresh = dashboard.Spec.RefreshInterval
	}

	// Set timepicker
	grafana.Timepicker = map[string]interface{}{
		"refresh_intervals": []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"},
	}

	// Convert variables
	grafana.Templating = g.convertVariables(dashboard.Spec.Variables)

	// Convert annotations
	grafana.Annotations = g.convertAnnotations(dashboard.Spec.Annotations)

	// Convert panels
	panels, err := g.convertPanels(dashboard)
	if err != nil {
		return nil, err
	}
	grafana.Panels = panels

	return grafana, nil
}

func (g *Generator) convertVariables(variables []schema.Variable) map[string]interface{} {
	if len(variables) == 0 {
		return map[string]interface{}{"list": []interface{}{}}
	}

	varList := make([]interface{}, 0, len(variables))
	for _, v := range variables {
		varDef := map[string]interface{}{
			"name":  v.Name,
			"type":  v.Type,
			"hide":  v.Hide,
		}

		if v.Label != "" {
			varDef["label"] = v.Label
		}

		switch v.Type {
		case "datasource":
			varDef["query"] = v.Datasource
			varDef["refresh"] = 1
		case "query":
			varDef["datasource"] = map[string]interface{}{
				"type": "prometheus",
				"uid":  v.Datasource,
			}
			varDef["query"] = v.Query
			varDef["refresh"] = v.Refresh
			if varDef["refresh"] == 0 {
				varDef["refresh"] = 1
			}
			varDef["multi"] = v.Multi
			varDef["includeAll"] = v.IncludeAll
			if v.AllValue != "" {
				varDef["allValue"] = v.AllValue
			}
			if v.Regex != "" {
				varDef["regex"] = v.Regex
			}
			varDef["sort"] = v.Sort
		case "interval":
			if v.Options != nil {
				varDef["query"] = strings.Join(v.Options.Values, ",")
				varDef["auto"] = v.Options.Auto
				if v.Options.AutoCount > 0 {
					varDef["auto_count"] = v.Options.AutoCount
				}
				if v.Options.AutoMin != "" {
					varDef["auto_min"] = v.Options.AutoMin
				}
			}
			varDef["refresh"] = 2
		case "constant":
			varDef["query"] = v.Default
		case "custom":
			varDef["query"] = strings.Join(v.Values, ",")
			varDef["multi"] = v.Multi
			varDef["includeAll"] = v.IncludeAll
		case "textbox":
			if v.Default != "" {
				varDef["query"] = v.Default
			}
		}

		if v.Default != "" {
			varDef["current"] = map[string]interface{}{
				"text":  v.Default,
				"value": v.Default,
			}
		}

		varList = append(varList, varDef)
	}

	return map[string]interface{}{"list": varList}
}

func (g *Generator) convertAnnotations(annotations []schema.Annotation) map[string]interface{} {
	annotList := []interface{}{
		map[string]interface{}{
			"builtIn":    1,
			"datasource": map[string]interface{}{"type": "grafana", "uid": "-- Grafana --"},
			"enable":     true,
			"hide":       true,
			"iconColor":  "rgba(0, 211, 255, 1)",
			"name":       "Annotations & Alerts",
			"type":       "dashboard",
		},
	}

	for _, a := range annotations {
		annotDef := map[string]interface{}{
			"name":       a.Name,
			"datasource": a.Datasource,
			"enable":     a.Enable,
			"hide":       a.Hide,
		}
		if a.IconColor != "" {
			annotDef["iconColor"] = a.IconColor
		}
		if a.Query != "" {
			annotDef["expr"] = a.Query
		}
		annotList = append(annotList, annotDef)
	}

	return map[string]interface{}{"list": annotList}
}

func (g *Generator) convertPanels(dashboard *schema.Dashboard) ([]map[string]interface{}, error) {
	panels := []map[string]interface{}{}
	panelID := 1
	yPos := 0

	// Handle rows if present
	if len(dashboard.Spec.Rows) > 0 {
		for _, row := range dashboard.Spec.Rows {
			// Add row panel
			rowPanel := map[string]interface{}{
				"id":        panelID,
				"type":      "row",
				"title":     row.Title,
				"collapsed": row.Collapsed,
				"gridPos": map[string]interface{}{
					"h": 1,
					"w": 24,
					"x": 0,
					"y": yPos,
				},
			}
			if row.Repeat != "" {
				rowPanel["repeat"] = row.Repeat
			}
			panels = append(panels, rowPanel)
			panelID++
			yPos++

			// Add panels in the row
			rowPanels, newY, newID, err := g.convertPanelList(row.Panels, panelID, yPos, row.Collapsed)
			if err != nil {
				return nil, err
			}
			
			if row.Collapsed {
				rowPanel["panels"] = rowPanels
			} else {
				panels = append(panels, rowPanels...)
			}
			panelID = newID
			yPos = newY
		}
	}

	// Handle flat panels list
	if len(dashboard.Spec.Panels) > 0 {
		panelList, _, _, err := g.convertPanelList(dashboard.Spec.Panels, panelID, yPos, false)
		if err != nil {
			return nil, err
		}
		panels = append(panels, panelList...)
	}

	return panels, nil
}

func (g *Generator) convertPanelList(panelDefs []schema.Panel, startID, startY int, collapsed bool) ([]map[string]interface{}, int, int, error) {
	panels := []map[string]interface{}{}
	panelID := startID
	currentY := startY
	currentX := 0

	for _, p := range panelDefs {
		panel, err := g.convertPanel(&p, panelID, currentX, currentY)
		if err != nil {
			return nil, 0, 0, err
		}

		// Handle grid positioning
		gridPos := panel["gridPos"].(map[string]interface{})
		w := gridPos["w"].(int)
		h := gridPos["h"].(int)

		// Auto-layout if position not specified
		if currentX+w > 24 {
			currentX = 0
			currentY += h
		}
		gridPos["x"] = currentX
		gridPos["y"] = currentY

		currentX += w
		if currentX >= 24 {
			currentX = 0
			currentY += h
		}

		panels = append(panels, panel)
		panelID++
	}

	return panels, currentY, panelID, nil
}

func (g *Generator) convertPanel(p *schema.Panel, id, x, y int) (map[string]interface{}, error) {
	panel := map[string]interface{}{
		"id":    id,
		"type":  p.Type,
		"title": p.Title,
	}

	if p.Description != "" {
		panel["description"] = p.Description
	}

	if p.Transparent {
		panel["transparent"] = true
	}

	// Apply template if specified
	var tmpl *templates.PanelTemplate
	if p.Template != "" {
		tmpl = templates.GetTemplate(p.Template)
		if tmpl != nil {
			panel["type"] = tmpl.Type
		}
	}

	// Set grid position
	gridPos := map[string]interface{}{
		"x": x,
		"y": y,
		"w": 12, // default width
		"h": 8,  // default height
	}
	if p.GridPos.W > 0 {
		gridPos["w"] = p.GridPos.W
	}
	if p.GridPos.H > 0 {
		gridPos["h"] = p.GridPos.H
	}
	if p.GridPos.X > 0 {
		gridPos["x"] = p.GridPos.X
	}
	if p.GridPos.Y > 0 {
		gridPos["y"] = p.GridPos.Y
	}
	panel["gridPos"] = gridPos

	// Set datasource
	if p.Datasource != "" {
		if p.Datasource[0] == '$' {
			panel["datasource"] = map[string]interface{}{
				"type": "prometheus",
				"uid":  p.Datasource,
			}
		} else {
			panel["datasource"] = map[string]interface{}{
				"type": "prometheus",
				"uid":  p.Datasource,
			}
		}
	}

	// Convert queries
	if len(p.Queries) > 0 {
		targets := make([]map[string]interface{}, 0, len(p.Queries))
		for i, q := range p.Queries {
			target := map[string]interface{}{}
			if q.RefID != "" {
				target["refId"] = q.RefID
			} else {
				target["refId"] = string(rune('A' + i))
			}
			if q.Expr != "" {
				target["expr"] = q.Expr
			}
			if q.LegendFormat != "" {
				target["legendFormat"] = q.LegendFormat
			}
			if q.Instant {
				target["instant"] = true
			}
			if q.Range {
				target["range"] = true
			}
			if q.Interval != "" {
				target["interval"] = q.Interval
			}
			if q.Format != "" {
				target["format"] = q.Format
			}
			if q.Hide {
				target["hide"] = true
			}
			if q.Datasource != "" {
				target["datasource"] = map[string]interface{}{
					"type": "prometheus",
					"uid":  q.Datasource,
				}
			}
			if q.RawSQL != "" {
				target["rawSql"] = q.RawSQL
			}
			targets = append(targets, target)
		}
		panel["targets"] = targets
	}

	// Convert transformations
	if len(p.Transformations) > 0 {
		transforms := make([]map[string]interface{}, 0, len(p.Transformations))
		for _, t := range p.Transformations {
			transform := map[string]interface{}{
				"id":      t.ID,
				"options": t.Options,
			}
			transforms = append(transforms, transform)
		}
		panel["transformations"] = transforms
	}

	// Build options
	options := map[string]interface{}{}
	
	// Apply template options first
	if tmpl != nil && tmpl.Options != nil {
		g.applyOptions(options, tmpl.Options)
	}
	
	// Then apply panel-specific options (overrides template)
	if p.Options != nil {
		g.applyOptions(options, p.Options)
	}
	
	if len(options) > 0 {
		panel["options"] = options
	}

	// Build field config
	fieldConfig := map[string]interface{}{
		"defaults":  map[string]interface{}{},
		"overrides": []interface{}{},
	}

	// Apply template field config first
	if tmpl != nil && tmpl.FieldConfig != nil {
		g.applyFieldConfig(fieldConfig, tmpl.FieldConfig)
	}

	// Then apply panel-specific field config
	if p.FieldConfig != nil {
		g.applyFieldConfig(fieldConfig, p.FieldConfig)
	}

	// Apply thresholds from panel level
	if p.Thresholds != nil {
		defaults := fieldConfig["defaults"].(map[string]interface{})
		defaults["thresholds"] = g.convertThresholds(p.Thresholds)
	}

	panel["fieldConfig"] = fieldConfig

	// Handle repeat
	if p.Repeat != "" {
		panel["repeat"] = p.Repeat
		if p.RepeatDirection != "" {
			panel["repeatDirection"] = p.RepeatDirection
		}
		if p.MaxPerRow > 0 {
			panel["maxPerRow"] = p.MaxPerRow
		}
	}

	// Handle links
	if len(p.Links) > 0 {
		links := make([]map[string]interface{}, 0, len(p.Links))
		for _, l := range p.Links {
			link := map[string]interface{}{
				"title": l.Title,
				"url":   l.URL,
			}
			if l.TargetBlank {
				link["targetBlank"] = true
			}
			links = append(links, link)
		}
		panel["links"] = links
	}

	return panel, nil
}

func (g *Generator) applyOptions(options map[string]interface{}, opts *schema.PanelOptions) {
	if opts.Legend != nil {
		legend := map[string]interface{}{}
		if opts.Legend.DisplayMode != "" {
			legend["displayMode"] = opts.Legend.DisplayMode
		}
		if opts.Legend.Placement != "" {
			legend["placement"] = opts.Legend.Placement
		}
		legend["showLegend"] = opts.Legend.ShowLegend
		if len(opts.Legend.Calcs) > 0 {
			legend["calcs"] = opts.Legend.Calcs
		}
		options["legend"] = legend
	}

	if opts.Tooltip != nil {
		tooltip := map[string]interface{}{}
		if opts.Tooltip.Mode != "" {
			tooltip["mode"] = opts.Tooltip.Mode
		}
		if opts.Tooltip.Sort != "" {
			tooltip["sort"] = opts.Tooltip.Sort
		}
		options["tooltip"] = tooltip
	}

	if opts.Content != "" {
		options["content"] = opts.Content
	}
	if opts.Mode != "" {
		options["mode"] = opts.Mode
	}
	if opts.ShowThresholdLabels {
		options["showThresholdLabels"] = true
	}
	if opts.ShowThresholdMarkers {
		options["showThresholdMarkers"] = true
	}
	if opts.Orientation != "" {
		options["orientation"] = opts.Orientation
	}
	if opts.TextMode != "" {
		options["textMode"] = opts.TextMode
	}
	if opts.ColorMode != "" {
		options["colorMode"] = opts.ColorMode
	}
	if opts.GraphMode != "" {
		options["graphMode"] = opts.GraphMode
	}
	if opts.JustifyMode != "" {
		options["justifyMode"] = opts.JustifyMode
	}
	if opts.ShowHeader {
		options["showHeader"] = true
	}
	if opts.PieType != "" {
		options["pieType"] = opts.PieType
	}

	if opts.ReduceOptions != nil {
		reduceOpts := map[string]interface{}{}
		reduceOpts["values"] = opts.ReduceOptions.Values
		if len(opts.ReduceOptions.Calcs) > 0 {
			reduceOpts["calcs"] = opts.ReduceOptions.Calcs
		}
		if opts.ReduceOptions.Fields != "" {
			reduceOpts["fields"] = opts.ReduceOptions.Fields
		}
		if opts.ReduceOptions.Limit > 0 {
			reduceOpts["limit"] = opts.ReduceOptions.Limit
		}
		options["reduceOptions"] = reduceOpts
	}
}

func (g *Generator) applyFieldConfig(fieldConfig map[string]interface{}, fc *schema.FieldConfig) {
	if fc.Defaults != nil {
		defaults := fieldConfig["defaults"].(map[string]interface{})
		d := fc.Defaults

		if d.Unit != "" {
			defaults["unit"] = d.Unit
		}
		if d.Decimals != nil {
			defaults["decimals"] = *d.Decimals
		}
		if d.Min != nil {
			defaults["min"] = *d.Min
		}
		if d.Max != nil {
			defaults["max"] = *d.Max
		}
		if d.NoValue != "" {
			defaults["noValue"] = d.NoValue
		}
		if d.DisplayName != "" {
			defaults["displayName"] = d.DisplayName
		}

		if d.Color != nil {
			color := map[string]interface{}{}
			if d.Color.Mode != "" {
				color["mode"] = d.Color.Mode
			}
			if d.Color.FixedColor != "" {
				color["fixedColor"] = d.Color.FixedColor
			}
			defaults["color"] = color
		}

		if d.Thresholds != nil {
			defaults["thresholds"] = g.convertThresholds(d.Thresholds)
		}

		if d.Custom != nil {
			custom := map[string]interface{}{}
			c := d.Custom

			if c.DrawStyle != "" {
				custom["drawStyle"] = c.DrawStyle
			}
			if c.LineInterpolation != "" {
				custom["lineInterpolation"] = c.LineInterpolation
			}
			if c.BarAlignment != 0 {
				custom["barAlignment"] = c.BarAlignment
			}
			if c.LineWidth > 0 {
				custom["lineWidth"] = c.LineWidth
			}
			if c.FillOpacity > 0 {
				custom["fillOpacity"] = c.FillOpacity
			}
			if c.GradientMode != "" {
				custom["gradientMode"] = c.GradientMode
			}
			custom["spanNulls"] = c.SpanNulls
			if c.PointSize > 0 {
				custom["pointSize"] = c.PointSize
			}
			if c.ShowPoints != "" {
				custom["showPoints"] = c.ShowPoints
			}
			if c.AxisPlacement != "" {
				custom["axisPlacement"] = c.AxisPlacement
			}
			if c.AxisLabel != "" {
				custom["axisLabel"] = c.AxisLabel
			}
			if c.AxisSoftMin != nil {
				custom["axisSoftMin"] = *c.AxisSoftMin
			}
			if c.AxisSoftMax != nil {
				custom["axisSoftMax"] = *c.AxisSoftMax
			}

			if c.Stacking != nil {
				stacking := map[string]interface{}{}
				if c.Stacking.Mode != "" {
					stacking["mode"] = c.Stacking.Mode
				}
				if c.Stacking.Group != "" {
					stacking["group"] = c.Stacking.Group
				}
				custom["stacking"] = stacking
			}

			if c.ScaleDistribution != nil {
				scale := map[string]interface{}{}
				if c.ScaleDistribution.Type != "" {
					scale["type"] = c.ScaleDistribution.Type
				}
				if c.ScaleDistribution.Log > 0 {
					scale["log"] = c.ScaleDistribution.Log
				}
				custom["scaleDistribution"] = scale
			}

			if c.ThresholdsStyle != nil && c.ThresholdsStyle.Mode != "" {
				custom["thresholdsStyle"] = map[string]interface{}{
					"mode": c.ThresholdsStyle.Mode,
				}
			}

			if c.HideFrom != nil {
				hideFrom := map[string]interface{}{}
				hideFrom["tooltip"] = c.HideFrom.Tooltip
				hideFrom["viz"] = c.HideFrom.Viz
				hideFrom["legend"] = c.HideFrom.Legend
				custom["hideFrom"] = hideFrom
			}

			defaults["custom"] = custom
		}

		if len(d.Links) > 0 {
			links := make([]map[string]interface{}, 0, len(d.Links))
			for _, l := range d.Links {
				link := map[string]interface{}{
					"title": l.Title,
					"url":   l.URL,
				}
				if l.TargetBlank {
					link["targetBlank"] = true
				}
				links = append(links, link)
			}
			defaults["links"] = links
		}

		if len(d.Mappings) > 0 {
			mappings := make([]map[string]interface{}, 0, len(d.Mappings))
			for _, m := range d.Mappings {
				mapping := map[string]interface{}{
					"type":    m.Type,
					"options": m.Options,
				}
				mappings = append(mappings, mapping)
			}
			defaults["mappings"] = mappings
		}
	}

	if len(fc.Overrides) > 0 {
		overrides := make([]map[string]interface{}, 0, len(fc.Overrides))
		for _, o := range fc.Overrides {
			override := map[string]interface{}{
				"matcher": map[string]interface{}{
					"id":      o.Matcher.ID,
					"options": o.Matcher.Options,
				},
			}
			props := make([]map[string]interface{}, 0, len(o.Properties))
			for _, p := range o.Properties {
				props = append(props, map[string]interface{}{
					"id":    p.ID,
					"value": p.Value,
				})
			}
			override["properties"] = props
			overrides = append(overrides, override)
		}
		fieldConfig["overrides"] = overrides
	}
}

func (g *Generator) convertThresholds(t *schema.Thresholds) map[string]interface{} {
	result := map[string]interface{}{
		"mode": "absolute",
	}
	if t.Mode != "" {
		result["mode"] = t.Mode
	}

	steps := make([]map[string]interface{}, 0, len(t.Steps))
	for _, s := range t.Steps {
		step := map[string]interface{}{
			"color": s.Color,
		}
		if s.Value != nil {
			step["value"] = *s.Value
		} else {
			step["value"] = nil
		}
		steps = append(steps, step)
	}
	result["steps"] = steps

	return result
}

// GrafanaDashboardCR represents the Grafana Operator CRD structure.
type GrafanaDashboardCR struct {
	APIVersion string                   `yaml:"apiVersion"`
	Kind       string                   `yaml:"kind"`
	Metadata   map[string]interface{}   `yaml:"metadata"`
	Spec       map[string]interface{}   `yaml:"spec"`
}

// ToGrafanaDashboardCR generates a GrafanaDashboard CR from the dashboard definition.
func (g *Generator) ToGrafanaDashboardCR(dashboard *schema.Dashboard, grafanaJSON *GrafanaJSON) (*GrafanaDashboardCR, error) {
	// Marshal JSON to string
	jsonData, err := json.MarshalIndent(grafanaJSON, "", "  ")
	if err != nil {
		return nil, err
	}

	cr := &GrafanaDashboardCR{
		APIVersion: "grafana.integreatly.org/v1beta1",
		Kind:       "GrafanaDashboard",
		Metadata: map[string]interface{}{
			"name": dashboard.Metadata.Name,
		},
		Spec: map[string]interface{}{
			"json": string(jsonData),
		},
	}

	// Set namespace if provided
	if dashboard.Metadata.Namespace != "" {
		cr.Metadata["namespace"] = dashboard.Metadata.Namespace
	}

	// Set labels if provided
	if len(dashboard.Metadata.Labels) > 0 {
		cr.Metadata["labels"] = dashboard.Metadata.Labels
	}

	// Set annotations if provided
	if len(dashboard.Metadata.Annotations) > 0 {
		cr.Metadata["annotations"] = dashboard.Metadata.Annotations
	}

	// Set instance selector
	if dashboard.Spec.InstanceSelector != nil {
		selector := map[string]interface{}{}
		if len(dashboard.Spec.InstanceSelector.MatchLabels) > 0 {
			selector["matchLabels"] = dashboard.Spec.InstanceSelector.MatchLabels
		}
		if len(dashboard.Spec.InstanceSelector.MatchExpressions) > 0 {
			exprs := make([]map[string]interface{}, 0, len(dashboard.Spec.InstanceSelector.MatchExpressions))
			for _, e := range dashboard.Spec.InstanceSelector.MatchExpressions {
				expr := map[string]interface{}{
					"key":      e.Key,
					"operator": e.Operator,
				}
				if len(e.Values) > 0 {
					expr["values"] = e.Values
				}
				exprs = append(exprs, expr)
			}
			selector["matchExpressions"] = exprs
		}
		cr.Spec["instanceSelector"] = selector
	}

	// Set folder reference
	if dashboard.Spec.FolderRef != "" {
		cr.Spec["folderRef"] = dashboard.Spec.FolderRef
	} else if dashboard.Spec.Folder != "" {
		cr.Spec["folder"] = dashboard.Spec.Folder
	}

	// Set resync period
	if dashboard.Spec.ResyncPeriod != "" {
		cr.Spec["resyncPeriod"] = dashboard.Spec.ResyncPeriod
	}

	// Set cross-namespace import
	if dashboard.Spec.AllowCrossNamespaceImport {
		cr.Spec["allowCrossNamespaceImport"] = true
	}

	return cr, nil
}
