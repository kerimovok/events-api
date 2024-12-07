package requests

type CreateEventRequest struct {
	Properties map[string]interface{} `json:"properties" validate:"required"`
}
