package truenas

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAppClient_Methods(t *testing.T) {
	// This is a unit test to verify the AppClient methods compile correctly
	// Integration tests would require a live TrueNAS instance

	client := &Client{}
	appClient := NewAppClient(client)

	// Verify the client is initialized properly
	if appClient == nil {
		t.Fatal("NewAppClient returned nil")
	}

	if appClient.client != client {
		t.Error("AppClient.client not set correctly")
	}
}

func TestAppState_Constants(t *testing.T) {
	// Test that our AppState constants are defined correctly
	states := []AppState{
		AppStateRunning,
		AppStateStopped,
		AppStateDeploying,
		AppStateError,
		AppStateUpgrading,
		AppStatePending,
	}

	expected := []string{
		"RUNNING",
		"STOPPED",
		"DEPLOYING",
		"ERROR",
		"UPGRADING",
		"PENDING",
	}

	for i, state := range states {
		if string(state) != expected[i] {
			t.Errorf("AppState constant %d: expected %s, got %s", i, expected[i], string(state))
		}
	}
}

func TestApp_StructFields(t *testing.T) {
	// Test that App struct has all expected fields
	app := App{
		ID:           "test-app",
		Name:         "Test App",
		State:        AppStateRunning,
		Version:      "1.0.0",
		ChartName:    "test-chart",
		Namespace:    "ix-test",
		Catalog:      "TRUENAS",
		CatalogTrain: "stable",
	}

	if app.ID != "test-app" {
		t.Errorf("App.ID: expected 'test-app', got %s", app.ID)
	}

	if app.State != AppStateRunning {
		t.Errorf("App.State: expected %s, got %s", AppStateRunning, app.State)
	}
}

func TestAppClient_List(t *testing.T) {
	endpoint := os.Getenv("TRUENAS_ENDPOINT")
	apiKey := os.Getenv("TRUENAS_API_KEY")
	t.Logf("Using endpoint: %s", endpoint)

	client, err := NewClient(endpoint, Options{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	t.Logf("Client connected to %s", endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	apps, err := client.App.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, app := range apps {
		t.Logf("App: %s, State: %s, Version: %s\n", app.Name, app.State, app.Version)
	}
}

func ExampleAppClient_Get() {
	// This example shows how to get a specific app by name

	// client, err := truenas.NewClient("ws://truenas.local/websocket", truenas.Options{
	//     APIKey: "your-api-key",
	// })
	// if err != nil {
	//     panic(err)
	// }
	// defer client.Close()
	//
	// ctx := context.Background()
	// app, err := client.App.Get(ctx, "plex")
	// if err != nil {
	//     panic(err)
	// }
	//
	// fmt.Printf("Plex app state: %s\n", app.State)
}

func ExampleAppClient_QueryByState() {
	// This example shows how to query apps by their state

	// client, err := truenas.NewClient("ws://truenas.local/websocket", truenas.Options{
	//     Username: "admin",
	//     Password: "password",
	// })
	// if err != nil {
	//     panic(err)
	// }
	// defer client.Close()
	//
	// ctx := context.Background()
	// runningApps, err := client.App.QueryByState(ctx, truenas.AppStateRunning)
	// if err != nil {
	//     panic(err)
	// }
	//
	// fmt.Printf("Found %d running apps\n", len(runningApps))
}
