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
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardRewardStore interface {
	GetCardRewardsByCardID(ctx context.Context, cardID string, status commonM.Status) ([]*cardDTO.CardRewardDTO, error)
	GetAllCardRewardsByCardRewardType(ctx context.Context, status commonM.Status, cardRewardType cardDTO.CardRewardType, time int64) ([]*cardDTO.CardRewardDTO, error)

	ModifiedCardReward(ctx context.Context, dto *cardDTO.CardRewardDTO) error
	DeleteCardReward(ctx context.Context, ID string) error
	GetAllCardRewardsByCardIDAndCardRewardType(ctx context.Context, cardID string, cardRewardType cardDTO.CardRewardType) ([]*cardDTO.CardRewardDTO, error)

	GetCardRewardByID(ctx context.Context, ID string) (*cardDTO.CardRewardDTO, error)
}

type cardRewardImpl struct {
	dig.In

	primary   *pgx.ConnPool
	migration *pgx.ConnPool
}

func NewCardReward(sql *psql.Psql) CardRewardStore {
	logPos := "[card.reward.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Infof("init card reward store")

	return &cardRewardImpl{
		primary:   sql.Primary,
		migration: sql.Migration,
	}
}

const CARD_REWARD = "card_reward"
const ALL_CARD_REWARD_COLUMNS = " \"id\", \"card_id\", \"name\", \"description\", " +
	" \"start_date\", \"end_date\", \"card_reward_type\", \"feedback_type_id\", \"task_labels\", " +
	" \"order\", \"card_reward_status\", \"create_date\", \"update_date\" "

var MODIFIED_CARD_REWARD_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) "+
		" ON CONFLICT(id) DO UPDATE SET  "+
		" card_id = $14, name = $15, description = $16, start_date = $17, end_date = $18, "+
		" card_reward_type = $19, feedback_type_id = $20, \"task_labels\" = $21,  "+
		" \"order\" = $22, card_reward_status = $23, create_date = $24, update_date = $25 ",
	CARD_REWARD, ALL_CARD_REWARD_COLUMNS,
)

