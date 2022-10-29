package stats

import (
	"github.com/convee/adcreative/internal/service"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	materialTaskCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "material_task_status_gauge",
	}, []string{"consumer", "advertiser", "publisher", "status"})
	taskChanGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "material_task_chan_gauge",
	}, []string{"publisher", "type"})
	publisherApiSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "publisher_api_summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001},
	}, []string{"publisher", "api"})
)

func init() {
	prometheus.MustRegister(materialTaskCounter)
	prometheus.MustRegister(publisherApiSummary)
	prometheus.MustRegister(taskChanGauge)
}

func CollectMaterialTaskStatus(consumer string, advertiser int, pub string, status string, count int) {
	materialTaskCounter.WithLabelValues(consumer, strconv.Itoa(advertiser), pub, status).Set(float64(count))
}

func CollectMaterialTaskReset() {
	materialTaskCounter.Reset()
}
func PublisherApiObserve(pub int, api string, costMs int64) {
	publisherApiSummary.WithLabelValues(service.GetPubName(pub), api).Observe(float64(costMs))
}

func CollectMaterialChanSize(pub int, typ string, count int) {
	taskChanGauge.WithLabelValues(service.GetPubName(pub), typ).Set(float64(count))
}
