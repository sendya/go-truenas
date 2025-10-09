# TrueNAS Go Client - Application Management

This document provides details on using the AppClient for managing TrueNAS applications.

## AppClient Usage Examples

### Basic App Queries

```go
// List all apps
apps, err := client.App.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, app := range apps {
    fmt.Printf("App: %s, State: %s, Version: %s\n", app.Name, app.State, app.Version)
}

// Get a specific app by name
app, err := client.App.Get(ctx, "plex")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Plex state: %s\n", app.State)

// Get app by ID
app, err := client.App.GetByID(ctx, "ix-plex")
if err != nil {
    log.Fatal(err)
}
```

### Filtering Apps

```go
// Get running apps
runningApps, err := client.App.ListRunning(ctx)

// Get stopped apps
stoppedApps, err := client.App.ListStopped(ctx)

// Get apps in error state
errorApps, err := client.App.ListWithErrors(ctx)

// Filter by catalog
truenasApps, err := client.App.QueryByCatalog(ctx, "TRUENAS")

// Filter by state
deployingApps, err := client.App.QueryByState(ctx, truenas.AppStateDeploying)
```

### Advanced Queries

```go
// Custom filters
filters := [][]any{
    {"state", "=", "RUNNING"},
    {"catalog", "=", "TRUENAS"},
}
options := &truenas.AppQueryOptions{
    ExtraOptions: map[string]interface{}{
        "include_history": true,
    },
}
apps, err := client.App.QueryWithFilters(ctx, filters, options)
```

## App States

The following app states are available:

- `AppStateRunning` - App is running normally
- `AppStateStopped` - App is stopped
- `AppStateDeploying` - App is being deployed
- `AppStateError` - App is in error state
- `AppStateUpgrading` - App is being upgraded
- `AppStatePending` - App deployment is pending

## App Structure

The `App` struct contains comprehensive information about TrueNAS applications:

```go
type App struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    State       AppState               `json:"state"`
    Version     string                 `json:"version"`
    Upgrade     *bool                  `json:"upgrade,omitempty"`
    ChartName   string                 `json:"chart_name"`
    Namespace   string                 `json:"namespace"`
    Catalog     string                 `json:"catalog"`
    CatalogTrain string                `json:"catalog_train"`
    Config      map[string]interface{} `json:"config,omitempty"`
    History     []AppHistory           `json:"history,omitempty"`
    Resources   *AppResources          `json:"resources,omitempty"`
    Metadata    *AppMetadata           `json:"metadata,omitempty"`
    // ... additional fields
}
```

## API Methods Implemented

- `app.query` - Query applications with filters and options

## Error Handling

The AppClient methods return `NewNotFoundError` when apps are not found:

```go
app, err := client.App.Get(ctx, "nonexistent-app")
if err != nil {
    var notFoundErr *truenas.NotFoundError
    if errors.As(err, &notFoundErr) {
        fmt.Printf("App not found: %s\n", notFoundErr.Error())
    } else {
        log.Printf("Other error: %v\n", err)
    }
}
```
