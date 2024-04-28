package store

import (
	"context"
	"errors"
	"fmt"

	commonM "pickrewardapi/internal/shared/common/model"

	payDTO "pickrewardapi/internal/domain/pay/dto"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	psql "pickrewardapi/internal/pkg/postgres"
)

type PayStore interface {
	ModifiedPay(ctx context.Context, payDTO *payDTO.PayDTO) error
	GetPayByID(ctx context.Context, ID string) (*payDTO.PayDTO, error)
	GetAllPays(ctx context.Context, status commonM.Status) ([]*payDTO.PayDTO, error)
}

type payImpl struct {
	dig.In

	primary   *pgx.ConnPool
	migration *pgx.ConnPool
}

func NewPay(sql *psql.Psql) PayStore {
	logPos := "[pay.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Infof("init bank store")

	return &payImpl{
		primary:   sql.Primary,
		migration: sql.Migration,
	}
}

const PAY = "pay"
const ALL_COLUMNS = " \"id\", \"name\", \"order\", \"pay_status\", \"create_date\", \"update_date\" "

var MODIFIED_PAY_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5, $6) "+
		" ON CONFLICT(id) DO UPDATE SET  "+
		" \"name\" = $7, \"order\" = $8, "+
		" \"pay_status\" = $9, \"create_date\" = $10 , \"update_date\" = $11 ",
	PAY, ALL_COLUMNS,
)

func (im *payImpl) ModifiedPay(ctx context.Context, payDTO *payDTO.PayDTO) error {
	logPos := "[pay.store][ModifiedPay]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":    logPos,
			"pay.ID": payDTO.ID,
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
		payDTO.ID,
		payDTO.Name,
		payDTO.Order,
		payDTO.PayStatus,
		payDTO.CreateDate,
		payDTO.UpdateDate,

		payDTO.Name,
		payDTO.Order,
		payDTO.PayStatus,
		payDTO.CreateDate,
		payDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_PAY_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":    logPos,
			"pay.ID": payDTO.ID,
		}).Errorf("tx.Exec failed: ", err)
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

var SELECT_PAY_BY_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"id\" = $1 ",
	ALL_COLUMNS, PAY,
)

func (im *payImpl) GetPayByID(ctx context.Context, ID string) (*payDTO.PayDTO, error) {
	logPos := "[pay.store][GetPayByID]"

	rows, err := im.primary.Query(SELECT_PAY_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		payResult := &payDTO.PayDTO{}
		selector := []interface{}{
			&payResult.ID,
			&payResult.Name,
			&payResult.Order,
			&payResult.PayStatus,
			&payResult.CreateDate,
			&payResult.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"bank.ID": ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		return payResult, nil
	}

	return nil, errors.New("not found pay")

}

var SELECT_ALL_PAY_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE pay_status = $1 "+
		" ORDER BY \"order\" ",
	ALL_COLUMNS, PAY,
)

func (im *payImpl) GetAllPays(ctx context.Context, status commonM.Status) ([]*payDTO.PayDTO, error) {
	logPos := "[pay.store][GetAllPays]"

	payDTOs := []*payDTO.PayDTO{}
	rows, err := im.primary.Query(SELECT_ALL_PAY_STAT, status)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		payDTO := &payDTO.PayDTO{}
		selector := []interface{}{
			&payDTO.ID,
			&payDTO.Name,
			&payDTO.Order,
			&payDTO.CreateDate,
			&payDTO.UpdateDate,
			&payDTO.PayStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		payDTOs = append(payDTOs, payDTO)
	}

	return payDTOs, nil
}
