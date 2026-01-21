# Grafana Dashboard Generator

A user-friendly solution for creating, versioning, restoring, and managing Grafana dashboards using the [Grafana Operator](https://github.com/grafana/grafana-operator).

## Problem Statement

Managing Grafana dashboards as JSON files is challenging:
- Dashboard JSON files are typically 1000+ lines, making them hard to read and maintain
- Manual JSON editing is error-prone
- Comparing dashboard versions in Git is difficult
- Reusing common panel configurations requires copy-paste

## Solution

This tool provides:
- **User-friendly YAML DSL**: Define dashboards in a clean, readable YAML format
- **Reusable Panel Templates**: Pre-built templates for common visualizations
- **Automatic JSON Generation**: Convert YAML to Grafana-compatible JSON
- **Grafana Operator Integration**: Generate GrafanaDashboard CRs for Kubernetes
- **GitOps-Ready**: Full Kustomize support for environment-based deployments
- **Dashboard Versioning**: Git-friendly format with easy diff and merge

## Quick Start

### 1. Initialize a Project

```bash
go run ./cmd init -o my-dashboards
cd my-dashboards
```

This creates:
```
my-dashboards/
├── dashboards/
│   └── example-dashboard.yaml    # Example dashboard definition
├── panels/                        # Custom panel templates (optional)
├── output/                        # Generated files
└── kustomize/
    ├── base/
    │   └── kustomization.yaml
    └── overlays/
        ├── dev/
        ├── staging/
        └── prod/
```

### 2. Create a Dashboard

Create a new file `dashboards/my-app.yaml`:

```yaml
apiVersion: grafana-dashboards/v1
kind: GrafanaDashboard
metadata:
  name: my-application
  namespace: monitoring
spec:
  title: "My Application Dashboard"
  description: "Monitors my application health and performance"
  tags:
    - application
    - monitoring
  
  variables:
    - name: datasource
      type: datasource
      datasource: prometheus
    
    - name: namespace
      type: query
      datasource: "$datasource"
      query: "label_values(up, namespace)"
      includeAll: true
  
  instanceSelector:
    matchLabels:
      dashboards: grafana
  
  rows:
    - title: "Overview"
      panels:
        - title: "Request Rate"
          template: timeseries-rate  # Use pre-built template
          gridPos: { w: 12, h: 8 }
          datasource: "$datasource"
          queries:
            - expr: 'sum(rate(http_requests_total{namespace=~"$namespace"}[5m]))'
              legendFormat: "Requests/s"
        
        - title: "Error Rate"
          template: stat-percent
          gridPos: { w: 6, h: 4 }
          datasource: "$datasource"
          queries:
            - expr: 'sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) * 100'
```

### 3. Generate Output Files

```bash
go run ./cmd generate -d dashboards/ -o output/
```

This creates:
- `output/my-application.json` - Grafana JSON dashboard
- `output/my-application-cr.yaml` - GrafanaDashboard Custom Resource

### 4. Deploy with GitOps

```bash
# Add generated CRs to kustomize base
# Then deploy to your environment
kubectl apply -k kustomize/overlays/dev
```

## YAML DSL Reference

### Dashboard Structure

```yaml
apiVersion: grafana-dashboards/v1
kind: GrafanaDashboard
metadata:
  name: dashboard-name           # Unique identifier (required)
  namespace: monitoring          # Kubernetes namespace
  labels: {}                     # Additional labels
  annotations: {}                # Additional annotations
spec:
  title: "Dashboard Title"       # Display name (required)
  description: "Description"     # Dashboard description
  tags: [tag1, tag2]            # Tags for organization
  folder: "FolderName"          # Grafana folder
  folderRef: "folder-cr-name"   # Or reference a GrafanaFolder CR
  editable: true                 # Allow UI editing
  refreshInterval: "30s"         # Auto-refresh interval
  
  timeRange:
    from: "now-6h"
    to: "now"
  
  variables: []                  # Template variables
  annotations: []                # Dashboard annotations
  rows: []                       # Panel rows (optional grouping)
  panels: []                     # Flat panel list (alternative to rows)
  
  # Grafana Operator settings
  instanceSelector:
    matchLabels:
      dashboards: grafana
  resyncPeriod: "10m"
  allowCrossNamespaceImport: false
```

### Variables

```yaml
variables:
  # Datasource selector
  - name: datasource
    label: "Data Source"
    type: datasource
    datasource: prometheus
  
  # Query-based variable
  - name: namespace
    label: "Namespace"
    type: query
    datasource: "$datasource"
    query: "label_values(up, namespace)"
    multi: true              # Allow multiple selection
    includeAll: true         # Add "All" option
    allValue: ".*"           # Value when "All" selected
    refresh: 2               # 1=on load, 2=on time change
    sort: 1                  # 1=asc, 2=desc
    regex: ""                # Filter results
  
  # Interval variable
  - name: interval
    label: "Interval"
    type: interval
    options:
      values: ["1m", "5m", "15m", "1h"]
      auto: true
      autoCount: 30
      autoMin: "10s"
  
  # Custom variable
  - name: environment
    type: custom
    values: ["dev", "staging", "prod"]
    multi: false
  
  # Constant
  - name: app_name
    type: constant
    default: "my-app"
```

### Panels

```yaml
panels:
  - title: "Panel Title"           # Required
    description: "Description"      # Optional
    type: timeseries               # Panel type (or use template)
    template: timeseries-rate      # Use pre-built template
    
    gridPos:                        # Position and size
      x: 0                          # Column (0-23)
      y: 0                          # Row
      w: 12                         # Width (1-24)
      h: 8                          # Height
    
    datasource: "$datasource"       # Data source reference
    
    queries:                        # Data queries
      - refId: A                    # Query identifier
        expr: 'prometheus_query'    # PromQL expression
        legendFormat: "{{label}}"   # Legend template
        instant: false              # Instant query
        range: true                 # Range query
        interval: "1m"              # Min interval
    
    # Panel-specific options
    options:
      legend:
        displayMode: table          # list, table, hidden
        placement: bottom           # bottom, right
        showLegend: true
        calcs: [mean, max, last]
      tooltip:
        mode: all                   # single, all, none
        sort: desc
    
    # Field configuration
    fieldConfig:
      defaults:
        unit: bytes                 # Unit (bytes, percent, s, ops, etc.)
        decimals: 2
        min: 0
        max: 100
        color:
          mode: thresholds          # thresholds, palette-classic, fixed
        thresholds:
          mode: absolute            # absolute, percentage
          steps:
            - color: green
            - value: 80
              color: yellow
            - value: 90
              color: red
        custom:                     # Visualization-specific
          drawStyle: line           # line, bars, points
          lineInterpolation: smooth
          fillOpacity: 20
          stacking:
            mode: normal            # none, normal, percent
    
    # Repeat panel for each value of variable
    repeat: namespace
    repeatDirection: horizontal     # horizontal, vertical
    maxPerRow: 4
```

## Available Panel Templates

| Template | Type | Description |
|----------|------|-------------|
| `timeseries-basic` | timeseries | Basic line chart |
| `timeseries-stacked` | timeseries | Stacked area chart |
| `timeseries-rate` | timeseries | For rate metrics (ops/sec) |
| `timeseries-bytes` | timeseries | For byte-based metrics |
| `timeseries-percent` | timeseries | For percentage metrics (0-100%) |
| `timeseries-duration` | timeseries | For latency/duration metrics |
| `gauge-percent` | gauge | Percentage gauge with thresholds |
| `gauge-basic` | gauge | Basic gauge |
| `stat-value` | stat | Large single value |
| `stat-percent` | stat | Percentage stat with coloring |
| `stat-uptime` | stat | Uptime percentage (SLA thresholds) |
| `stat-count` | stat | Counter/count value |
| `table-basic` | table | Basic data table |
| `table-logs` | table | Table for log data |
| `barchart-horizontal` | barchart | Horizontal bar chart |
| `barchart-vertical` | barchart | Vertical bar chart |
| `piechart-basic` | piechart | Pie chart |
| `text-markdown` | text | Markdown text panel |
| `logs-basic` | logs | Log viewer |
| `heatmap-basic` | heatmap | Heatmap visualization |
| `alertlist-basic` | alertlist | Alert list |

List templates: `go run ./cmd list-templates`

## GitOps & Versioning

### Workflow

```
┌─────────────────┐    ┌────────────────┐    ┌─────────────────┐
│  YAML Dashboard │ -> │   Generator    │ -> │ GrafanaDashboard│
│   Definitions   │    │                │    │       CRs       │
└─────────────────┘    └────────────────┘    └─────────────────┘
         │                                           │
         │                                           │
         v                                           v
┌─────────────────┐                        ┌─────────────────┐
│    Git Repo     │                        │   Kubernetes    │
│  (Versioning)   │                        │(Grafana Operator)
└─────────────────┘                        └─────────────────┘
```

### Version Control Best Practices

1. **Store YAML definitions in Git** - These are the source of truth
2. **Generate CRs in CI/CD** - Don't commit generated JSON/CRs
3. **Use branches for changes** - Create PRs for dashboard updates
4. **Tag releases** - Version your dashboard definitions

### Kustomize Multi-Environment Setup

```
kustomize/
├── base/
│   └── kustomization.yaml      # Common configuration
└── overlays/
    ├── dev/                    # Development Grafana
    │   └── kustomization.yaml
    ├── staging/                # Staging Grafana
    │   └── kustomization.yaml
    └── prod/                   # Production Grafana
        └── kustomization.yaml
```

Each overlay can:
- Target different Grafana instances
- Override namespaces
- Add environment-specific labels
- Make dashboards read-only in production

### Dashboard Restoration

Since dashboards are defined in YAML and stored in Git:

1. **Restore from Git history**:
   ```bash
   git checkout <commit> -- dashboards/my-dashboard.yaml
   go run ./cmd generate -i dashboards/my-dashboard.yaml -o output/
   kubectl apply -f output/my-dashboard-cr.yaml
   ```

2. **Use Git tags for releases**:
   ```bash
   git tag -a v1.0.0 -m "Dashboard release 1.0.0"
   git push origin v1.0.0
   ```

3. **Rollback with Kustomize**:
   ```bash
   git checkout v1.0.0
   kubectl apply -k kustomize/overlays/prod
   ```

## Converting Existing Dashboards

Convert existing Grafana JSON dashboards to YAML:

```bash
go run ./cmd convert -i existing-dashboard.json -o converted.yaml
```

Then review and enhance the converted YAML with templates.

## CLI Commands

```bash
# Generate dashboards
go run ./cmd generate -i dashboard.yaml -o output/     # Single file
go run ./cmd generate -d dashboards/ -o output/        # Directory

# List available templates
go run ./cmd list-templates
go run ./cmd list-templates -f json                    # JSON output

# Validate dashboard definition
go run ./cmd validate -i dashboard.yaml

# Initialize new project
go run ./cmd init -o my-project

# Convert existing JSON to YAML
go run ./cmd convert -i dashboard.json -o dashboard.yaml
```

## Integration with Argo CD

This tool works seamlessly with Argo CD for GitOps deployment:

```yaml
# Argo CD Application for dashboards
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: grafana-dashboards
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/your-org/dashboards.git
    targetRevision: main
    path: kustomize/overlays/prod
  destination:
    server: https://kubernetes.default.svc
    namespace: monitoring
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

## Directory Structure

```
hack/grafana-dashboards/
├── cmd/
│   └── main.go                 # CLI tool
├── pkg/
│   ├── generator/
│   │   └── generator.go        # YAML to JSON converter
│   ├── schema/
│   │   └── types.go            # YAML DSL type definitions
│   └── templates/
│       └── panels.go           # Panel templates
├── examples/
│   └── dashboards/
│       ├── argocd-overview.yaml
│       └── kubernetes-pods.yaml
├── kustomize/
│   ├── base/
│   └── overlays/
│       ├── dev/
│       ├── staging/
│       └── prod/
└── output/                      # Generated files
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add new panel templates or enhance the DSL
4. Submit a pull request

## License

Apache 2.0
