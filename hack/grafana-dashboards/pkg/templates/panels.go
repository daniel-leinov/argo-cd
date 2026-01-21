// Package templates provides predefined panel configurations for common use cases.
package templates

import "github.com/argoproj/argo-cd/v3/hack/grafana-dashboards/pkg/schema"

// PanelTemplate represents a reusable panel configuration.
type PanelTemplate struct {
	// Name is the template identifier
	Name string
	// Type is the panel type
	Type string
	// Options are the default panel options
	Options *schema.PanelOptions
	// FieldConfig are the default field configurations
	FieldConfig *schema.FieldConfig
	// Description of what the template is for
	Description string
}

// GetTemplate returns a panel template by name.
func GetTemplate(name string) *PanelTemplate {
	templates := GetAllTemplates()
	if t, ok := templates[name]; ok {
		return t
	}
	return nil
}

// GetAllTemplates returns all available panel templates.
func GetAllTemplates() map[string]*PanelTemplate {
	return map[string]*PanelTemplate{
		"timeseries-basic":     TimeseriesBasic(),
		"timeseries-stacked":   TimeseriesStacked(),
		"timeseries-rate":      TimeseriesRate(),
		"timeseries-bytes":     TimeseriesBytes(),
		"timeseries-percent":   TimeseriesPercent(),
		"timeseries-duration":  TimeseriesDuration(),
		"gauge-percent":        GaugePercent(),
		"gauge-basic":          GaugeBasic(),
		"stat-value":           StatValue(),
		"stat-percent":         StatPercent(),
		"stat-uptime":          StatUptime(),
		"stat-count":           StatCount(),
		"table-basic":          TableBasic(),
		"table-logs":           TableLogs(),
		"barchart-horizontal":  BarChartHorizontal(),
		"barchart-vertical":    BarChartVertical(),
		"piechart-basic":       PieChartBasic(),
		"text-markdown":        TextMarkdown(),
		"logs-basic":           LogsBasic(),
		"heatmap-basic":        HeatmapBasic(),
		"alertlist-basic":      AlertListBasic(),
	}
}

// TimeseriesBasic returns a basic time series panel template.
func TimeseriesBasic() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "timeseries-basic",
		Type:        "timeseries",
		Description: "Basic time series visualization with standard defaults",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "single",
				Sort: "none",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "linear",
					LineWidth:         1,
					FillOpacity:       0,
					GradientMode:      "none",
					ShowPoints:        "auto",
					AxisPlacement:     "auto",
					SpanNulls:         false,
				},
			},
		},
	}
}

// TimeseriesStacked returns a stacked time series panel template.
func TimeseriesStacked() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "timeseries-stacked",
		Type:        "timeseries",
		Description: "Stacked area time series for cumulative metrics",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "all",
				Sort: "desc",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "linear",
					LineWidth:         1,
					FillOpacity:       30,
					GradientMode:      "none",
					ShowPoints:        "never",
					AxisPlacement:     "auto",
					Stacking: &schema.StackingConfig{
						Mode:  "normal",
						Group: "A",
					},
				},
			},
		},
	}
}

// TimeseriesRate returns a time series panel for rate metrics.
func TimeseriesRate() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "timeseries-rate",
		Type:        "timeseries",
		Description: "Time series optimized for rate/counter metrics (ops/sec, req/sec)",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
				Calcs:       []string{"mean", "max"},
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "all",
				Sort: "desc",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "ops",
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "linear",
					LineWidth:         1,
					FillOpacity:       10,
					GradientMode:      "none",
					ShowPoints:        "never",
					AxisPlacement:     "auto",
				},
			},
		},
	}
}

// TimeseriesBytes returns a time series panel for byte metrics.
func TimeseriesBytes() *PanelTemplate {
	showLegend := true
	zero := float64(0)
	return &PanelTemplate{
		Name:        "timeseries-bytes",
		Type:        "timeseries",
		Description: "Time series for byte-based metrics (memory, disk, network)",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "table",
				Placement:   "bottom",
				ShowLegend:  showLegend,
				Calcs:       []string{"mean", "max", "last"},
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "all",
				Sort: "desc",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "bytes",
				Min:  &zero,
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "smooth",
					LineWidth:         2,
					FillOpacity:       20,
					GradientMode:      "opacity",
					ShowPoints:        "never",
					AxisPlacement:     "auto",
				},
			},
		},
	}
}

