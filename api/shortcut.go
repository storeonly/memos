package api

type Shortcut struct {
	Id        int   `json:"id"`
	CreatedTs int64 `json:"createdTs"`
	UpdatedTs int64 `json:"updatedTs"`

	Title     string `json:"title"`
	Payload   string `json:"payload"`
	RowStatus string `json:"rowStatus"`
	CreatorId int
}

type ShortcutCreate struct {
	// Standard fields
	CreatorId int

	// Domain specific fields
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

type ShortcutPatch struct {
	Id int

	Title     *string `json:"title"`
	Payload   *string `json:"payload"`
	RowStatus *string `json:"rowStatus"`
}

type ShortcutFind struct {
	Id *int

	// Standard fields
	CreatorId *int

	// Domain specific fields
	Title *string `json:"title"`
}

type ShortcutDelete struct {
	Id int
}

type ShortcutService interface {
	CreateShortcut(create *ShortcutCreate) (*Shortcut, error)
	PatchShortcut(patch *ShortcutPatch) (*Shortcut, error)
	FindShortcutList(find *ShortcutFind) ([]*Shortcut, error)
	FindShortcut(find *ShortcutFind) (*Shortcut, error)
	DeleteShortcut(delete *ShortcutDelete) error
}
