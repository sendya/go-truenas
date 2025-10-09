# TrueNAS Go Client - Application Management

This document provides details on using the AppClient for managing TrueNAS applications based on the official TrueNAS API documentation.

## AppClient Usage Examples

### Basic App Queries

```go
// List all apps
apps, err := client.App.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, app := range apps {
    fmt.Printf("App: %s, State: %s, Version: %s\n", app.Name, app.State, app.HumanVersion)
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

````go
// Get running apps
runningApps, err := client.App.ListRunning(ctx)

// Get stopped apps
stoppedApps, err := client.App.ListStopped(ctx)

// Get apps in crashed state
crashedApps, err := client.App.ListCrashed(ctx)

// Filter by state
deployingApps, err := client.App.QueryByState(ctx, truenas.AppStateDeploying)
```### Advanced Queries

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
````

## App States

The following app states are available according to TrueNAS API documentation:

- `AppStateCrashed` - App has crashed
- `AppStateDeploying` - App is being deployed
- `AppStateRunning` - App is running normally
- `AppStateStopped` - App is stopped

## App Structure

The `App` struct contains information about TrueNAS applications based on the official API schema:

```go
type App struct {
    Name             string                 `json:"name"`
    ID               string                 `json:"id"`
    State            AppState               `json:"state"`
    UpgradeAvailable bool                   `json:"upgrade_available"`
    HumanVersion     string                 `json:"human_version"`
    Version          string                 `json:"version"`
    Metadata         map[string]interface{} `json:"metadata"`
    ActiveWorkloads  *AppActiveWorkloads    `json:"active_workloads"`
}
```

### Active Workloads

The `ActiveWorkloads` field provides detailed runtime information:

```go
type AppActiveWorkloads struct {
    Containers       int                   `json:"containers"`
    UsedPorts        []AppUsedPort         `json:"used_ports"`
    ContainerDetails []AppContainerDetail  `json:"container_details"`
    Volumes          []AppVolume           `json:"volumes"`
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
