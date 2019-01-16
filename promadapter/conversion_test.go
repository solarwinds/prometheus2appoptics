package promadapter

import (
	promremote "github.com/prometheus/prometheus/storage/remote"
	"testing"
	"time"

	"github.com/prometheus/common/model"
)

var metricNameFixture = "rpc_widget_count"
var valueFixture = 3.1415
var timestampFixture = time.Now().UTC().Unix()

var labels = model.LabelSet{
	model.LabelName(model.MetricNameLabel): model.LabelValue(metricNameFixture),
	model.LabelName("environment"):         model.LabelValue("production"),
	model.LabelName("job"):                 model.LabelValue("inventory-service"),
}

var promSamples = model.Samples{
	&model.Sample{
		Metric:    model.Metric(labels),
		Value:     model.SampleValue(valueFixture),
		Timestamp: model.Time(timestampFixture),
	},
}

func TestSamplesToMeasurements(t *testing.T) {
	adapter := NewPromAdapter()
	ms := adapter.SamplesToMeasurements(promSamples)

	for i, measurement := range ms {
		castTime := int64((time.Duration(promSamples[i].Timestamp)) / time.Microsecond)
		castValue := float64(promSamples[i].Value)
		castName := string(promSamples[i].Metric[model.MetricNameLabel])

		if castTime != measurement.Time {
			t.Errorf("expected %d to match %d", castTime, measurement.Time)
		}

		if castValue != measurement.Value {
			t.Errorf("expected %f to match %f", castValue, measurement.Value)
		}

		if castName != measurement.Name {
			t.Errorf("expected %s to match %s", castName, measurement.Name)
		}
	}
}

func TestLabelsToTags(t *testing.T) {
	sample := promSamples[0]
	adapter := NewPromAdapter()
	tags := adapter.LabelsToTags(sample)
	for k, v := range tags {
		labelName := model.LabelName(k)
		if model.LabelValue(v) != sample.Metric[labelName] {
			t.Errorf("expected %s to map to %s but it didn't", k, sample.Metric[labelName])
		}

	}
}

func TestWriteRequestToSamples(t *testing.T) {
	samples := []*promremote.Sample{{
		Value:       valueFixture,
		TimestampMs: timestampFixture,
	}}

	labels := []*promremote.LabelPair{{Name: "Label", Value: "Value"}}

	tss := []*promremote.TimeSeries{{Samples: samples, Labels: labels}}
	wr := promremote.WriteRequest{tss}
	adapter := NewPromAdapter()

	modelSamples := adapter.PromDataToAppOpticsMeasurements(&wr)

	for _, s := range modelSamples {
		if int(s.Value.(float64)) != int(samples[0].Value) {
			t.Errorf("expected %s to map to %s but it didn't", s, samples[0].Value)
		}
	}
}
