# Grafana dashboard management with Grafana Operator

This document proposes a user-friendly workflow for creating, versioning,
restoring, and persisting Grafana dashboards when you use the Grafana Operator.
It focuses on making dashboards easy to edit and review by avoiding large,
hand-edited JSON files.

## Goals

- Provide a clean authoring experience for dashboards at scale.
- Keep dashboards versioned and reviewable in Git.
- Make restore/rollback predictable with GitOps.
- Ensure Grafana data is persisted across restarts.
- Work with the Grafana Operator (no bespoke controllers required).

## Non-goals

- Replace the Grafana UI or build a custom dashboard editor.
- Cover all Grafana resources (alerting rules, data sources, etc.).
- Support every Grafana Operator version (use the fields your operator exposes).

## Proposed architecture

1. **Dashboard sources as code**: Store dashboards as Jsonnet (Grafonnet) or
   another dashboard DSL. Keep JSON generated, not hand-written.
2. **Build step**: Render dashboard sources into JSON files (or YAML
   manifests) in a `generated/` directory.
3. **Grafana Operator resources**: Use `GrafanaDashboard` resources that either
   embed JSON directly or reference a ConfigMap/Secret generated from the JSON
   files (depending on the operator version you run).
4. **GitOps delivery**: Use Argo CD (or another GitOps tool) to apply the
   Grafana Operator resources.
5. **Persistence**: Configure Grafana to use persistent storage (PVC or an
   external database). The operator keeps dashboards in sync; storage ensures
   the Grafana instance and related metadata persist across restarts.

## Repository layout (example)

```
dashboards/
  lib/                     # reusable panels, queries, variables
    panels/
    queries/
  dashboards/              # top-level dashboards
    argocd/
      overview.jsonnet
      app-health.jsonnet
  env/                     # optional environment overlays
    base.libsonnet
    prod.libsonnet
  generated/               # rendered JSON (generated output)
    argocd-overview.json
    argocd-app-health.json
  k8s/                     # Grafana Operator resources
    grafana-folder.yaml
    dashboard-template.yaml
    kustomization.yaml
```

## Authoring workflow

1. **Author** a dashboard in Jsonnet (Grafonnet) and compose from reusable
   components in `lib/`.
2. **Render** dashboards to JSON:
   - `make dashboards-build` (jsonnet or grizzly render)
3. **Validate** and format:
   - `make dashboards-validate` (jsonnetfmt, schema checks, or unit tests)
4. **Review** the small, readable source changes (Jsonnet), not huge JSON
   diffs.
5. **Deploy** with Argo CD. The Grafana Operator reconciles dashboards into
   Grafana.

Tip: Keep `generated/` committed if you want fully declarative GitOps and
easy diff visibility in CI. If you prefer generated artifacts out of Git,
render them in CI and apply via a build step.

## Versioning strategy

- Use Git history as the source of truth.
- Optionally add a version label on `GrafanaDashboard` resources, for example:
  - `app.kubernetes.io/version: "1.3.0"`
- Give dashboards stable `uid` values so URLs and references remain intact.

## Restore and rollback

- **Rollback** by reverting Git and re-syncing. The operator reconciles the
  prior dashboards into Grafana.
- **Disaster recovery**:
  - Restore the Grafana database or PVC backup.
  - Re-apply dashboards from Git to ensure everything is recreated.

## Persistence

Use one of the following options in your Grafana configuration (via the
Grafana Operator):

- **Persistent volume** for the Grafana data directory.
- **External database** (PostgreSQL/MySQL) for metadata, users, and dashboard
  history.

Even if dashboards are fully Git-managed, persistence protects:
users, annotations, org settings, and other metadata.

## UI editing without losing changes

If you want a user-friendly path from UI to code:

- Create a **sandbox** folder in Grafana for ad-hoc edits.
- Export dashboards from the UI using a tool such as `grizzly` or the Grafana
  API.
- Convert exports into Jsonnet modules and commit them.
- Keep production folders read-only by setting `editable: false` in dashboard
  templates and by using folder permissions.

## Minimal GrafanaDashboard example

The exact fields depend on your Grafana Operator version, but the shape
typically looks like this:

```yaml
apiVersion: grafana.integreatly.org/v1beta1
kind: GrafanaDashboard
metadata:
  name: argocd-overview
  labels:
    app.kubernetes.io/part-of: argocd
spec:
  folder: Argo CD
  json: |
    { ...generated dashboard JSON... }
```

## Risks and mitigations

- **Drift between Grafana UI and Git**:
  - Mitigate by restricting edits in managed folders and exporting from a
    sandbox folder only.
- **Large dashboard JSON**:
  - Mitigate by storing Jsonnet sources and keeping JSON generated.
- **Operator version differences**:
  - Mitigate by documenting the supported spec fields for your operator
    version and standardizing the build pipeline.

## Summary

This design keeps dashboards user-friendly and version-controlled by
separating authoring (Jsonnet) from delivery (Grafana Operator), while GitOps
handles rollout, rollback, and consistency. It scales to large teams by
encouraging modular sources, automated validation, and a clear UI-to-code
workflow.
