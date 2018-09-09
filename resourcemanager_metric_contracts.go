package main

// ResourceManagerMetricResponse Top level response body of metrics behind Azure Resource Manager when reading metric data
type ResourceManagerMetricResponse struct {
	Cost           int                     `json:"cost"`
	Timespan       string                  `json:"timespan"`
	Interval       string                  `json:"interval"`
	Namespace      string                  `json:"namespace"`
	ResourceRegion string                  `json:"resourceregion"`
	Metrics        []ResourceManagerMetric `json:"value"`
}

// ResourceManagerMetric Metric object
type ResourceManagerMetric struct {
	ID       string                        `json:"id"`
	Name     ResourceManagerLocalizedName  `json:"name"`
	Unit     string                        `json:"unit"`
	Timespan string                        `json:"timespan"`
	Series   []ResourceManagerMetricSeries `json:"timeseries"`
}

// ResourceManagerLocalizedName Localized value
type ResourceManagerLocalizedName struct {
	Value         string `json:"value"`
	LocalizedName string `json:"localizedValue"`
}

// ResourceManagerMetricSeries Metric series data
type ResourceManagerMetricSeries struct {
	Metadata []ResourceManagerMetadataValue   `json:"metadatavalues"`
	Data     []ResourceManagerMetricDataValue `json:"data"`
}

// ResourceManagerMetadataValue Metric dimensions
type ResourceManagerMetadataValue struct {
	Name  ResourceManagerLocalizedName `json:"name"`
	Value string                       `json:"value"`
}

// ResourceManagerMetricDataValue Metric data
type ResourceManagerMetricDataValue struct {
	Timestamp string  `json:"timeStamp"`
	Total     float64 `json:"total"`
	Count     float64 `json:"count"`
	Average   float64 `json:"average"`
	Maximum   float64 `json:"maximum"`
	Minimum   float64 `json:"minimum"`
}
