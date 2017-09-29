package promadapter

import (
	"github.com/prometheus/common/model"
	promremote "github.com/prometheus/prometheus/storage/remote"
	"github.com/solarwinds/p2l/librato"
)

//
// The adapter package holds information necessary to convert from Prometheus types to types defined in the Librato
// client library.
//

func PromDataToLibratoMeasurements(req *promremote.WriteRequest) librato.MeasurementCollection {
	return samplesToMeasurementSubmission(writeRequestToSamples(req))
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

// samplesToMeasurementSubmission converts Prometheus common model Samples to a collection of Librato Measurements
func samplesToMeasurementSubmission(samples model.Samples) librato.MeasurementCollection {
	measurements := make(librato.MeasurementCollection, len(samples))

	for _, s := range samples {
		m := &librato.Measurement{
			Name:  s.Metric.String(),
			Value: float64(s.Value),
			Time:  int64(s.Timestamp),
		}
		measurements = append(measurements, m)
	}
	return measurements
}
