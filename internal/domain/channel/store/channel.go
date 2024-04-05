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
	GetChannelsByType(ctx context.Context, channelCategoryType int32, status commonM.Status, limit, offset int32) ([]*channelDTO.ChannelDTO, error)
	GetChannelByID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error)
	GetChannelByIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error)
	SearchChannel(ctx context.Context, keyword string, status commonM.Status) ([]*channelDTO.ChannelDTO, error)
}

type impl struct {
	dig.In

	primary   *pgx.ConnPool
	migration *pgx.ConnPool
}

func New(sql *psql.Psql) ChannelStore {
	logPos := "[channel.store][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init channel store")

	return &impl{
		primary:   sql.Primary,
		migration: sql.Migration,
	}
}

const CHANNEL = "channel"
const ALL_COLUMNS = " \"id\", \"name\", \"link_url\", \"channel_type\", \"create_date\", \"update_date\", " +
	" \"channel_labels\", \"order\", \"channel_status\" "

var MODIFIED_CHANNEL_STAT = fmt.Sprintf("INSERT INTO %s "+
	" (%s) "+
	" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"+
	" ON CONFLICT(id) DO UPDATE SET  "+
	" \"name\" = $10, \"link_url\" = $11, \"channel_type\" = $12, "+
	" \"create_date\" = $13, \"update_date\" = $14, "+
	" \"channel_labels\" = $15, \"order\" = $16, "+
	" \"channel_status\" = $17 ", CHANNEL, ALL_COLUMNS)

func (im *impl) ModifiedChannel(ctx context.Context, channelDTO *channelDTO.ChannelDTO) error {
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
		channelDTO.LinkURL,
		channelDTO.ChannelType,
		channelDTO.ChannelLabels,
		channelDTO.Order,
		channelDTO.ChannelStatus,

		channelDTO.Name,
		channelDTO.LinkURL,
		channelDTO.ChannelType,
		channelDTO.ChannelLabels,
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

var SELECT_CHANNELS_BY_CHANNEL_TYPE_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s "+
	" WHERE \"channel_type\" = $1 "+
	" AND channel_status = $2 "+
	" ORDER BY \"order\" "+
	" LIMIT $3 OFFSET $4 ", ALL_COLUMNS, CHANNEL)

func (im *impl) GetChannelsByType(ctx context.Context, ctype int32, status commonM.Status, limit, offset int32) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.store][GetChannelsByType]"

	channelDTOs := []*channelDTO.ChannelDTO{}

	rows, err := im.primary.Query(SELECT_CHANNELS_BY_CHANNEL_TYPE_STAT, ctype, status, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"channel.type": ctype,
		}).Error("psql.Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		channelDTO := &channelDTO.ChannelDTO{}
		selector := []interface{}{
			&channelDTO.ID,
			&channelDTO.Name,
			&channelDTO.LinkURL,
			&channelDTO.ChannelType,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
			&channelDTO.Order,
			&channelDTO.ChannelStatus,
		}

		if err := rows.Scan(selector...); err != nil {
			log.WithFields(log.Fields{
				"pos":          logPos,
				"channel.type": ctype,
			}).Error("rows.Scan failed: ", err)
			return nil, err
		}
		channelDTOs = append(channelDTOs, channelDTO)
	}

	return channelDTOs, nil
}

var SELECT_CHANNEL_BY_ID_STAT = fmt.Sprintf("SELECT %s "+
	" FROM %s WHERE \"id\" = $1", ALL_COLUMNS, CHANNEL)

func (im *impl) GetChannelByID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error) {

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
			&channelDTO.LinkURL,
			&channelDTO.ChannelType,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
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
	" FROM %s WHERE \"id\" = ANY($1) ", ALL_COLUMNS, CHANNEL)

func (im *impl) GetChannelByIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error) {

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
			&channelDTO.LinkURL,
			&channelDTO.ChannelType,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
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
	" LIMIT 20 ", ALL_COLUMNS, CHANNEL)

func (im *impl) SearchChannel(ctx context.Context, keyword string, status commonM.Status) ([]*channelDTO.ChannelDTO, error) {
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
			&channelDTO.LinkURL,
			&channelDTO.ChannelType,
			&channelDTO.CreateDate,
			&channelDTO.UpdateDate,
			&channelDTO.ChannelLabels,
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
