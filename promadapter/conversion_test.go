package promadapter

import (
	"testing"
	"time"

	promremote "github.com/prometheus/prometheus/storage/remote"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

var (
	metricNameFixture = "rpc_widget_count"
	valueFixture      = 3.1415
	timestampFixture  = time.Now().UTC().Unix()
)

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

		assert.Equal(t, castTime, measurement.Time)
		assert.Equal(t, castName, measurement.Name)
		assert.Equal(t, castValue, measurement.Value)

	}
}

func TestLabelsToTags(t *testing.T) {
	sample := promSamples[0]
	adapter := NewPromAdapter()
	tags := adapter.LabelsToTags(sample)
	for k, v := range tags {
		labelName := model.LabelName(k)
		assert.Equal(t, model.LabelValue(v), sample.Metric[labelName])
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
		assert.Equal(t, int(s.Value.(float64)), int(samples[0].Value))
	}
}
