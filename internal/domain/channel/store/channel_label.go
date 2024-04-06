package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	psql "pickrewardapi/internal/pkg/postgres"

	channelDTO "pickrewardapi/internal/domain/channel/dto"
)

type ChannelLabelStore interface {
	ModifiedChannelLabel(ctx context.Context, channelDTO *channelDTO.ChannelLabelDTO) error
	GetAllChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error)
	GetChannelLabelByLabel(ctx context.Context, label int32) (*channelDTO.ChannelLabelDTO, error)
}

type channelLabelImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewChannelLabel(sql *psql.Psql) ChannelLabelStore {
	logPos := "[channel_label.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init channel label store")

	return &channelLabelImpl{
		primary: sql.Primary,
	}
}

const CHANNEL_LABEL = "channel_label"
const ALL_CHANNEL_LABEL_COLUMNS = " \"label\", \"name\", \"show\" "

var MODIFIED_CHANNEL_LABEL_STAT = fmt.Sprintf(
	"INSERT INTO %s "+
		"(%s) VALUES ($1, $2, $3)"+
		" ON CONFLICT(label) DO UPDATE SET  "+
		" \"name\" = $4, \"show\" = $5 ",
	CHANNEL_LABEL, ALL_CHANNEL_LABEL_COLUMNS,
)

func (im *channelLabelImpl) ModifiedChannelLabel(ctx context.Context, channelLabelDTO *channelDTO.ChannelLabelDTO) error {
	logPos := "[channel_label.store][ModifiedChannelLabel]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":                   logPos,
			"channelLabelDTO.Label": channelLabelDTO.Label,
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
		channelLabelDTO.Label,
		channelLabelDTO.Name,
		channelLabelDTO.Show,

		channelLabelDTO.Name,
		channelLabelDTO.Show,
	}

	if _, err := tx.Exec(MODIFIED_CHANNEL_LABEL_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":                   logPos,
			"channelLabelDTO.Label": channelLabelDTO.Label,
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

var SELECT_CHANNEL_LABELS_STAT = fmt.Sprintf(
	" SELECT %s "+
		" FROM %s ORDER BY label ",
	ALL_CHANNEL_LABEL_COLUMNS, CHANNEL_LABEL,
)

func (im *channelLabelImpl) GetAllChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error) {
	logPos := "[channel_label.store][GetAllChannelLabels]"

	channelLabelDTOs := []*channelDTO.ChannelLabelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNEL_LABELS_STAT)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		channelLabelDTO := &channelDTO.ChannelLabelDTO{}

		selector := []interface{}{
			&channelLabelDTO.Label,
			&channelLabelDTO.Name,
			&channelLabelDTO.Show,
		}
		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		channelLabelDTOs = append(channelLabelDTOs, channelLabelDTO)
	}

	return channelLabelDTOs, nil
}

var SELECT_CHANNEL_LABEL_BY_LABEL_STAT = fmt.Sprintf(
	" SELECT %s "+
		" FROM %s "+
		" WHERE label = $1 ",
	ALL_CHANNEL_LABEL_COLUMNS, CHANNEL_LABEL,
)

func (im *channelLabelImpl) GetChannelLabelByLabel(ctx context.Context, label int32) (*channelDTO.ChannelLabelDTO, error) {
	logPos := "[channel_label.store][GetAllChannelLabels]"

	channelLabelDTO := &channelDTO.ChannelLabelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNEL_LABEL_BY_LABEL_STAT, label)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		selector := []interface{}{
			&channelLabelDTO.Label,
			&channelLabelDTO.Name,
			&channelLabelDTO.Show,
		}
		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":   logPos,
				"label": label,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

		if rows.Next() {
			log.WithFields(log.Fields{
				"pos":   logPos,
				"label": label,
			}).Error("There have more than one record.")
			return nil, errors.New("There have more than one record.")
		}
	}

	return channelLabelDTO, nil

}
