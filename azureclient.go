package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	supportedRegions = map[string]int{
		"westus2":        1,
		"westcentralus":  1,
		"northeurope":    1,
		"westeurope":     1,
		"eastus":         1,
		"southcentralus": 1,
		"southeastasia":  1,
	}
)

// AzureClient Azure client to read and write metrics to Azure Monitor
type AzureClient struct {
	client        *http.Client
	tokenAccessor *AzureTokenAccessor
	config        *Config
	environment   *Environment
}

// NewAzureClient Create new Azure client
func NewAzureClient(config *Config, environment *Environment) *AzureClient {
	tokenAccessor := NewAzureTokenAccessor(environment, config.Credentials.TenantID)
	return &AzureClient{
		client: &http.Client{
			Timeout: time.Second * 45,
		},
		tokenAccessor: tokenAccessor,
		config:        config,
		environment:   environment,
	}
}

func (azureClient *AzureClient) sendArmHTTPMessage(method string, url string, body interface{}) (*http.Response, error) {
	token := azureClient.tokenAccessor.GetAccessToken(
		azureClient.config.Credentials.ClientID,
		azureClient.config.Credentials.ClientSecret,
		azureClient.environment.armURL,
	)
	token.EnsureFresh()
	return azureClient.sendHTTPMessage(method, url, token.Token().AccessToken, body)
}

func (azureClient *AzureClient) sendCustomMetricHTTPMessage(method string, url string, body interface{}) (*http.Response, error) {
	token := azureClient.tokenAccessor.GetAccessToken(
		azureClient.config.Credentials.ClientID,
		azureClient.config.Credentials.ClientSecret,
		azureClient.environment.ingestionURL,
	)
	token.EnsureFresh()
	return azureClient.sendHTTPMessage(method, url, token.Token().AccessToken, body)
}

func (azureClient *AzureClient) sendHTTPMessage(method string, url string, accessToken string, body interface{}) (*http.Response, error) {
	loginfof("Running %s %s", method, url)
	var requestBodyReader io.Reader
	if body != nil {
		bodyAsBytes, _ := json.Marshal(body)
		requestBodyReader = bytes.NewBuffer(bodyAsBytes)
	}

	request, err := http.NewRequest(method, url, requestBodyReader)
	if err != nil {
		logerrorf("Error creating HTTP request: %v", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := azureClient.client.Do(request)
	if err != nil {
		logerrorf("Error sending HTTP request: %v", err)
		return nil, err
	}

	loginfof("Status code: %d (%s %s)", response.StatusCode, method, url)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		logerrorf("Error sending HTTP request with status code: %d\n", response.StatusCode)

		// If unsuccessful HTTP status code, print out the error response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logerrorf("Error reading body of response: %v", err)
			return nil, err
		}

		loginfof("Length: %d", len(body))
		loginfof(string(body))
	}

	return response, nil
}

// ReadCustomMetric Read custom metrics
func (azureClient *AzureClient) ReadCustomMetric(resourceID string, metric string, namespace string, filter string) (*ResourceManagerMetricResponse, error) {
	timespan := time.Now().Add(time.Duration(-15)*time.Minute).UTC().Format(time.RFC3339) + "/" + time.Now().UTC().Format(time.RFC3339)
	url := fmt.Sprintf("%s%s/providers/microsoft.insights/metrics?api-version=2018-01-01&timespan=%s&interval=PT1M&aggregation=average,total,count&metricnamespace=%s&metricnames=%s",
		strings.TrimSuffix(azureClient.environment.armURL, "/"),
		resourceID,
		timespan,
		namespace,
		metric)

	if len(filter) > 0 {
		url += "&$filter=" + filter
	}

	response, err := azureClient.sendArmHTTPMessage("GET", url, nil)
	defer response.Body.Close()
	if err != nil {
		logerrorf("Error reading metric: %v", err)
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)
	loginfof("Length: %d", len(body))
	loginfof(string(body))

	var r ResourceManagerMetricResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		logerrorf("Error reading metric: %v", err)
		return nil, err
	}

	return &r, nil
}

// ReadAndPrintCustomMetric Read and print metrics
func (azureClient *AzureClient) ReadAndPrintCustomMetric(resourceID string, metric string, namespace string, filter string) {
	r, err := azureClient.ReadCustomMetric(resourceID, metric, namespace, filter)
	if err != nil {
		logerrorf("Unable to read custom metrics")
		return
	}

	for _, m := range r.Metrics {
		loginfof("Metric: %s", m.Name.Value)
		for _, s := range m.Series {
			for _, meta := range s.Metadata {
				loginfof("Dimension: %s=%s", meta.Name.Value, meta.Value)
			}
			for _, d := range s.Data {
				loginfof("%s: avg:%f, total:%f, count:%f", d.Timestamp, d.Average, d.Total, d.Count)
			}
		}
	}
}

// SendCustomMetric Write custom metrics without dimensions using current timestamp
func (azureClient *AzureClient) SendCustomMetric(resourceID string, resourceRegion string, metric string, namespace string, min float64, max float64, sum float64, count int64) error {
	_, ok := supportedRegions[resourceRegion]
	if !ok {
		logerrorf("%s region is not currently supported", resourceRegion)
		return errors.New("Invalid region")
	}

	request := CustomMetricRequest{
		Time: time.Now().Format(time.RFC3339),
		Data: CustomMetricDataPayload{
			BaseData: CustomMetricBaseDataPayload{
				Metric:    metric,
				Namespace: namespace,
				Series: []CustomMetricSeriesPayload{
					CustomMetricSeriesPayload{
						Min:   min,
						Max:   max,
						Sum:   sum,
						Count: count,
					},
				},
			},
		},
	}

	url := fmt.Sprintf(strings.TrimSuffix(azureClient.environment.ingestionURLTemplate, "/")+"%s/metrics",
		resourceRegion,
		resourceID)
	_, err := azureClient.sendCustomMetricHTTPMessage("POST", url, request)
	return err
}

// SendCustomMetricWithDimensions Write custom metrics with dimensions using current timestamp
func (azureClient *AzureClient) SendCustomMetricWithDimensions(resourceID string, resourceRegion string, metric string, namespace string, dimName string, dimValue string, min float64, max float64, sum float64, count int64) error {
	_, ok := supportedRegions[resourceRegion]
	if !ok {
		logerrorf("%s region is not currently supported", resourceRegion)
		return errors.New("Invalid region")
	}

	request := CustomMetricRequest{
		Time: time.Now().Format(time.RFC3339),
		Data: CustomMetricDataPayload{
			BaseData: CustomMetricBaseDataPayload{
				Metric:    metric,
				Namespace: namespace,
				DimNames:  []string{dimName},
				Series: []CustomMetricSeriesPayload{
					CustomMetricSeriesPayload{
						DimValues: []string{dimValue},
						Min:       min,
						Max:       max,
						Sum:       sum,
						Count:     count,
					},
				},
			},
		},
	}

	url := fmt.Sprintf(strings.TrimSuffix(azureClient.environment.ingestionURLTemplate, "/")+"%s/metrics",
		resourceRegion,
		resourceID)
	_, err := azureClient.sendCustomMetricHTTPMessage("POST", url, request)
	return err
}