func (im *cardRewardImpl) ModifiedCardReward(ctx context.Context, cardRewardDTO *cardDTO.CardRewardDTO) error {

	logPos := "[card.reward.store][ModifiedCardReward]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":            logPos,
			"card.reward.ID": cardRewardDTO.ID,
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
		cardRewardDTO.ID,
		cardRewardDTO.CardID,
		cardRewardDTO.Name,
		cardRewardDTO.Descriptions,
		cardRewardDTO.StartDate,
		cardRewardDTO.EndDate,
		cardRewardDTO.CardRewardType,
		cardRewardDTO.FeedbackType.ID,
		cardRewardDTO.TaskLabelDTOs,
		cardRewardDTO.Order,
		cardRewardDTO.CardRewardStatus,
		cardRewardDTO.CreateDate,
		cardRewardDTO.UpdateDate,

		cardRewardDTO.CardID,
		cardRewardDTO.Name,
		cardRewardDTO.Descriptions,
		cardRewardDTO.StartDate,
		cardRewardDTO.EndDate,
		cardRewardDTO.CardRewardType,
		cardRewardDTO.FeedbackType.ID,
		cardRewardDTO.TaskLabelDTOs,
		cardRewardDTO.Order,
		cardRewardDTO.CardRewardStatus,
		cardRewardDTO.CreateDate,
		cardRewardDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_CARD_REWARD_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":            logPos,
			"card.reward.ID": cardRewardDTO.ID,
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

var SELECT_CARD_REWARDS_BY_CARD_REWARD_TYPE_STAT = fmt.Sprintf(
	" SELECT %s FROM %s "+
		" WHERE card_reward_status = $1 "+
		" AND card_reward_type = $2 "+
		" AND start_date <= $3 "+
		" AND end_date >= $4 "+
		" ORDER BY \"order\" ",
	ALL_CARD_REWARD_COLUMNS, CARD_REWARD,
)

func (im *cardRewardImpl) GetAllCardRewardsByCardRewardType(ctx context.Context, status commonM.Status, cardRewardType cardDTO.CardRewardType, time int64) ([]*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.store][GetAllCardRewardsByCardRewardType]"

	rows, err := im.primary.Query(SELECT_CARD_REWARDS_BY_CARD_REWARD_TYPE_STAT, status, cardRewardType, time, time)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	dtos := []*cardDTO.CardRewardDTO{}

	for rows.Next() {

		dto := &cardDTO.CardRewardDTO{}
		selector := []interface{}{
			&dto.ID,
			&dto.CardID,
			&dto.Name,
			&dto.Descriptions,
			&dto.StartDate,
			&dto.EndDate,
			&dto.CardRewardType,
			&dto.FeedbackType,
			&dto.TaskLabelDTOs,
			&dto.Order,
			&dto.CardRewardStatus,
			&dto.CreateDate,
			&dto.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}

var SELECT_CARD_REWARD_BY_CARD_ID_STAT = fmt.Sprintf(
	" SELECT %s FROM %s "+
		" WHERE \"card_id\" = $1 "+
		" AND card_reward_status = $2 "+
		" ORDER BY \"order\" ",
	ALL_CARD_REWARD_COLUMNS, CARD_REWARD,
)

func (im *cardRewardImpl) GetCardRewardsByCardID(ctx context.Context, cardID string, status commonM.Status) ([]*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.store][GetCardRewardsByCardID]"

	rows, err := im.primary.Query(SELECT_CARD_REWARD_BY_CARD_ID_STAT, cardID, status)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	dtos := []*cardDTO.CardRewardDTO{}

	for rows.Next() {

		dto := &cardDTO.CardRewardDTO{
			FeedbackType: &cardDTO.FeedbackTypeDTO{},
		}
		selector := []interface{}{
			&dto.ID,
			&dto.CardID,
			&dto.Name,
			&dto.Descriptions,
			&dto.StartDate,
			&dto.EndDate,
			&dto.CardRewardType,
			&dto.FeedbackType.ID,
			&dto.TaskLabelDTOs,
			&dto.Order,
			&dto.CardRewardStatus,
			&dto.CreateDate,
			&dto.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"card.ID": cardID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

var SELECT_ALL_CARD_REWARD_BY_CARD_ID_AND_CARD_REWARD_TYPE_STAT = fmt.Sprintf(
	"SELECT %s FROM %s"+
		" WHERE \"card_id\" = $1 "+
		" AND card_reward_type = $2 "+
		" ORDER BY \"order\" ",
	ALL_CARD_REWARD_COLUMNS, CARD_REWARD,
)

func (im *cardRewardImpl) GetAllCardRewardsByCardIDAndCardRewardType(ctx context.Context, cardID string, cardRewardType cardDTO.CardRewardType) ([]*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.store][GetAllCardRewardsByCardIDAndCardRewardType]"

	rows, err := im.primary.Query(SELECT_ALL_CARD_REWARD_BY_CARD_ID_AND_CARD_REWARD_TYPE_STAT, cardID, cardRewardType)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	dtos := []*cardDTO.CardRewardDTO{}

	for rows.Next() {
		dto := &cardDTO.CardRewardDTO{
			FeedbackType: &cardDTO.FeedbackTypeDTO{},
		}
		selector := []interface{}{
			&dto.ID,
			&dto.CardID,
			&dto.Name,
			&dto.Descriptions,
			&dto.StartDate,
			&dto.EndDate,
			&dto.CardRewardType,
			&dto.FeedbackType.ID,
			&dto.TaskLabelDTOs,
			&dto.Order,
			&dto.CardRewardStatus,
			&dto.CreateDate,
			&dto.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"card.ID": cardID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		dtos = append(dtos, dto)
	}
	return dtos, nil
}

var DELETE_CARD_REWARD_BY_ID_STAT = fmt.Sprintf(
	"DELETE FROM %s WHERE \"id\" = $1 ",
	CARD_REWARD,
)

func (im *cardRewardImpl) DeleteCardReward(ctx context.Context, ID string) error {
	logPos := "[card.reward.store][DeleteCardReward]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": ID,
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
		ID,
	}

	if _, err := tx.Exec(DELETE_CARD_REWARD_BY_ID_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": ID,
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

var SELECT_CARD_REWARD_BY_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"id\" = $1 ",
	ALL_CARD_COLUMNS, CARD_REWARD,
)

func (im *cardRewardImpl) GetCardRewardByID(ctx context.Context, ID string) (*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.store][GetCardRewardByID]"

	rows, err := im.primary.Query(SELECT_CARD_REWARD_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
			"ID":  ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		dto := &cardDTO.CardRewardDTO{
			FeedbackType: &cardDTO.FeedbackTypeDTO{},
		}
		selector := []interface{}{
			&dto.ID,
			&dto.CardID,
			&dto.Name,
			&dto.Descriptions,
			&dto.StartDate,
			&dto.EndDate,
			&dto.CardRewardType,
			&dto.FeedbackType.ID,
			&dto.TaskLabelDTOs,
			&dto.Order,
			&dto.CardRewardStatus,
			&dto.CreateDate,
			&dto.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
				"ID":  ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		return dto, nil
	}

	return nil, errors.New("not found card reward")

}
