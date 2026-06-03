package things

// Todo represents a single Things 3 to-do, serialized to JSON for the caller.
type Todo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Notes   string `json:"notes"`
	Status  string `json:"status"`
	Tags    string `json:"tags"`
	Due     string `json:"due,omitempty"`
	Start   string `json:"start,omitempty"`
	Project string `json:"project,omitempty"`
	Area    string `json:"area,omitempty"`
}

type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Area   string `json:"area,omitempty"`
	Status string `json:"status"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListKey is a canonical, language-independent identifier for one of the
// built-in Things 3 lists.
type ListKey string

const (
	ListInbox    ListKey = "inbox"
	ListToday    ListKey = "today"
	ListUpcoming ListKey = "upcoming"
	ListAnytime  ListKey = "anytime"
	ListSomeday  ListKey = "someday"
	ListLogbook  ListKey = "logbook"
	ListTrash    ListKey = "trash"
)