// TimeseriesPercent returns a time series panel for percentage metrics.
func TimeseriesPercent() *PanelTemplate {
	showLegend := true
	zero := float64(0)
	hundred := float64(100)
	return &PanelTemplate{
		Name:        "timeseries-percent",
		Type:        "timeseries",
		Description: "Time series for percentage metrics (0-100%)",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
				Calcs:       []string{"mean", "max"},
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "all",
				Sort: "desc",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "percent",
				Min:  &zero,
				Max:  &hundred,
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "green"},
						{Value: floatPtr(80), Color: "yellow"},
						{Value: floatPtr(90), Color: "red"},
					},
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "linear",
					LineWidth:         1,
					FillOpacity:       10,
					ShowPoints:        "never",
					AxisPlacement:     "auto",
				},
			},
		},
	}
}

// TimeseriesDuration returns a time series panel for duration metrics.
func TimeseriesDuration() *PanelTemplate {
	showLegend := true
	zero := float64(0)
	return &PanelTemplate{
		Name:        "timeseries-duration",
		Type:        "timeseries",
		Description: "Time series for latency/duration metrics",
		Options: &schema.PanelOptions{
			Legend: &schema.LegendOptions{
				DisplayMode: "table",
				Placement:   "bottom",
				ShowLegend:  showLegend,
				Calcs:       []string{"mean", "max", "p99"},
			},
			Tooltip: &schema.TooltipOptions{
				Mode: "all",
				Sort: "desc",
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "s",
				Min:  &zero,
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
				Custom: &schema.CustomFieldConfig{
					DrawStyle:         "line",
					LineInterpolation: "smooth",
					LineWidth:         2,
					FillOpacity:       10,
					ShowPoints:        "never",
					AxisPlacement:     "auto",
				},
			},
		},
	}
}

// GaugePercent returns a gauge panel template for percentages.
func GaugePercent() *PanelTemplate {
	zero := float64(0)
	hundred := float64(100)
	return &PanelTemplate{
		Name:        "gauge-percent",
		Type:        "gauge",
		Description: "Gauge for percentage values with colored thresholds",
		Options: &schema.PanelOptions{
			ShowThresholdLabels:  false,
			ShowThresholdMarkers: true,
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "percent",
				Min:  &zero,
				Max:  &hundred,
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "green"},
						{Value: floatPtr(70), Color: "yellow"},
						{Value: floatPtr(85), Color: "orange"},
						{Value: floatPtr(95), Color: "red"},
					},
				},
			},
		},
	}
}

// GaugeBasic returns a basic gauge panel template.
func GaugeBasic() *PanelTemplate {
	return &PanelTemplate{
		Name:        "gauge-basic",
		Type:        "gauge",
		Description: "Basic gauge visualization",
		Options: &schema.PanelOptions{
			ShowThresholdLabels:  false,
			ShowThresholdMarkers: true,
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "green"},
					},
				},
			},
		},
	}
}

// StatValue returns a basic stat panel template.
func StatValue() *PanelTemplate {
	return &PanelTemplate{
		Name:        "stat-value",
		Type:        "stat",
		Description: "Large single value display",
		Options: &schema.PanelOptions{
			TextMode:    "auto",
			ColorMode:   "value",
			GraphMode:   "area",
			JustifyMode: "auto",
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "green"},
					},
				},
			},
		},
	}
}

// StatPercent returns a stat panel for percentage values.
func StatPercent() *PanelTemplate {
	zero := float64(0)
	hundred := float64(100)
	return &PanelTemplate{
		Name:        "stat-percent",
		Type:        "stat",
		Description: "Percentage stat with colored thresholds",
		Options: &schema.PanelOptions{
			TextMode:    "auto",
			ColorMode:   "value",
			GraphMode:   "area",
			JustifyMode: "auto",
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "percent",
				Min:  &zero,
				Max:  &hundred,
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "red"},
						{Value: floatPtr(80), Color: "yellow"},
						{Value: floatPtr(95), Color: "green"},
					},
				},
			},
		},
	}
}

// StatUptime returns a stat panel for uptime percentage.
func StatUptime() *PanelTemplate {
	decimals := 2
	zero := float64(0)
	hundred := float64(100)
	return &PanelTemplate{
		Name:        "stat-uptime",
		Type:        "stat",
		Description: "Uptime percentage with SLA thresholds",
		Options: &schema.PanelOptions{
			TextMode:    "auto",
			ColorMode:   "background",
			GraphMode:   "none",
			JustifyMode: "auto",
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"mean"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit:     "percent",
				Decimals: &decimals,
				Min:      &zero,
				Max:      &hundred,
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "red"},
						{Value: floatPtr(99), Color: "yellow"},
						{Value: floatPtr(99.9), Color: "green"},
					},
				},
			},
		},
	}
}

