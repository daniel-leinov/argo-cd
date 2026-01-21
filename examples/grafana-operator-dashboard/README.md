# Grafana Dashboard Management with GitOps

This example demonstrates a user-friendly way to manage Grafana Dashboards using the Grafana Operator and Argo CD (or any GitOps tool).

## The Problem

The `GrafanaDashboard` Custom Resource Definition (CRD) typically requires embedding the dashboard JSON as a string within the YAML `spec.json` field. This makes it difficult to:
- Edit the dashboard (no syntax highlighting, escaping issues).
- Review changes in Git (diffs are messy).
- Copy-paste from Grafana UI exports.

## The Solution

We use **Kustomize** to separate the dashboard JSON from the CRD.

1.  **`dashboard.json`**: The raw Grafana dashboard JSON file. You can export this directly from Grafana or generate it.
2.  **`kustomization.yaml`**: 
    - Generates a `ConfigMap` containing the JSON file.
    - Uses `replacements` to inject the generated ConfigMap name (with hash) into the `GrafanaDashboard` CRD.
3.  **`grafanadashboard.yaml`**: The CRD definition that references the ConfigMap.

## How to Use

1.  **Export** your dashboard from Grafana as JSON and save it as `dashboard.json`.
2.  **Commit** the changes to Git.
3.  **Sync** via Argo CD.

Argo CD will run `kustomize build` which produces:
- A `ConfigMap` named `dashboard-cm-<hash>` containing the JSON.
- A `GrafanaDashboard` resource referencing that specific ConfigMap.

When you update `dashboard.json`:
1.  Kustomize generates a new ConfigMap name (new hash).
2.  The `GrafanaDashboard` is updated to point to the new ConfigMap.
3.  The Grafana Operator detects the change and updates the dashboard in Grafana.

## Files

- `dashboard.json`: Sample dashboard.
- `grafanadashboard.yaml`: The Custom Resource.
- `kustomization.yaml`: The glue that binds them.
