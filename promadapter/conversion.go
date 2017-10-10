package promadapter

import (
	"math"

	"time"

	"github.com/prometheus/common/model"
	promremote "github.com/prometheus/prometheus/storage/remote"
	"github.com/solarwinds/p2l/librato"
)

//
// The adapter package holds information necessary to convert from Prometheus types to types defined in the Librato
// client library, as well as for creating API-compliant batches and using the Librato client to send them.
//

func PromDataToLibratoMeasurements(req *promremote.WriteRequest) []*librato.Measurement {
	return samplesToMeasurements(writeRequestToSamples(req))
}

// writeRequestToSamples converts a Prometheus remote storage WriteRequest to a collection of Prometheus common model Samples
func writeRequestToSamples(req *promremote.WriteRequest) model.Samples {
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

// samplesToMeasurements converts Prometheus common model Samples to a collection of Librato Measurements
func samplesToMeasurements(samples model.Samples) []*librato.Measurement {
	var measurements []*librato.Measurement
	for _, s := range samples {
		if math.IsNaN(float64(s.Value)) {
			continue
		}

		msTime := time.Duration(s.Timestamp) / time.Microsecond

		m := &librato.Measurement{
			Name:  string(s.Metric[model.MetricNameLabel]),
			Value: float64(s.Value),
			Time:  int64(msTime),
			Tags:  labelsToTags(s),
		}
		measurements = append(measurements, m)
	}
	return measurements
}

// labelsToTags converts the Metric's associated Labels to Librato Tags
func labelsToTags(sample *model.Sample) librato.MeasurementTags {
	var mt = make(librato.MeasurementTags)
	for k, v := range sample.Metric {
		if k == model.MetricNameLabel {
			continue
		}
		mt[string(k)] = string(v)
	}
	return mt
}
