package clubhouse

// EntityType holds the type of Clubhouse Entity that was changed.
type EntityType string

// Entity types.
var (
	EntityTypeEpic EntityType = "epic"
)

// Action is an event that triggered a Webhook.
type Action string

// Webhook actions.
var (
	ActionUpdate Action = "update"
)

// Webhook holds the information for a Clubhouse Webhook.
type Webhook struct {
	Actions []WebhookAction
}

// WebhookAction holds the information for the Action that triggered the Webhook.
type WebhookAction struct {
	EntityType EntityType `json:"entity_type"`
	Action     Action
	Name       string
	Changes    WebhookActionChanges
}

// WebhookActionChanges hold the information for the changes that triggered the Webhook.
type WebhookActionChanges struct {
	State WebhookActionState
}

// WebhookActionState holds the state information for an Entity.
type WebhookActionState struct {
	New EpicState
}
