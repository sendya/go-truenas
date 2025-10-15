package truenas

import (
	"context"
	"fmt"
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
	// Test that our AppState constants are defined correctly according to TrueNAS API
	states := []AppState{
		AppStateCrashed,
		AppStateDeploying,
		AppStateRunning,
		AppStateStopped,
	}

	expected := []string{
		"CRASHED",
		"DEPLOYING",
		"RUNNING",
		"STOPPED",
	}

	for i, state := range states {
		if string(state) != expected[i] {
			t.Errorf("AppState constant %d: expected %s, got %s", i, expected[i], string(state))
		}
	}
}

func TestApp_StructFields(t *testing.T) {
	// Test that App struct has all expected fields according to TrueNAS API
	app := App{
		ID:               "test-app",
		Name:             "Test App",
		State:            AppStateRunning,
		Version:          "1.0.0",
		HumanVersion:     "1.0.0",
		UpgradeAvailable: false,
		Metadata:         map[string]interface{}{},
	}

	if app.ID != "test-app" {
		t.Errorf("App.ID: expected 'test-app', got %s", app.ID)
	}

	if app.State != AppStateRunning {
		t.Errorf("App.State: expected %s, got %s", AppStateRunning, app.State)
	}
}

func TestAppClient_GetApp(t *testing.T) {
	extra := map[string]any{
		// "host_ip":            "nas.tooko.io",
		"include_app_schema": true,
		"retrieve_config":    true,
	}

	endpoint := "" // put your TrueNAS endpoint here for testing
	apiKey := ""   // put your TrueNAS API key here for testing

	t.Logf("Using endpoint: %s", endpoint)
	client, err := NewClient(endpoint, Options{
		APIKey: apiKey,
		Debug:  false,
	})
	if err != nil {
		t.Fatal(err)
	}

	appClient := NewAppClient(client)
	app, err := appClient.Get(context.Background(), "grafana", extra)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("app = \n%s", tryMarshal(app))
}

func ExampleAppClient_List() {
	endpoint := os.Getenv("TRUENAS_ENDPOINT")
	apiKey := os.Getenv("TRUENAS_API_KEY")
	fmt.Printf("Using endpoint: %s", endpoint)

	client, err := NewClient(endpoint, Options{
		APIKey: apiKey,
		Debug:  false,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Printf("Client connected to %s", endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.App.SubscribeStats(ctx, func(apps []AppStats) error {
		for _, app := range apps {
			fmt.Printf("App: %s, CPU Usage: %.2f%%, Memory: %dMiB \n", app.AppName, app.CPUUsage, app.Memory/1024/1024)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	time.Sleep(8 * time.Second)

	if err := client.App.UnsubscribeStats(context.Background()); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
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

// Query App details
// {"jsonrpc":"2.0","id":"a11f55f4-7b9f-2e4e-9cb5-ec4423d13e3b","method":"app.query","params":[[["name","=","grafana"]],{"extra":{"include_app_schema":true,"retrieve_config":true,"host_ip":"nas.tooko.io"}}]}