// StatCount returns a stat panel for count values.
func StatCount() *PanelTemplate {
	return &PanelTemplate{
		Name:        "stat-count",
		Type:        "stat",
		Description: "Counter/count value display",
		Options: &schema.PanelOptions{
			TextMode:    "value",
			ColorMode:   "value",
			GraphMode:   "area",
			JustifyMode: "auto",
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Unit: "short",
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
			},
		},
	}
}

// TableBasic returns a basic table panel template.
func TableBasic() *PanelTemplate {
	return &PanelTemplate{
		Name:        "table-basic",
		Type:        "table",
		Description: "Basic data table",
		Options: &schema.PanelOptions{
			ShowHeader: true,
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "thresholds",
				},
				Thresholds: &schema.Thresholds{
					Mode: "absolute",
					Steps: []schema.ThresholdStep{
						{Color: "green"},
					},
				},
				Custom: &schema.CustomFieldConfig{
					// Table-specific settings would go here
				},
			},
		},
	}
}

// TableLogs returns a table optimized for log data.
func TableLogs() *PanelTemplate {
	return &PanelTemplate{
		Name:        "table-logs",
		Type:        "table",
		Description: "Table optimized for log/event data",
		Options: &schema.PanelOptions{
			ShowHeader: true,
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "fixed",
				},
			},
		},
	}
}

// BarChartHorizontal returns a horizontal bar chart template.
func BarChartHorizontal() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "barchart-horizontal",
		Type:        "barchart",
		Description: "Horizontal bar chart for comparisons",
		Options: &schema.PanelOptions{
			Orientation: "horizontal",
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
			},
		},
	}
}

// BarChartVertical returns a vertical bar chart template.
func BarChartVertical() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "barchart-vertical",
		Type:        "barchart",
		Description: "Vertical bar chart for comparisons",
		Options: &schema.PanelOptions{
			Orientation: "vertical",
			Legend: &schema.LegendOptions{
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  showLegend,
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
			},
		},
	}
}

// PieChartBasic returns a basic pie chart template.
func PieChartBasic() *PanelTemplate {
	showLegend := true
	return &PanelTemplate{
		Name:        "piechart-basic",
		Type:        "piechart",
		Description: "Basic pie chart for proportional data",
		Options: &schema.PanelOptions{
			PieType: "pie",
			Legend: &schema.LegendOptions{
				DisplayMode: "table",
				Placement:   "right",
				ShowLegend:  showLegend,
				Calcs:       []string{"percent", "value"},
			},
			ReduceOptions: &schema.ReduceOptions{
				Calcs: []string{"lastNotNull"},
			},
		},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
			},
		},
	}
}

// TextMarkdown returns a text panel template for markdown content.
func TextMarkdown() *PanelTemplate {
	return &PanelTemplate{
		Name:        "text-markdown",
		Type:        "text",
		Description: "Text panel for documentation/notes",
		Options: &schema.PanelOptions{
			Mode: "markdown",
		},
	}
}

// LogsBasic returns a basic logs panel template.
func LogsBasic() *PanelTemplate {
	return &PanelTemplate{
		Name:        "logs-basic",
		Type:        "logs",
		Description: "Log viewer panel",
		Options:     &schema.PanelOptions{},
	}
}

// HeatmapBasic returns a basic heatmap panel template.
func HeatmapBasic() *PanelTemplate {
	return &PanelTemplate{
		Name:        "heatmap-basic",
		Type:        "heatmap",
		Description: "Heatmap for distribution visualization",
		Options:     &schema.PanelOptions{},
		FieldConfig: &schema.FieldConfig{
			Defaults: &schema.FieldDefaults{
				Color: &schema.ColorConfig{
					Mode: "palette-classic",
				},
			},
		},
	}
}

// AlertListBasic returns an alert list panel template.
func AlertListBasic() *PanelTemplate {
	return &PanelTemplate{
		Name:        "alertlist-basic",
		Type:        "alertlist",
		Description: "Panel showing firing alerts",
		Options:     &schema.PanelOptions{},
	}
}

// Helper function to create float pointers
func floatPtr(v float64) *float64 {
	return &v
}
