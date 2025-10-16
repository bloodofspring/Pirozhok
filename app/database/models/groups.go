package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Groups struct {
	TgId      int64 `pg:",pk,notnull"`
	CreatedAt int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt int64 `pg:",default:extract(epoch from now())"`

	Users []*Users `pg:"many2many:group_participants,fk:group_tg_id,join_fk:user_tg_id"`
}

func (g *Groups) AfterInsert(tx *pg.Tx) error {
	g.UpdatedAt = time.Now().Unix()
	return nil
}

type GroupParticipants struct {
	UserTgId  int64 `pg:",pk"`
	GroupTgId int64 `pg:",pk"`

	IsAdmin bool `pg:"default:false"`
}
