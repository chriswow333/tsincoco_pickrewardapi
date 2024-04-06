package service

import (
	"context"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	cardStore "pickrewardapi/internal/domain/card/store"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type TaskLabelService interface {
	GetShowTaskLabels(ctx context.Context) ([]*cardDTO.TaskLabelDTO, error)
}

type taskLabelImpl struct {
	dig.In

	taskLabelStore cardStore.TaskLabelStore
}

func NewTaskLabel(
	taskLabelStore cardStore.TaskLabelStore,
) TaskLabelService {

	impl := &taskLabelImpl{
		taskLabelStore: taskLabelStore,
	}
	return impl
}

func (im *taskLabelImpl) GetShowTaskLabels(ctx context.Context) ([]*cardDTO.TaskLabelDTO, error) {
	logPos := "[task.label.service][GetShowTaskLabels]"

	labels, err := im.taskLabelStore.GetAllTaskLabels(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("taskLabelStore.GetAllTaskLabels failed", err)
		return nil, err
	}

	showLabels := []*cardDTO.TaskLabelDTO{}
	for _, l := range labels {
		if l.Show == 1 {
			showLabels = append(showLabels, l)
		}
	}
	return showLabels, nil
}
