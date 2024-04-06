package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	psql "pickrewardapi/internal/pkg/postgres"
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardStore interface {
	ModifiedCard(ctx context.Context, cardDTO *cardDTO.CardDTO) error

	GetByCardID(ctx context.Context, ID string) (*cardDTO.CardDTO, error)
	GetLatestCards(ctx context.Context) ([]*cardDTO.CardDTO, error)
	GetAllCards(ctx context.Context) ([]*cardDTO.CardDTO, error)
	GetCardsByBankID(ctx context.Context, bankID string, status commonM.Status) ([]*cardDTO.CardDTO, error)
	SearchCard(ctx context.Context, keyword string, status commonM.Status) ([]*cardDTO.CardDTO, error)
}

type cardImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewCard(sql *psql.Psql) CardStore {
	logPos := "[card.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Infof("init card store")

	return &cardImpl{
		primary: sql.Primary,
	}
}

const CARD = "card"
const ALL_CARD_COLUMNS = " \"id\", \"name\", \"descriptions\",\"link_url\", " +
	" \"bank_id\", \"order\", \"card_status\", \"create_date\", \"update_date\" "

var MODIFIED_CARD_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"+
		" ON CONFLICT(id) DO UPDATE SET  "+
		" \"name\" = $10, \"descriptions\" = $11, \"link_url\" = $12, \"bank_id\" = $13, \"order\" = $14, "+
		" \"card_status\" = $15 ,\"create_date\" = $16, \"update_date\" = $17 ",
	CARD, ALL_CARD_COLUMNS,
)

func (im *cardImpl) ModifiedCard(ctx context.Context, cardDTO *cardDTO.CardDTO) error {
	logPos := "[card.store][ModifiedCard]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardDTO.ID,
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
		cardDTO.ID,
		cardDTO.Name,
		cardDTO.Descriptions,
		cardDTO.LinkURL,
		cardDTO.BankID,
		cardDTO.Order,
		cardDTO.CardStatus,
		cardDTO.CreateDate,
		cardDTO.UpdateDate,

		cardDTO.Name,
		cardDTO.Descriptions,
		cardDTO.LinkURL,
		cardDTO.BankID,
		cardDTO.Order,
		cardDTO.CardStatus,
		cardDTO.CreateDate,
		cardDTO.UpdateDate,
	}

	if _, err := tx.Exec(MODIFIED_CARD_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardDTO.ID,
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

var INSERT_CARD_STAT = fmt.Sprintf(
	"INSERT INTO %s (%s) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
	CARD, ALL_CARD_COLUMNS,
)

func (im *cardImpl) CreateCard(ctx context.Context, cardDTO *cardDTO.CardDTO) error {
	logPos := "[card.store][CreateCard]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardDTO.ID,
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
		cardDTO.ID,
		cardDTO.Name,
		cardDTO.Descriptions,
		cardDTO.LinkURL,
		cardDTO.BankID,
		cardDTO.Order,
		cardDTO.CardStatus,
		cardDTO.CreateDate,
		cardDTO.UpdateDate,
	}

	if _, err := tx.Exec(INSERT_CARD_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": cardDTO.ID,
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

var SELECT_CARD_BY_ID_STAT = fmt.Sprintf(
	"SELECT %s "+
		" FROM %s WHERE \"id\" = $1",
	ALL_CARD_COLUMNS, CARD,
)

func (im *cardImpl) GetByCardID(ctx context.Context, ID string) (*cardDTO.CardDTO, error) {
	logPos := "[card.store][GetByCardID]"

	var c *cardDTO.CardDTO

	rows, err := im.primary.Query(SELECT_CARD_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.ID": ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		c = &cardDTO.CardDTO{}
		selector := []interface{}{
			&c.ID,
			&c.Name,
			&c.Descriptions,
			&c.LinkURL,
			&c.BankID,
			&c.Order,
			&c.CardStatus,
			&c.CreateDate,
			&c.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"card.ID": ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		if rows.Next() {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"card.ID": ID,
			}).Error("There have more than one record.")
			return nil, errors.New("There have more than one record.")
		}

	}

	return c, nil
}

var SELECT_CARDS_BY_BANK_ID_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE \"bank_id\" = $1 "+
		" AND \"card_status\" = $2 ",
	ALL_CARD_COLUMNS, CARD,
)

