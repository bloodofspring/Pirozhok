package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Groups struct {
	GroupTgId int64
	CreatedAt            int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt            int64 `pg:",default:extract(epoch from now())"`

	Users []*Users `pg:"many2many:group_participants"`
}

func (g *Groups) AfterInsert(tx *pg.Tx) error {
	g.UpdatedAt = time.Now().Unix()
	return nil
}

type GroyupParticipants struct {
	UserTgId int64
	GroupTgId int64

	IsAdmin bool `pg:"default:false"`
}
