package model

type CreateToDoItemRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateToDoItemResponse struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateToDoItemRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateToDoItemResponse struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeleteToDoItemResponse struct {
	Success bool `json:"success,omitempty"`
}

type GetToDoItemsData struct {
	ID          int   `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetToDoItemsResponse struct {
	Data  []GetToDoItemsData `json:"data,omitempty"`
	Page  int          `json:"page,omitempty"`
	Limit int          `json:"limit,omitempty"`
	Total int          `json:"total,omitempty"`
}

type RequestError struct {
	Error string `json:"error,omitempty"`
}
