package store

import (
	"context"
	"errors"
	"fmt"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	commonM "pickrewardapi/internal/shared/common/model"

	psql "pickrewardapi/internal/pkg/postgres"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type BankStore interface {
	ModifiedBank(ctx context.Context, bankDTO *cardDTO.BankDTO) error
	GetBankByID(ctx context.Context, ID string) (*cardDTO.BankDTO, error)
	GetAllBanks(ctx context.Context, status commonM.Status) ([]*cardDTO.BankDTO, error)
	GetBankNameByBankID(ctx context.Context, ID string) (*cardDTO.BankDTO, error)
}

type bankImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewBank(sql *psql.Psql) BankStore {
	logPos := "[bank.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Infof("init bank store")

	return &bankImpl{
		primary: sql.Primary,
	}
}

const BANK = "bank"
const ALL_COLUMNS = " \"id\", \"name\", \"order\", \"bank_status\", \"create_date\", \"update_date\" "

var MODIFIED_BANK_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5, $6) "+
		" ON CONFLICT(id) DO UPDATE SET  "+
		" \"name\" = $7, \"order\" = $8, "+
		" \"bank_status\" = $9, \"create_date\" = $10 , \"update_date\" = $11 ",
	BANK, ALL_COLUMNS,
)

func (im *bankImpl) ModifiedBank(ctx context.Context, bankDTO *cardDTO.BankDTO) error {
	logPos := "[card.store][ModifiedBank]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": bankDTO.ID,
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
		bankDTO.ID,
		bankDTO.Name,
		bankDTO.Order,
		bankDTO.BankStatus,
		bankDTO.CreateDate,
		bankDTO.UpdateDate,

		bankDTO.Name,
		bankDTO.Order,
		bankDTO.BankStatus,
		bankDTO.CreateDate,
		bankDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_BANK_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": bankDTO.ID,
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

var SELECT_BANK_BY_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"id\" = $1 ",
	ALL_COLUMNS, BANK,
)

func (im *bankImpl) GetBankByID(ctx context.Context, ID string) (*cardDTO.BankDTO, error) {
	logPos := "[bank.store][GetBankByID]"

	rows, err := im.primary.Query(SELECT_BANK_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		bankResult := &cardDTO.BankDTO{}
		selector := []interface{}{
			&bankResult.ID,
			&bankResult.Name,
			&bankResult.Order,
			&bankResult.BankStatus,
			&bankResult.CreateDate,
			&bankResult.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"bank.ID": ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		return bankResult, nil
	}

	return nil, errors.New("not found bank")

}

var SELECT_ALL_BANK_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE bank_status = $1 "+
		" ORDER BY \"order\" ",
	ALL_COLUMNS, BANK,
)

func (im *bankImpl) GetAllBanks(ctx context.Context, status commonM.Status) ([]*cardDTO.BankDTO, error) {
	logPos := "[bank.store][GetAllBanks]"

	bankDTOs := []*cardDTO.BankDTO{}
	rows, err := im.primary.Query(SELECT_ALL_BANK_STAT, status)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		bankDTO := &cardDTO.BankDTO{}
		selector := []interface{}{
			&bankDTO.ID,
			&bankDTO.Name,
			&bankDTO.Order,
			&bankDTO.CreateDate,
			&bankDTO.UpdateDate,
			&bankDTO.BankStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		bankDTOs = append(bankDTOs, bankDTO)
	}

	return bankDTOs, nil
}

var SELECT_BANK_NAME_BY_ID_STAT = fmt.Sprintf(
	"SELECT \"id\", \"name\" FROM %s "+
		" WHERE bank_status = $1 "+
		" ORDER BY \"order\" ",
	BANK,
)

func (im *bankImpl) GetBankNameByBankID(ctx context.Context, ID string) (*cardDTO.BankDTO, error) {
	logPos := "[bank.store][GetBankNameByBankID]"

	rows, err := im.primary.Query(SELECT_BANK_NAME_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	var bankResult *cardDTO.BankDTO
	defer rows.Close()

	if rows.Next() {

		bankResult := &cardDTO.BankDTO{}
		selector := []interface{}{
			&bankResult.ID,
			&bankResult.Name,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"bank.ID": ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		if rows.Next() {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"bank.ID": ID,
			}).Error("there have more than one record")
			return nil, errors.New("there have more than one record")
		}

	} else {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": ID,
		}).Info("no bank found with the given ID")
		return nil, errors.New("no bank found with the given ID")
	}

	return bankResult, nil
}
