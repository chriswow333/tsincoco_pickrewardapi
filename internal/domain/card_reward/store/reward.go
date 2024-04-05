package store

import (
	"context"

	rewardDTO "pickrewardapi/internal/domain/card_reward/dto"

	"errors"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	psql "pickrewardapi/internal/pkg/postgres"
)

type RewardStore interface {
	// CreateReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error
	// UpdateReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error

	ModifiedReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error

	GetAllRewards(ctx context.Context) ([]*rewardDTO.RewardDTO, error)
	GetRewardByID(ctx context.Context, ID string) (*rewardDTO.RewardDTO, error)
	MigrateReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error
}

type impl struct {
	dig.In

	primary   *pgx.ConnPool
	migration *pgx.ConnPool
}

func New(sql *psql.Psql) RewardStore {
	logPos := "[reward.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init reward store")

	return &impl{
		primary:   sql.Primary,
		migration: sql.Migration,
	}
}

const MODIFIED_REWARD_STAT = "INSERT INTO reward " +
	"(\"id\", \"name\", \"reward_type\", \"create_date\", \"update_date\") " +
	" VALUES ($1, $2, $3, $4, $5) " +
	" ON CONFLICT(id) DO UPDATE SET " +
	" \"name\" = $6, \"reward_type\" = $7, \"create_date\" = $8, \"update_date\" = $9 "

func (im *impl) ModifiedReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error {
	logPos := "[reward.store][CreateReward]"

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
		rewardDTO.ID,
		rewardDTO.Name,
		rewardDTO.RewardType,
		rewardDTO.CreateDate,
		rewardDTO.UpdateDate,

		rewardDTO.Name,
		rewardDTO.RewardType,
		rewardDTO.CreateDate,
		rewardDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_REWARD_STAT, updater...); err != nil {
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

const SELECT_ALL_REWARDS_STAT = "SELECT \"id\", \"name\", \"reward_type\", create_date, update_date " +
	" FROM reward "

func (im *impl) GetAllRewards(ctx context.Context) ([]*rewardDTO.RewardDTO, error) {
	logPos := "[reward.store][GetAllRewards]"

	rewardDTOs := []*rewardDTO.RewardDTO{}

	rows, err := im.primary.Query(SELECT_ALL_REWARDS_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		rewardDTO := &rewardDTO.RewardDTO{}
		selector := []interface{}{
			&rewardDTO.ID,
			&rewardDTO.Name,
			&rewardDTO.RewardType,
			&rewardDTO.CreateDate,
			&rewardDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		rewardDTOs = append(rewardDTOs, rewardDTO)
	}

	return rewardDTOs, nil
}

const SELECT_REWARD_BY_ID_STAT = "SELECT \"id\", \"name\", \"reward_type\", create_date, update_date " +
	" FROM reward " +
	" WHERE \"id\" = $1 "

func (im *impl) GetRewardByID(ctx context.Context, ID string) (*rewardDTO.RewardDTO, error) {
	logPos := "[point.store][GetPointByID]"

	var r *rewardDTO.RewardDTO

	rows, err := im.primary.Query(SELECT_REWARD_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		r = &rewardDTO.RewardDTO{}
		selector := []interface{}{
			&r.ID,
			&r.Name,
			&r.RewardType,
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
			}).Error("There have more than one record.")
			return nil, errors.New("There have more than one record.")
		}

	}
	return r, nil
}

const MIGRATION_REWARD_STAT = "INSERT INTO reward " +
	"(\"id\", \"name\", \"reward_type\", \"create_date\", \"update_date\") " +
	" VALUES ($1, $2, $3, $4, $5) " +
	" ON CONFLICT(id) DO UPDATE SET " +
	" \"name\" = $6, \"reward_type\" = $7, \"create_date\" = $8, \"update_date\" = $9 "

func (im *impl) MigrateReward(ctx context.Context, rewardDTO *rewardDTO.RewardDTO) error {

	logPos := "[reward.store][MigrateReward]"

	tx, err := im.migration.Begin()
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
		rewardDTO.ID,
		rewardDTO.Name,
		rewardDTO.RewardType,
		rewardDTO.CreateDate,
		rewardDTO.UpdateDate,

		rewardDTO.Name,
		rewardDTO.RewardType,
		rewardDTO.CreateDate,
		rewardDTO.UpdateDate,
	}

	if _, err := tx.Exec(MIGRATION_REWARD_STAT, updater...); err != nil {
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
