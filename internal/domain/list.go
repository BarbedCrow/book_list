package domain

type ListType string

const (
	ListTypeWantToRead ListType = "want_to_read"
	ListTypeReading    ListType = "reading"
	ListTypeFinished   ListType = "finished"
	ListTypeCustom     ListType = "custom"
)

type List struct {
	ID      string   `json:"id"`
	OwnerID string   `json:"owner_id"`
	Name    string   `json:"name"`
	Type    ListType `json:"type"`
	Books   []string `json:"books"`
}
