package pipeline

import "github.com/leilei3167/copy_design_pattern/monitor/plugin"

type SimplePipeline struct {
	pipelineTemplate
}

func (s *SimplePipeline) SetContext(ctx plugin.Context) {
	s.run = func() {
		go func() {
			s.doRun()
		}()
	}
}
