package domain

type ListType string

const (
	ListTypeWantToRead ListType = "want_to_read"
	ListTypeReading    ListType = "reading"
	ListTypeFinished   ListType = "finished"
	ListTypeCustom     ListType = "custom"
)

type List struct {
	ID      string
	OwnerID string
	Name    string
	Type    ListType
	Books   []string
}
