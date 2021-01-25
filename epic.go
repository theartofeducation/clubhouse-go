package clubhouse

// EpicState holds the State of the Epic.
type EpicState string

// Epic states.
var (
	EpicStateDone       EpicState = "done"
	EpicStateInProgress EpicState = "in progress"
)

// Epic holds the information for a Clubhouse Epic.
type Epic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
