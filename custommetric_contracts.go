package main

// CustomMetricRequest Top level request body to push custom metrics to Azure Monitor
type CustomMetricRequest struct {
	Time string                  `json:"time"`
	Data CustomMetricDataPayload `json:"data"`
}

// CustomMetricDataPayload Data payload envelope of custom metrics
type CustomMetricDataPayload struct {
	BaseData CustomMetricBaseDataPayload `json:"baseData"`
}

// CustomMetricBaseDataPayload Metric data of custom metrics
type CustomMetricBaseDataPayload struct {
	Metric    string                      `json:"metric"`
	Namespace string                      `json:"namespace"`
	DimNames  []string                    `json:"dimNames"`
	Series    []CustomMetricSeriesPayload `json:"series"`
}

// CustomMetricSeriesPayload Time series data for custom metrics
type CustomMetricSeriesPayload struct {
	DimValues []string `json:"dimValues"`
	Min       float64  `json:"min"`
	Max       float64  `json:"max"`
	Sum       float64  `json:"sum"`
	Count     int64    `json:"count"`
}
