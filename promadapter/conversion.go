package promadapter

import (
	"math"

	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/prometheus/common/model"
	promremote "github.com/prometheus/prometheus/storage/remote"
)

//
// The adapter package holds information necessary to convert from Prometheus types to types defined in the AppOptics
// client library, as well as for creating API-compliant batches and using the AppOptics client to send them.
//

type PrometheusAdapter interface {
	WriteRequestToSamples(req *promremote.WriteRequest) model.Samples
	PromDataToAppOpticsMeasurements(req *promremote.WriteRequest) []appoptics.Measurement
	SamplesToMeasurements(samples model.Samples) []appoptics.Measurement
	LabelsToTags(sample *model.Sample) map[string]string
}

type Adapter struct {
	PrometheusAdapter
}

func NewPromAdapter() PrometheusAdapter {
	p := Adapter{}
	return PrometheusAdapter(&p)
}

func (p *Adapter) PromDataToAppOpticsMeasurements(req *promremote.WriteRequest) []appoptics.Measurement {
	return p.SamplesToMeasurements(p.WriteRequestToSamples(req))
}

// WriteRequestToSamples converts a Prometheus remote storage WriteRequest to a collection of Prometheus common model Samples
func (p *Adapter) WriteRequestToSamples(req *promremote.WriteRequest) model.Samples {
	var samples model.Samples
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, label := range ts.Labels {
			metric[model.LabelName(label.Name)] = model.LabelValue(label.Value)
		}

		for _, sample := range ts.Samples {
			s := &model.Sample{
				Metric:    metric,
				Value:     model.SampleValue(sample.Value),
				Timestamp: model.Time(sample.TimestampMs),
			}
			samples = append(samples, s)
		}
	}
	return samples
}

// SamplesToMeasurements converts Prometheus common model Samples to a collection of AppOptics Measurements
func (p *Adapter) SamplesToMeasurements(samples model.Samples) []appoptics.Measurement {
	var measurements []appoptics.Measurement
	for _, s := range samples {
		if math.IsNaN(float64(s.Value)) {
			continue
		}

		msTime := time.Duration(s.Timestamp) / time.Microsecond

		m := appoptics.Measurement{
			Name:  string(s.Metric[model.MetricNameLabel]),
			Value: float64(s.Value),
			Time:  int64(msTime),
			Tags:  p.LabelsToTags(s),
		}
		measurements = append(measurements, m)
	}
	return measurements
}

// LabelsToTags converts the Metric's associated Labels to AppOptics Tags
func (p *Adapter) LabelsToTags(sample *model.Sample) map[string]string {
	var mt = make(map[string]string)
	for k, v := range sample.Metric {
		if k == model.MetricNameLabel {
			continue
		}
		mt[string(k)] = string(v)
	}
	return mt
}
