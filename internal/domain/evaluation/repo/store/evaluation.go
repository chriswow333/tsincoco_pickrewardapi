package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	evaluationDTO "pickrewardapi/internal/domain/evaluation/dto"
	psql "pickrewardapi/internal/pkg/postgres"
)

type EvaluationStore interface {
	GetAllEvaluations(ctx context.Context) ([]*evaluationDTO.EvaluationDTO, error)
	GetEvaluationByID(ctx context.Context, ID string) (*evaluationDTO.EvaluationDTO, error)
	GetEvaluationByOwnerID(ctx context.Context, ownerID string) (*evaluationDTO.EvaluationDTO, error)
	ModifiedEvaluation(ctx context.Context, dto *evaluationDTO.EvaluationDTO) error
	DeleteEvaluationByID(ctx context.Context, ownerID string) error
}

type evaluationImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewEvaluation(sql *psql.Psql) EvaluationStore {
	logPos := "[evaluation.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")

	return &evaluationImpl{
		primary: sql.Primary,
	}
}

const EVALUATION = "evaluation"
const ALL_COLUMNS = " \"id\", \"feedback_id\", \"owner_id\", start_date, end_date, create_date, update_date, payload "

var SELECT_ALL_EVALUATIONS_BY_STAT = fmt.Sprintf(
	" SELECT %s FROM %s ",
	ALL_COLUMNS, EVALUATION,
)

func (im *evaluationImpl) GetAllEvaluations(ctx context.Context) ([]*evaluationDTO.EvaluationDTO, error) {
	logPos := "[evaluation.store][GetAllEvaluations]"

	evaluationDTOs := []*evaluationDTO.EvaluationDTO{}

	rows, err := im.primary.Query(SELECT_ALL_EVALUATIONS_BY_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Errorf("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		evaluationDTO := &evaluationDTO.EvaluationDTO{}
		selector := []interface{}{
			&evaluationDTO.ID,
			&evaluationDTO.FeedbackID,
			&evaluationDTO.OwnerID,
			&evaluationDTO.StartDate,
			&evaluationDTO.EndDate,
			&evaluationDTO.CreateDate,
			&evaluationDTO.UpdateDate,
			&evaluationDTO.Payload,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		evaluationDTOs = append(evaluationDTOs, evaluationDTO)
	}

	return evaluationDTOs, nil
}

var SELECT_EVALUATION_BY_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"id\" = $1 ",
	ALL_COLUMNS, EVALUATION,
)

func (im *evaluationImpl) GetEvaluationByID(ctx context.Context, ID string) (*evaluationDTO.EvaluationDTO, error) {
	logPos := "[evaluation.store][GetEvaluationByID]"

	rows, err := im.primary.Query(SELECT_EVALUATION_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	evaluationDTO := &evaluationDTO.EvaluationDTO{}
	for rows.Next() {
		selector := []interface{}{
			&evaluationDTO.ID,
			&evaluationDTO.FeedbackID,
			&evaluationDTO.OwnerID,
			&evaluationDTO.StartDate,
			&evaluationDTO.EndDate,
			&evaluationDTO.CreateDate,
			&evaluationDTO.UpdateDate,
			&evaluationDTO.Payload,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

	}
	return evaluationDTO, nil
}

var SELECT_EVALUATION_BY_OWNER_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"owner_id\" = $1 ",
	ALL_COLUMNS, EVALUATION,
)

func (im *evaluationImpl) GetEvaluationByOwnerID(ctx context.Context, ownerID string) (*evaluationDTO.EvaluationDTO, error) {
	logPos := "[evaluation.store][GetEvaluationByOwnerID]"

	rows, err := im.primary.Query(SELECT_EVALUATION_BY_OWNER_ID_STAT, ownerID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	evaluationDTO := &evaluationDTO.EvaluationDTO{}
	for rows.Next() {
		selector := []interface{}{
			&evaluationDTO.ID,
			&evaluationDTO.FeedbackID,
			&evaluationDTO.OwnerID,
			&evaluationDTO.StartDate,
			&evaluationDTO.EndDate,
			&evaluationDTO.CreateDate,
			&evaluationDTO.UpdateDate,
			&evaluationDTO.Payload,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

	}
	return evaluationDTO, nil
}

var MODIFIED_EVALUATION_STAT = fmt.Sprintf(
	" INSERT INTO %s "+
		" (%s) VALUES "+
		" ($1, $2, $3, $4, $5, $6, $7, $8) "+
		" ON CONFLICT(id) DO UPDATE SET "+
		" feedback_id = $9, owner_id = $10, "+
		" start_date = $11, end_date = $12, create_date = $13, update_date = $14, "+
		" payload = $15 ",
	EVALUATION, ALL_COLUMNS,
)

func (im *evaluationImpl) ModifiedEvaluation(ctx context.Context, evaluationDTO *evaluationDTO.EvaluationDTO) error {
	logPos := "[evaluation.store][ModifiedEvaluation]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
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

	const ALL_COLUMNS = " \"id\", \"feedback_id\", \"owner_id\", start_date, end_date, create_date, update_date, payload "

	updater := []interface{}{
		evaluationDTO.ID,
		evaluationDTO.FeedbackID,
		evaluationDTO.OwnerID,
		evaluationDTO.StartDate,
		evaluationDTO.EndDate,
		evaluationDTO.CreateDate,
		evaluationDTO.UpdateDate,
		evaluationDTO.Payload,

		evaluationDTO.FeedbackID,
		evaluationDTO.OwnerID,
		evaluationDTO.StartDate,
		evaluationDTO.EndDate,
		evaluationDTO.CreateDate,
		evaluationDTO.UpdateDate,
		evaluationDTO.Payload,
	}

	if _, err := tx.Exec(MODIFIED_EVALUATION_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
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

var DELETE_EVALUATION_BY_ID_STAT = fmt.Sprintf(
	"DELETE FROM %s WHERE \"id\" = $1 ",
	EVALUATION,
)

func (im *evaluationImpl) DeleteEvaluationByID(ctx context.Context, id string) error {
	logPos := "[evaluation.store][DeleteEvaluationByID]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
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
		id,
	}

	if _, err := tx.Exec(DELETE_EVALUATION_BY_ID_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("tx.Exec failed: ", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("tx.Commit failed: ", err)
		return err
	}
	return nil
}
