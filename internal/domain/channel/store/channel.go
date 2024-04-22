package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	channelDTO "pickrewardapi/internal/domain/channel/dto"
	psql "pickrewardapi/internal/pkg/postgres"
	commonM "pickrewardapi/internal/shared/common/model"
)

type ChannelStore interface {
	ModifiedChannel(ctx context.Context, channelDTO *channelDTO.ChannelDTO) error
	GetChannelsByShowLabel(ctx context.Context, showLabel string, status commonM.Status, limit, offset int32) ([]*channelDTO.ChannelDTO, error)
	GetChannelByID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error)
	GetChannelByIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error)
	SearchChannel(ctx context.Context, keyword string, status commonM.Status) ([]*channelDTO.ChannelDTO, error)
}

type channelImpl struct {
	dig.In

	primary *pgx.ConnPool
}

func NewChannel(sql *psql.Psql) ChannelStore {
	logPos := "[channel.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init channel store")

	return &channelImpl{
		primary: sql.Primary,
	}
}

const CHANNEL = "channel"
const ALL_CAHNNEL_COLUMNS = " \"id\", \"name\", \"create_date\", \"update_date\", " +
	" \"channel_labels\", \"show_label\", \"order\", \"channel_status\" "

var MODIFIED_CHANNEL_STAT = fmt.Sprintf("INSERT INTO %s "+
	" (%s) "+
	" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"+
	" ON CONFLICT(id) DO UPDATE SET  "+
	" \"name\" = $9, \"create_date\" = $10, \"update_date\" = $11, "+
	" \"channel_labels\" = $12, \"show_label\" = $13, \"order\" = $14, "+
	" \"channel_status\" = $15 ", CHANNEL, ALL_CAHNNEL_COLUMNS)

func (im *channelImpl) ModifiedChannel(ctx context.Context, channelDTO *channelDTO.ChannelDTO) error {
	logPos := "[channel.store][ModifiedChannel]"

	tx, err := im.primary.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"pos":        logPos,
			"chennel.ID": channelDTO.ID,
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
		channelDTO.ID,
		channelDTO.Name,
		channelDTO.CreateDate,
		channelDTO.UpdateDate,
		channelDTO.ChannelLabels,
		channelDTO.ShowLabel,
		channelDTO.Order,
		channelDTO.ChannelStatus,

		channelDTO.Name,
		channelDTO.CreateDate,
		channelDTO.UpdateDate,
		channelDTO.ChannelLabels,
		channelDTO.ShowLabel,
		channelDTO.Order,
		channelDTO.ChannelStatus,
	}

	if _, err := tx.Exec(MODIFIED_CHANNEL_STAT, updater...); err != nil {
		log.WithFields(log.Fields{
			"pos":        logPos,
			"chennel.ID": channelDTO.ID,
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

var SELECT_CHANNELS_BY_SHOW_LABEL_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s "+
	" WHERE \"show_label\" = $1 "+
	" AND channel_status = $2 "+
	" ORDER BY \"order\" "+
	" LIMIT $3 OFFSET $4 ", ALL_CAHNNEL_COLUMNS, CHANNEL)

func (im *channelImpl) GetChannelsByShowLabel(ctx context.Context, showLabel string, status commonM.Status, limit, offset int32) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.store][GetChannelsByType]"

	channelDTOs := []*channelDTO.ChannelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNELS_BY_SHOW_LABEL_STAT, showLabel, status, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":                logPos,
			"channel.show_label": showLabel,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		channelDTO := &channelDTO.ChannelDTO{}
		selector := []interface{}{
			&channelDTO.ID,
			&channelDTO.Name,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
			&channelDTO.ShowLabel,
			&channelDTO.Order,
			&channelDTO.ChannelStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":                logPos,
				"channel.show_label": showLabel,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		channelDTOs = append(channelDTOs, channelDTO)
	}

	return channelDTOs, nil
}

var SELECT_CHANNEL_BY_ID_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s WHERE \"id\" = $1", ALL_CAHNNEL_COLUMNS, CHANNEL)

func (im *channelImpl) GetChannelByID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error) {

	logPos := "[channel.store][GetByChannelID]"

	channelDTO := &channelDTO.ChannelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNEL_BY_ID_STAT, ID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":        logPos,
			"channel.ID": ID,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		selector := []interface{}{
			&channelDTO.ID,
			&channelDTO.Name,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
			&channelDTO.ShowLabel,
			&channelDTO.Order,
			&channelDTO.ChannelStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":        logPos,
				"channel.ID": ID,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}

	}

	return channelDTO, nil
}

var SELECT_CHANNELS_BY_IDs_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s WHERE \"id\" = ANY($1) ", ALL_CAHNNEL_COLUMNS, CHANNEL)

func (im *channelImpl) GetChannelByIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error) {

	logPos := "[channel.store][GetChannelByIDs]"

	channelDTOs := []*channelDTO.ChannelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNELS_BY_IDs_STAT, pq.Array(IDs))
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		channelDTO := &channelDTO.ChannelDTO{}
		selector := []interface{}{
			&channelDTO.ID,
			&channelDTO.Name,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
			&channelDTO.ShowLabel,
			&channelDTO.Order,
			&channelDTO.ChannelStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		channelDTOs = append(channelDTOs, channelDTO)
	}
	return channelDTOs, nil
}

var SELECT_CHANNELS_BY_KEYWORD_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s "+
	" WHERE channel_status = $1 "+
	" AND name ~~* $2 "+
	" LIMIT 20 ", ALL_CAHNNEL_COLUMNS, CHANNEL)

func (im *channelImpl) SearchChannel(ctx context.Context, keyword string, status commonM.Status) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.store][SearchChannel]"

	channelDTOs := []*channelDTO.ChannelDTO{}

	var builder strings.Builder
	builder.WriteString("%")
	builder.WriteString(keyword)
	builder.WriteString("%")
	concatKeyword := builder.String()

	rows, err := im.primary.Query(SELECT_CHANNELS_BY_KEYWORD_STAT, status, concatKeyword)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		channelDTO := &channelDTO.ChannelDTO{}
		selector := []interface{}{
			&channelDTO.ID,
			&channelDTO.Name,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
			&channelDTO.ShowLabel,
			&channelDTO.Order,
			&channelDTO.ChannelStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		channelDTOs = append(channelDTOs, channelDTO)
	}
	return channelDTOs, nil
}