func (im *cardImpl) GetCardsByBankID(ctx context.Context, bankID string, status commonM.Status) ([]*cardDTO.CardDTO, error) {
	logPos := "[card.store][GetCardsByBankID]"

	cardDTOs := []*cardDTO.CardDTO{}

	rows, err := im.primary.Query(SELECT_CARDS_BY_BANK_ID_STAT, bankID, status)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"bank.ID": bankID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		cardDTO := &cardDTO.CardDTO{}
		selector := []interface{}{
			&cardDTO.ID,
			&cardDTO.Name,
			&cardDTO.Descriptions,
			&cardDTO.LinkURL,
			&cardDTO.BankID,
			&cardDTO.Order,
			&cardDTO.CardStatus,
			&cardDTO.CreateDate,
			&cardDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":     logPos,
				"bank.ID": bankID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		cardDTOs = append(cardDTOs, cardDTO)
	}

	return cardDTOs, nil
}

var SELECT_ALL_CARDS_STAT = fmt.Sprintf(
	"SELECT %s FROM %s ",
	ALL_CARD_COLUMNS, CARD,
)

func (im *cardImpl) GetAllCards(ctx context.Context) ([]*cardDTO.CardDTO, error) {

	logPos := "[card.store][GetAllCards]"

	cardDTOs := []*cardDTO.CardDTO{}

	rows, err := im.primary.Query(SELECT_ALL_CARDS_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		cardDTO := &cardDTO.CardDTO{}
		selector := []interface{}{
			&cardDTO.ID,
			&cardDTO.Name,
			&cardDTO.Descriptions,
			&cardDTO.LinkURL,
			&cardDTO.BankID,
			&cardDTO.Order,
			&cardDTO.CardStatus,
			&cardDTO.CreateDate,
			&cardDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		cardDTOs = append(cardDTOs, cardDTO)
	}

	return cardDTOs, nil
}

var SELECT_LATEST_CARDS_STAT = fmt.Sprintf(
	"SELECT %s FROM %s "+
		" WHERE card_status = 2 order by update_date desc LIMIT 20 ",
	ALL_CARD_COLUMNS, CARD,
)

func (im *cardImpl) GetLatestCards(ctx context.Context) ([]*cardDTO.CardDTO, error) {

	logPos := "[card.store][GetLatestCards]"

	cardDTOs := []*cardDTO.CardDTO{}

	rows, err := im.primary.Query(SELECT_LATEST_CARDS_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		cardDTO := &cardDTO.CardDTO{}
		selector := []interface{}{
			&cardDTO.ID,
			&cardDTO.Name,
			&cardDTO.Descriptions,
			&cardDTO.LinkURL,
			&cardDTO.BankID,
			&cardDTO.Order,
			&cardDTO.CardStatus,
			&cardDTO.CreateDate,
			&cardDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		cardDTOs = append(cardDTOs, cardDTO)
	}

	return cardDTOs, nil
}

var SELECT_CARDS_BY_KEYWORD_STAT = fmt.Sprintf(
	"SELECT DISTINCT on(id) %s "+
		" FROM %s, "+
		" LATERAL json_array_elements_text(card.descriptions) AS d "+
		" WHERE card_status = $1 "+
		" AND ( d ~~* $2 "+
		" OR name ~~* $3) ",

	ALL_CARD_COLUMNS, CARD,
)

func (im *cardImpl) SearchCard(ctx context.Context, keyword string, status commonM.Status) ([]*cardDTO.CardDTO, error) {
	logPos := "[card.store][GetLatestCards]"

	cardDTOs := []*cardDTO.CardDTO{}

	var builder strings.Builder
	builder.WriteString("%")
	builder.WriteString(keyword)
	builder.WriteString("%")
	concatKeyword := builder.String()

	rows, err := im.primary.Query(SELECT_CARDS_BY_KEYWORD_STAT, status, concatKeyword, concatKeyword)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		cardDTO := &cardDTO.CardDTO{}
		selector := []interface{}{
			&cardDTO.ID,
			&cardDTO.Name,
			&cardDTO.Descriptions,
			&cardDTO.LinkURL,
			&cardDTO.BankID,
			&cardDTO.Order,
			&cardDTO.CardStatus,
			&cardDTO.CreateDate,
			&cardDTO.UpdateDate,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		cardDTOs = append(cardDTOs, cardDTO)
	}

	return cardDTOs, nil
}
