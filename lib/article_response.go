package readraptor

type ArticleResponse struct {
	Created     *Timestamp `db:"created_at"       json:"created"`
	Updated     *Timestamp `db:"updated_at"       json:"updated"`
	Key         string     `db:"key"              json:"key"`
	FirstReadAt *Timestamp `db:"first_read_at" json:"first_read_at,omitempty"`
	LastReadAt  *Timestamp `db:"last_read_at"  json:"last_read_at,omitempty"`
}
