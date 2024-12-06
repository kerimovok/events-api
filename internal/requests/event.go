package requests

type CreateEventRequest struct {
	Name       string                 `json:"name" validate:"required,min=1,max=255"`
	Properties map[string]interface{} `json:"properties" validate:"required"`
}
