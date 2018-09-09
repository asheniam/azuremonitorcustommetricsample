package main

import (
	"os"
)

func main() {
	if len(os.Args) < 3 {
		loginfof("%s <resourceId> <resourceRegion>", os.Args[0])
		os.Exit(2)
	}

	resourceID := os.Args[1]
	resourceRegion := os.Args[2]

	config := &Config{}
	err := config.loadConfig("secret.yml")
	if err != nil {
		logerror(err)
		os.Exit(1)
	}

	client := NewAzureClient(config, &publicAzureEnvironment)
	client.SendCustomMetric(
		resourceID,
		resourceRegion,
		"SampleMetric1",
		"SampleMetricNamespace",
		1,
		1,
		1,
		1)

	client.SendCustomMetricWithDimensions(
		resourceID,
		resourceRegion,
		"SampleMetric2",
		"SampleMetricNamespace",
		"SampleDimension",
		"SampleDimensionValue",
		1,
		1,
		1,
		1)

	client.ReadAndPrintCustomMetric(
		resourceID,
		"SampleMetric1",
		"SampleMetricNamespace",
		"")

	client.ReadAndPrintCustomMetric(
		resourceID,
		"SampleMetric2",
		"SampleMetricNamespace",
		"SampleDimension%20eq%20'*'")
}
