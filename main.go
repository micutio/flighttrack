package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// lat is Latitude of SIN Airport.
	lat float64 = 1.359297
	// lon is Longitude of SIN Airport.
	lon float64 = 103.989348
	// aircraftUpdateInterval determines the update rate for general aircraft.
	aircraftUpdateInterval = 30 * time.Second
	// milAircraftUpdateInterval determines the update rate for military aircraft.
	milAircraftUpdateInterval = 15 * time.Minute
	// milAircraftUpdateDelay determines interleaving between general and mil aircraft api calls.
	milAircraftUpdateDelay = 15 * time.Second
	// summaryInterval determines how often the summary is show.
	summaryInterval = 1 * time.Hour
)

var (
	errNonOkResponse     error = errors.New("non-OK response")
	errEmptyResponseBody error = errors.New("empty response body")
	errNonJSONContent    error = errors.New("non-JSON content type")
)

func main() {
	logger := slog.Default()
	flightDash, dashboardErr := newDashboard()
	if dashboardErr != nil {
		logger.Error("unable to create dashboard, exiting", slog.Any("dashboard error", dashboardErr))
		os.Exit(1)
	}

	// Create a aircraftUpdateTicker that fires every 30 seconds
	aircraftUpdateTicker := time.NewTicker(aircraftUpdateInterval)
	defer aircraftUpdateTicker.Stop()

	// aircraft and military aircraft updates should not coincide to avoid exceeding the api limit.
	// Hence, stagger them by 15 seconds.
	milAircraftUpdateTicker := time.NewTicker(milAircraftUpdateInterval)
	defer milAircraftUpdateTicker.Stop()
	time.AfterFunc(milAircraftUpdateDelay, func() {
		milAircraftUpdateTicker.Reset(milAircraftUpdateInterval)
	})

	summaryTicker := time.NewTicker(summaryInterval)
	defer summaryTicker.Stop()

	// Use a channel to gracefully stop the program if needed.
	// (Though not strictly necessary for an infinite loop)
	done := make(chan bool)

	fmt.Println("aircraft Tracking dashboard")
	fmt.Println("Press Ctrl+C to stop the program.")

	// Start a goroutine to perform the requests
	go func() {
		for {
			select {
			case <-aircraftUpdateTicker.C:
				if err := requestAndProcessCivAircraft(flightDash); err != nil {
					logger.Error("main: ", slog.Any("error", err))
				}
			case <-milAircraftUpdateTicker.C:
				if err := requestAndProcessMilAircraft(flightDash); err != nil {
					logger.Error("main: %w", slog.Any("error", err))
				}
			case <-summaryTicker.C:
				flightDash.listTypesByRarity()
			case <-done:
				// This case allows for graceful shutdown (not used in this example but good practice)
				logger.Info("Stopping HTTP GET request routine.")

				return
			}
		}
	}()

	// Run once in the beginning.
	if err := requestAndProcessMilAircraft(flightDash); err != nil {
		logger.Error("main: ", slog.Any("error", err))
	}

	// Keep the main goroutine alive indefinitely, or until an interrupt signal is received
	// In a real application, you might have other logic here or use a wait group.
	select {} // Block indefinitely
}

func requestAndProcessCivAircraft(dashboard *dashboard) error {
	// Define the URL for the HTTP GET request
	targetURL := fmt.Sprintf(
		"https://opendata.adsb.fi/api/v2/lat/%.6f/lon/%.6f/dist/250",
		lat,
		lon,
	)

	// This case is executed every time the ticker "ticks"
	body, requestErr := sendRequest(targetURL)
	if requestErr != nil {
		return fmt.Errorf("requestAndProcessCivAircraft: error during request: %w", requestErr)
	}

	dashboard.processCivAircraftJSON(body)

	return nil
}

func requestAndProcessMilAircraft(dashboard *dashboard) error {
	// Define the URL for the HTTP GET request
	targetURL := "https://opendata.adsb.fi/api/v2/mil"
	// This case is executed every time the ticker "ticks"
	body, requestErr := sendRequest(targetURL)
	if requestErr != nil {
		return fmt.Errorf("error during request: %w", requestErr)
	}

	dashboard.processMilAircraftJSON(body)

	return nil
}

// sendRequest sends an HTTP GET request and returns a valid byte slice of the response body.
func sendRequest(url string) ([]byte, error) {
	ctx := context.Background()
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("sendRequest: invalid request error: %s : %w", url, reqErr)
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return nil, fmt.Errorf("sendRequest: failed to send GET request: %s: %w", url, respErr)
	}
	defer resp.Body.Close()
	// defer func(bodyReader io.ReadCloser) {
	// 	err := bodyReader.Close()
	// 	if err != nil {
	//		slog.Default().Error("sendRequest: ", slog.Any("close body reader", err))
	//	}
	// }(resp.Body) // Ensure the response body is closed

	// Check if the request was successful (status code 200 OK)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sendRequest: %w %s", errNonOkResponse, resp.Status)
	}

	// Read the response body
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to read response body: %w", bodyErr)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("sendRequest: %w", errEmptyResponseBody)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("senRequest: %w, %s", errNonJSONContent, contentType)
	}

	return body, nil
}
