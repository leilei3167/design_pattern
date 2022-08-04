package pipeline

import (
	"fmt"
	"github.com/leilei3167/copy_design_pattern/monitor/plugin"
	"github.com/panjf2000/ants"
)

var pool, _ = ants.NewPool(5)

type PoolPipeline struct {
	pipelineTemplate
}

func (p *PoolPipeline) SetContext(ctx plugin.Context) {
	p.run = func() {
		if err := pool.Submit(p.doRun); err != nil {
			fmt.Printf("PoolPipeine run error %s", err.Error())
		}
	}
}
