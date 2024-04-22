package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	psql "pickrewardapi/internal/pkg/postgres"
)

type TaskLabelStore interface {
	GetAllTaskLabels(ctx context.Context) ([]*cardDTO.TaskLabelDTO, error)
	ModifiedTaskLabel(ctx context.Context, taskLabel *cardDTO.TaskLabelDTO) error
	GetTaskLabelByID(ctx context.Context, id string) (*cardDTO.TaskLabelDTO, error)
}

type taskLabelImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewTaskLabel(sql *psql.Psql) TaskLabelStore {
	logPos := "[task.label.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Infof("init card reward store")

	return &taskLabelImpl{
		primary: sql.Primary,
	}
}

const TASK_LABEL = "task_label"
const ALL_TASK_LABEL_COLUMNS = " \"label\", \"name\", \"show\" "

var MODIFIED_TASK_LABEL_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3)"+
		" ON CONFLICT(label) DO UPDATE SET  "+
		" \"name\" = $4, \"show\" = $5 ",
	TASK_LABEL, ALL_TASK_LABEL_COLUMNS,
)

func (im *taskLabelImpl) ModifiedTaskLabel(ctx context.Context, taskLabelDTO *cardDTO.TaskLabelDTO) error {
	logPos := "[task.label.store][ModifiedTaskLabel]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":                logPos,
			"taskLabelDTO.Label": taskLabelDTO.Label,
		}).Error("psql.Begin failed: ", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.WithFields(log.Fields{
					"pos": logPos,
				}).Error("tx.Rollback failed: ", err)
			}
		}
	}()

	updater := []interface{}{
		taskLabelDTO.Label,
		taskLabelDTO.Name,
		taskLabelDTO.Show,

		taskLabelDTO.Name,
		taskLabelDTO.Show,
	}

	if _, err := tx.Exec(MODIFIED_TASK_LABEL_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":                logPos,
			"taskLabelDTO.Label": taskLabelDTO.Label,
		}).Error("tx.Exec failed: ", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("tx.Commit failed: ", err)
		return err
	}

	return nil
}

var SELECT_TASK_LABELS_STAT = fmt.Sprintf(" SELECT %s FROM %s ",
	ALL_TASK_LABEL_COLUMNS, TASK_LABEL,
)

func (im *taskLabelImpl) GetAllTaskLabels(ctx context.Context) ([]*cardDTO.TaskLabelDTO, error) {
	logPos := "[task.label.store][GetAllTaskLabels]"

	taskLabelDTOs := []*cardDTO.TaskLabelDTO{}

	rows, err := im.primary.Query(SELECT_TASK_LABELS_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		taskLabelDTO := &cardDTO.TaskLabelDTO{}
		selector := []interface{}{
			&taskLabelDTO.Label,
			&taskLabelDTO.Name,
			&taskLabelDTO.Show,
		}
		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		taskLabelDTOs = append(taskLabelDTOs, taskLabelDTO)
	}

	return taskLabelDTOs, nil
}

var SELECT_TASK_LABEL_BY_ID_STAT = fmt.Sprintf(" SELECT %s FROM %s "+
	" WHERE id = $1 ",
	ALL_TASK_LABEL_COLUMNS, TASK_LABEL,
)

func (im *taskLabelImpl) GetTaskLabelByID(ctx context.Context, id string) (*cardDTO.TaskLabelDTO, error) {
	logPos := "[task.label.store][GetTaskLabelByID]"

	rows, err := im.primary.Query(SELECT_TASK_LABEL_BY_ID_STAT, id)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	taskLabelDTO := &cardDTO.TaskLabelDTO{}
	if rows.Next() {

		selector := []interface{}{
			&taskLabelDTO.Label,
			&taskLabelDTO.Name,
			&taskLabelDTO.Show,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		return taskLabelDTO, nil
	}
	return nil, errors.New("not found record")
}
