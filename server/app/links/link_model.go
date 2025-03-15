package links

import (
	"database/sql"
	"time"
)

type LinkModel struct {
	Id          int64
	Uid         string
	Alias       string
	OriginalUrl string
	Name        string
	LifetimeSec int
	CreatedAt   time.Time
	OwnerId     sql.NullInt64
}
