package task

import (
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/stats"
)

func CollectChanStats() {
	for pub, ch := range uploadCreativeChan {
		stats.CollectMaterialChanSize(pub, model.HANDLER_METHOD_UPLOAD, len(ch))
	}
	for pub, ch := range queryCreativeChan {
		stats.CollectMaterialChanSize(pub, model.HANDLER_METHOD_QUERY, len(ch))
	}
}
