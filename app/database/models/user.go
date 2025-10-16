package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Users struct {
	TgId      int64 `pg:",pk"`
	UserName  string
	FullName  string
	CreatedAt int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt int64 `pg:",default:extract(epoch from now())"`

	IsBotAdmin bool `pg:"default:false"`

	Groups []*Groups `pg:"many2many:group_participants,fk:user_tg_id,join_fk:group_tg_id"`
}

func (u *Users) AfterInsert(tx *pg.Tx) error {
	u.UpdatedAt = time.Now().Unix()
	return nil
}
