package models

type Users struct {
	TgId                 int64 `pg:",pk"`
	UserName             string
	FullName             string
	CreatedAt            int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt            int64 `pg:",default:extract(epoch from now())"`

	IsAdmin               bool `pg:"default:false"`
	IsSuperAdmin          bool `pg:"default:false"`
}
