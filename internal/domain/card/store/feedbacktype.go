package store

import (
	"context"
	"errors"
	"fmt"
	cardDTO "pickrewardapi/internal/domain/card/dto"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	psql "pickrewardapi/internal/pkg/postgres"
)

type FeedbackTypeStore interface {
	ModifiedFeedbackType(ctx context.Context, rewardDTO *cardDTO.FeedbackTypeDTO) error
	GetAllFeedbackTypes(ctx context.Context) ([]*cardDTO.FeedbackTypeDTO, error)
	GetFeedbackTypeByID(ctx context.Context, ID string) (*cardDTO.FeedbackTypeDTO, error)
}

type feedbackTypeImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewFeedbackType(sql *psql.Psql) FeedbackTypeStore {
	logPos := "[feedback.type.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init feedback type store")

	return &feedbackTypeImpl{
		primary: sql.Primary,
	}
}

const FEEDBACK_TYPE = "feedback_type"
const ALL_FEEDBACK_TYPE_COLUMNS = " \"id\", \"name\", \"feedback_type\", " +
	" \"create_date\", \"update_date\" "

var MODIFIED_FEEDBACK_TYPE_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5) "+
		" ON CONFLICT(id) DO UPDATE SET "+
		" \"name\" = $6, \"feedback_type\" = $7, \"create_date\" = $8, \"update_date\" = $9 ",
	FEEDBACK_TYPE, ALL_FEEDBACK_TYPE_COLUMNS,
)

func (im *feedbackTypeImpl) ModifiedFeedbackType(ctx context.Context, feedbackDTO *cardDTO.FeedbackTypeDTO) error {
	logPos := "[reward.store][ModifiedFeedbackType]"

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
		feedbackDTO.ID,
		feedbackDTO.Name,
		feedbackDTO.FeedbackType,
		feedbackDTO.CreateDate,
		feedbackDTO.UpdateDate,

		feedbackDTO.Name,
		feedbackDTO.FeedbackType,
		feedbackDTO.CreateDate,
		feedbackDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_FEEDBACK_TYPE_STAT, updater...); err != nil {
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

var SELECT_ALL_FEEDBACK_TYPES_STAT = fmt.Sprintf(
	"SELECT %s FROM %s ",
	ALL_FEEDBACK_TYPE_COLUMNS, FEEDBACK_TYPE,
)

func (im *feedbackTypeImpl) GetAllFeedbackTypes(ctx context.Context) ([]*cardDTO.FeedbackTypeDTO, error) {
	logPos := "[feedback.type.store][GetAllFeedbackTypes]"

	feedbackDTOs := []*cardDTO.FeedbackTypeDTO{}

	rows, err := im.primary.Query(SELECT_ALL_FEEDBACK_TYPES_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		feedbackDTO := &cardDTO.FeedbackTypeDTO{}
		selector := []interface{}{
			&feedbackDTO.ID,
			&feedbackDTO.Name,
			&feedbackDTO.FeedbackType,
			&feedbackDTO.CreateDate,
			&feedbackDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		feedbackDTOs = append(feedbackDTOs, feedbackDTO)
	}

	return feedbackDTOs, nil
}

var SELECT_FEEDBACK_TYPE_BY_ID_STAT = fmt.Sprintf("SELECT %s FROM %s "+
	" WHERE \"id\" = $1 ",
	ALL_FEEDBACK_TYPE_COLUMNS, FEEDBACK_TYPE,
)

func (im *feedbackTypeImpl) GetFeedbackTypeByID(ctx context.Context, ID string) (*cardDTO.FeedbackTypeDTO, error) {
	logPos := "[feedback.type.store][GetFeedbackTypeByID]"

	var r *cardDTO.FeedbackTypeDTO

	rows, err := im.primary.Query(SELECT_FEEDBACK_TYPE_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		r = &cardDTO.FeedbackTypeDTO{}
		selector := []interface{}{
			&r.ID,
			&r.Name,
			&r.FeedbackType,
			&r.CreateDate,
			&r.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		if rows.Next() {
			log.WithFields(log.Fields{
				"pos": logPos,
				"ID":  ID,
			}).Error("has more than one record")
			return nil, errors.New("has more than one record")
		}

	}
	return r, nil
}
