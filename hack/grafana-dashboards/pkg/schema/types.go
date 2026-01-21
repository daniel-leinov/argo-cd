// Package schema defines the YAML DSL types for user-friendly Grafana dashboard creation.
// This schema provides a human-readable alternative to raw Grafana JSON.
package schema

// Dashboard represents a complete Grafana dashboard definition in YAML format.
type Dashboard struct {
	// APIVersion indicates the schema version (e.g., "grafana-dashboards/v1")
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	// Kind should always be "GrafanaDashboard"
	Kind string `yaml:"kind" json:"kind"`
	// Metadata contains dashboard identification information
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	// Spec contains the dashboard specification
	Spec DashboardSpec `yaml:"spec" json:"spec"`
}

// Metadata contains identification and labeling information for the dashboard.
type Metadata struct {
	// Name is the unique identifier for the dashboard
	Name string `yaml:"name" json:"name"`
	// Namespace is the Kubernetes namespace for the GrafanaDashboard CR
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	// Labels for the GrafanaDashboard CR
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	// Annotations for the GrafanaDashboard CR
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

// DashboardSpec contains the main dashboard configuration.
type DashboardSpec struct {
	// Title is the display name of the dashboard
	Title string `yaml:"title" json:"title"`
	// Description provides context about the dashboard
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Tags for organizing and searching dashboards
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	// Folder is the Grafana folder name or UID
	Folder string `yaml:"folder,omitempty" json:"folder,omitempty"`
	// FolderRef references a GrafanaFolder CR
	FolderRef string `yaml:"folderRef,omitempty" json:"folderRef,omitempty"`
	// Editable determines if the dashboard can be edited in the UI
	Editable bool `yaml:"editable,omitempty" json:"editable,omitempty"`
	// RefreshInterval is the default auto-refresh interval
	RefreshInterval string `yaml:"refreshInterval,omitempty" json:"refreshInterval,omitempty"`
	// TimeRange is the default time range for the dashboard
	TimeRange *TimeRange `yaml:"timeRange,omitempty" json:"timeRange,omitempty"`
	// Variables are dashboard-level template variables
	Variables []Variable `yaml:"variables,omitempty" json:"variables,omitempty"`
	// Annotations are dashboard annotations configuration
	Annotations []Annotation `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	// Rows organizes panels into logical groups (optional)
	Rows []Row `yaml:"rows,omitempty" json:"rows,omitempty"`
	// Panels is a flat list of panels (when not using rows)
	Panels []Panel `yaml:"panels,omitempty" json:"panels,omitempty"`
	// InstanceSelector selects which Grafana instance(s) to apply to
	InstanceSelector *InstanceSelector `yaml:"instanceSelector,omitempty" json:"instanceSelector,omitempty"`
	// ResyncPeriod for Grafana Operator
	ResyncPeriod string `yaml:"resyncPeriod,omitempty" json:"resyncPeriod,omitempty"`
	// AllowCrossNamespaceImport allows importing from other namespaces
	AllowCrossNamespaceImport bool `yaml:"allowCrossNamespaceImport,omitempty" json:"allowCrossNamespaceImport,omitempty"`
}

// TimeRange defines the default time window for the dashboard.
type TimeRange struct {
	// From is the start time (e.g., "now-1h", "now-7d")
	From string `yaml:"from" json:"from"`
	// To is the end time (e.g., "now")
	To string `yaml:"to" json:"to"`
}

// Variable defines a template variable for the dashboard.
type Variable struct {
	// Name is the variable identifier (used as $name in queries)
	Name string `yaml:"name" json:"name"`
	// Label is the display name in the UI
	Label string `yaml:"label,omitempty" json:"label,omitempty"`
	// Type is the variable type (query, datasource, interval, constant, custom, textbox)
	Type string `yaml:"type" json:"type"`
	// Datasource for query-type variables
	Datasource string `yaml:"datasource,omitempty" json:"datasource,omitempty"`
	// Query for query-type variables
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
	// Values for custom-type variables
	Values []string `yaml:"values,omitempty" json:"values,omitempty"`
	// Default value
	Default string `yaml:"default,omitempty" json:"default,omitempty"`
	// Multi allows selecting multiple values
	Multi bool `yaml:"multi,omitempty" json:"multi,omitempty"`
	// IncludeAll adds an "All" option
	IncludeAll bool `yaml:"includeAll,omitempty" json:"includeAll,omitempty"`
	// AllValue is the value used when "All" is selected
	AllValue string `yaml:"allValue,omitempty" json:"allValue,omitempty"`
	// Regex to filter results
	Regex string `yaml:"regex,omitempty" json:"regex,omitempty"`
	// Hide level (0=show, 1=hide label, 2=hide completely)
	Hide int `yaml:"hide,omitempty" json:"hide,omitempty"`
	// Refresh behavior (0=never, 1=on dashboard load, 2=on time range change)
	Refresh int `yaml:"refresh,omitempty" json:"refresh,omitempty"`
	// Sort order (0=disabled, 1=asc, 2=desc, 3=asc numerical, 4=desc numerical)
	Sort int `yaml:"sort,omitempty" json:"sort,omitempty"`
	// Options for interval-type variables
	Options *IntervalOptions `yaml:"options,omitempty" json:"options,omitempty"`
}

// IntervalOptions contains options for interval-type variables.
type IntervalOptions struct {
	// Values is a list of interval options (e.g., ["1m", "5m", "15m", "1h"])
	Values []string `yaml:"values" json:"values"`
	// Auto enables automatic interval calculation
	Auto bool `yaml:"auto,omitempty" json:"auto,omitempty"`
	// AutoCount is the target number of points for auto intervals
	AutoCount int `yaml:"autoCount,omitempty" json:"autoCount,omitempty"`
	// AutoMin is the minimum auto interval
	AutoMin string `yaml:"autoMin,omitempty" json:"autoMin,omitempty"`
}

// Annotation defines a dashboard annotation configuration.
type Annotation struct {
	// Name is the annotation name
	Name string `yaml:"name" json:"name"`
	// Datasource for the annotation
	Datasource string `yaml:"datasource" json:"datasource"`
	// Enable determines if the annotation is shown by default
	Enable bool `yaml:"enable,omitempty" json:"enable,omitempty"`
	// Hide determines if the annotation is hidden in the UI
	Hide bool `yaml:"hide,omitempty" json:"hide,omitempty"`
	// IconColor for the annotation markers
	IconColor string `yaml:"iconColor,omitempty" json:"iconColor,omitempty"`
	// Query for the annotation data
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
}

// Row is a logical grouping of panels with a collapsible header.
type Row struct {
	// Title of the row
	Title string `yaml:"title" json:"title"`
	// Collapsed determines if the row is collapsed by default
	Collapsed bool `yaml:"collapsed,omitempty" json:"collapsed,omitempty"`
	// Repeat variable name for row repetition
	Repeat string `yaml:"repeat,omitempty" json:"repeat,omitempty"`
	// Panels within the row
	Panels []Panel `yaml:"panels" json:"panels"`
}

// Panel represents a single visualization panel.
type Panel struct {
	// Title of the panel
	Title string `yaml:"title" json:"title"`
	// Description for the panel
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Type of panel (timeseries, gauge, stat, table, barchart, piechart, text, logs, etc.)
	Type string `yaml:"type" json:"type"`
	// Template references a predefined panel template
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
	// GridPos defines the panel position and size
	GridPos GridPos `yaml:"gridPos,omitempty" json:"gridPos,omitempty"`
	// Datasource for the panel queries
	Datasource string `yaml:"datasource,omitempty" json:"datasource,omitempty"`
	// Queries for the panel
	Queries []Query `yaml:"queries,omitempty" json:"queries,omitempty"`
	// Transformations to apply to the data
	Transformations []Transformation `yaml:"transformations,omitempty" json:"transformations,omitempty"`
	// Options are panel-specific options
	Options *PanelOptions `yaml:"options,omitempty" json:"options,omitempty"`
	// FieldConfig for field-level configuration
	FieldConfig *FieldConfig `yaml:"fieldConfig,omitempty" json:"fieldConfig,omitempty"`
	// Thresholds for the panel
	Thresholds *Thresholds `yaml:"thresholds,omitempty" json:"thresholds,omitempty"`
	// Links are panel-level links
	Links []Link `yaml:"links,omitempty" json:"links,omitempty"`
	// Repeat variable name for panel repetition
	Repeat string `yaml:"repeat,omitempty" json:"repeat,omitempty"`
	// RepeatDirection is horizontal or vertical
	RepeatDirection string `yaml:"repeatDirection,omitempty" json:"repeatDirection,omitempty"`
	// MaxPerRow when repeating horizontally
	MaxPerRow int `yaml:"maxPerRow,omitempty" json:"maxPerRow,omitempty"`
	// Transparent background
	Transparent bool `yaml:"transparent,omitempty" json:"transparent,omitempty"`
}

// GridPos defines the panel position and dimensions.
type GridPos struct {
	// X is the horizontal position (0-23)
	X int `yaml:"x,omitempty" json:"x,omitempty"`
	// Y is the vertical position
	Y int `yaml:"y,omitempty" json:"y,omitempty"`
	// W is the width (1-24)
	W int `yaml:"w,omitempty" json:"w,omitempty"`
	// H is the height
	H int `yaml:"h,omitempty" json:"h,omitempty"`
}

// Query defines a data query for a panel.
type Query struct {
	// RefID is the query identifier (A, B, C, etc.)
	RefID string `yaml:"refId,omitempty" json:"refId,omitempty"`
	// Expr is the query expression (for Prometheus, Loki, etc.)
	Expr string `yaml:"expr,omitempty" json:"expr,omitempty"`
	// LegendFormat is the legend template
	LegendFormat string `yaml:"legendFormat,omitempty" json:"legendFormat,omitempty"`
	// Instant is for instant queries
	Instant bool `yaml:"instant,omitempty" json:"instant,omitempty"`
	// Range is for range queries
	Range bool `yaml:"range,omitempty" json:"range,omitempty"`
	// Interval is the minimum time interval
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty"`
	// Format is the query result format
	Format string `yaml:"format,omitempty" json:"format,omitempty"`
	// Hide hides the query results
	Hide bool `yaml:"hide,omitempty" json:"hide,omitempty"`
	// Datasource overrides the panel datasource
	Datasource string `yaml:"datasource,omitempty" json:"datasource,omitempty"`
	// RawSQL for SQL datasources
	RawSQL string `yaml:"rawSql,omitempty" json:"rawSql,omitempty"`
}

// Transformation defines a data transformation.
type Transformation struct {
	// ID of the transformation (e.g., "organize", "calculateField", "filterByValue")
	ID string `yaml:"id" json:"id"`
	// Options are transformation-specific options
	Options map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// PanelOptions contains panel-specific display options.
type PanelOptions struct {
	// Legend configuration
	Legend *LegendOptions `yaml:"legend,omitempty" json:"legend,omitempty"`
	// Tooltip configuration
	Tooltip *TooltipOptions `yaml:"tooltip,omitempty" json:"tooltip,omitempty"`
	// Text content for text panels
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
	// Mode for text panels (markdown, html)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// ShowThresholdLabels for gauge panels
	ShowThresholdLabels bool `yaml:"showThresholdLabels,omitempty" json:"showThresholdLabels,omitempty"`
	// ShowThresholdMarkers for gauge panels
	ShowThresholdMarkers bool `yaml:"showThresholdMarkers,omitempty" json:"showThresholdMarkers,omitempty"`
	// Orientation (auto, horizontal, vertical)
	Orientation string `yaml:"orientation,omitempty" json:"orientation,omitempty"`
	// TextMode for stat panels (auto, value, name, value_and_name, none)
	TextMode string `yaml:"textMode,omitempty" json:"textMode,omitempty"`
	// ColorMode for stat panels (value, background, none)
	ColorMode string `yaml:"colorMode,omitempty" json:"colorMode,omitempty"`
	// GraphMode for stat panels (none, area)
	GraphMode string `yaml:"graphMode,omitempty" json:"graphMode,omitempty"`
	// JustifyMode for stat panels (auto, center)
	JustifyMode string `yaml:"justifyMode,omitempty" json:"justifyMode,omitempty"`
	// ReduceOptions for panels that reduce data to a single value
	ReduceOptions *ReduceOptions `yaml:"reduceOptions,omitempty" json:"reduceOptions,omitempty"`
	// Table-specific options
	ShowHeader bool `yaml:"showHeader,omitempty" json:"showHeader,omitempty"`
	// PieChart type (pie, donut)
	PieType string `yaml:"pieType,omitempty" json:"pieType,omitempty"`
}

// LegendOptions configures the panel legend.
type LegendOptions struct {
	// DisplayMode (list, table, hidden)
	DisplayMode string `yaml:"displayMode,omitempty" json:"displayMode,omitempty"`
	// Placement (bottom, right)
	Placement string `yaml:"placement,omitempty" json:"placement,omitempty"`
	// ShowLegend enables/disables the legend
	ShowLegend bool `yaml:"showLegend,omitempty" json:"showLegend,omitempty"`
	// Calcs are the calculations to show (min, max, mean, last, etc.)
	Calcs []string `yaml:"calcs,omitempty" json:"calcs,omitempty"`
}

// TooltipOptions configures the panel tooltip.
type TooltipOptions struct {
	// Mode (single, all, none)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Sort (none, asc, desc)
	Sort string `yaml:"sort,omitempty" json:"sort,omitempty"`
}

// ReduceOptions configures how data is reduced to a single value.
type ReduceOptions struct {
	// Values show all values instead of calculated
	Values bool `yaml:"values,omitempty" json:"values,omitempty"`
	// Calcs are the calculations (last, lastNotNull, min, max, mean, sum, count, etc.)
	Calcs []string `yaml:"calcs,omitempty" json:"calcs,omitempty"`
	// Fields pattern to match
	Fields string `yaml:"fields,omitempty" json:"fields,omitempty"`
	// Limit the number of values shown
	Limit int `yaml:"limit,omitempty" json:"limit,omitempty"`
}

// FieldConfig defines field-level configuration.
type FieldConfig struct {
	// Defaults are the default field settings
	Defaults *FieldDefaults `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	// Overrides are field-specific overrides
	Overrides []FieldOverride `yaml:"overrides,omitempty" json:"overrides,omitempty"`
}

// FieldDefaults defines default field settings.
type FieldDefaults struct {
	// Unit for the field (short, bytes, percent, s, ms, etc.)
	Unit string `yaml:"unit,omitempty" json:"unit,omitempty"`
	// Decimals to display
	Decimals *int `yaml:"decimals,omitempty" json:"decimals,omitempty"`
	// Min value
	Min *float64 `yaml:"min,omitempty" json:"min,omitempty"`
	// Max value
	Max *float64 `yaml:"max,omitempty" json:"max,omitempty"`
	// Color configuration
	Color *ColorConfig `yaml:"color,omitempty" json:"color,omitempty"`
	// NoValue text to display when no data
	NoValue string `yaml:"noValue,omitempty" json:"noValue,omitempty"`
	// DisplayName override
	DisplayName string `yaml:"displayName,omitempty" json:"displayName,omitempty"`
	// MappingsNote: simplified mapping support
	Mappings []Mapping `yaml:"mappings,omitempty" json:"mappings,omitempty"`
	// Thresholds for the field
	Thresholds *Thresholds `yaml:"thresholds,omitempty" json:"thresholds,omitempty"`
	// Custom field settings
	Custom *CustomFieldConfig `yaml:"custom,omitempty" json:"custom,omitempty"`
	// Links for data links
	Links []Link `yaml:"links,omitempty" json:"links,omitempty"`
}

// ColorConfig defines color settings.
type ColorConfig struct {
	// Mode (thresholds, palette-classic, fixed, etc.)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// FixedColor when mode is "fixed"
	FixedColor string `yaml:"fixedColor,omitempty" json:"fixedColor,omitempty"`
}

// CustomFieldConfig contains visualization-specific settings.
type CustomFieldConfig struct {
	// DrawStyle (line, bars, points)
	DrawStyle string `yaml:"drawStyle,omitempty" json:"drawStyle,omitempty"`
	// LineInterpolation (linear, smooth, stepBefore, stepAfter)
	LineInterpolation string `yaml:"lineInterpolation,omitempty" json:"lineInterpolation,omitempty"`
	// BarAlignment (-1=before, 0=center, 1=after)
	BarAlignment int `yaml:"barAlignment,omitempty" json:"barAlignment,omitempty"`
	// LineWidth for lines
	LineWidth int `yaml:"lineWidth,omitempty" json:"lineWidth,omitempty"`
	// FillOpacity (0-100)
	FillOpacity int `yaml:"fillOpacity,omitempty" json:"fillOpacity,omitempty"`
	// GradientMode (none, opacity, hue, scheme)
	GradientMode string `yaml:"gradientMode,omitempty" json:"gradientMode,omitempty"`
	// SpanNulls connects null values
	SpanNulls bool `yaml:"spanNulls,omitempty" json:"spanNulls,omitempty"`
	// PointSize for point visualization
	PointSize int `yaml:"pointSize,omitempty" json:"pointSize,omitempty"`
	// ShowPoints (auto, always, never)
	ShowPoints string `yaml:"showPoints,omitempty" json:"showPoints,omitempty"`
	// Stacking mode
	Stacking *StackingConfig `yaml:"stacking,omitempty" json:"stacking,omitempty"`
	// AxisPlacement (auto, left, right, hidden)
	AxisPlacement string `yaml:"axisPlacement,omitempty" json:"axisPlacement,omitempty"`
	// AxisLabel custom label
	AxisLabel string `yaml:"axisLabel,omitempty" json:"axisLabel,omitempty"`
	// AxisSoftMin is the soft minimum
	AxisSoftMin *float64 `yaml:"axisSoftMin,omitempty" json:"axisSoftMin,omitempty"`
	// AxisSoftMax is the soft maximum
	AxisSoftMax *float64 `yaml:"axisSoftMax,omitempty" json:"axisSoftMax,omitempty"`
	// ScaleDistribution (linear, log)
	ScaleDistribution *ScaleDistribution `yaml:"scaleDistribution,omitempty" json:"scaleDistribution,omitempty"`
	// ThresholdsStyle for displaying thresholds
	ThresholdsStyle *ThresholdsStyle `yaml:"thresholdsStyle,omitempty" json:"thresholdsStyle,omitempty"`
	// HideFrom hides from specific views
	HideFrom *HideFrom `yaml:"hideFrom,omitempty" json:"hideFrom,omitempty"`
}

// StackingConfig configures stacking behavior.
type StackingConfig struct {
	// Mode (none, normal, percent)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Group name for stacking
	Group string `yaml:"group,omitempty" json:"group,omitempty"`
}

// ScaleDistribution configures axis scale.
type ScaleDistribution struct {
	// Type (linear, log)
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	// Log base for logarithmic scale
	Log int `yaml:"log,omitempty" json:"log,omitempty"`
}

// ThresholdsStyle configures how thresholds are displayed.
type ThresholdsStyle struct {
	// Mode (off, line, area, line+area, dashed, dashed+area)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

// HideFrom configures hiding from specific views.
type HideFrom struct {
	// Tooltip hides from tooltip
	Tooltip bool `yaml:"tooltip,omitempty" json:"tooltip,omitempty"`
	// Viz hides from visualization
	Viz bool `yaml:"viz,omitempty" json:"viz,omitempty"`
	// Legend hides from legend
	Legend bool `yaml:"legend,omitempty" json:"legend,omitempty"`
}

// Mapping defines a value mapping.
type Mapping struct {
	// Type (value, range, regex, special)
	Type string `yaml:"type" json:"type"`
	// Options for the mapping
	Options map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// Thresholds defines threshold configuration.
type Thresholds struct {
	// Mode (absolute, percentage)
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Steps are the threshold steps
	Steps []ThresholdStep `yaml:"steps,omitempty" json:"steps,omitempty"`
}

// ThresholdStep defines a single threshold step.
type ThresholdStep struct {
	// Value is the threshold value (null for base)
	Value *float64 `yaml:"value,omitempty" json:"value,omitempty"`
	// Color for this threshold
	Color string `yaml:"color" json:"color"`
}

// FieldOverride defines a field-specific override.
type FieldOverride struct {
	// Matcher identifies which fields to override
	Matcher OverrideMatcher `yaml:"matcher" json:"matcher"`
	// Properties to override
	Properties []OverrideProperty `yaml:"properties" json:"properties"`
}

// OverrideMatcher identifies fields for override.
type OverrideMatcher struct {
	// ID is the matcher type (byName, byRegexp, byType, byFrameRefID)
	ID string `yaml:"id" json:"id"`
	// Options are matcher-specific options
	Options interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// OverrideProperty is a single property override.
type OverrideProperty struct {
	// ID is the property identifier
	ID string `yaml:"id" json:"id"`
	// Value is the property value
	Value interface{} `yaml:"value" json:"value"`
}

// Link defines a data link or panel link.
type Link struct {
	// Title of the link
	Title string `yaml:"title" json:"title"`
	// URL is the link target
	URL string `yaml:"url" json:"url"`
	// TargetBlank opens in new tab
	TargetBlank bool `yaml:"targetBlank,omitempty" json:"targetBlank,omitempty"`
}

// InstanceSelector selects Grafana instances.
type InstanceSelector struct {
	// MatchLabels selects instances by labels
	MatchLabels map[string]string `yaml:"matchLabels,omitempty" json:"matchLabels,omitempty"`
	// MatchExpressions for advanced selection
	MatchExpressions []MatchExpression `yaml:"matchExpressions,omitempty" json:"matchExpressions,omitempty"`
}

// MatchExpression is a label selector expression.
type MatchExpression struct {
	// Key is the label key
	Key string `yaml:"key" json:"key"`
	// Operator is the comparison operator (In, NotIn, Exists, DoesNotExist)
	Operator string `yaml:"operator" json:"operator"`
	// Values for In/NotIn operators
	Values []string `yaml:"values,omitempty" json:"values,omitempty"`
}
