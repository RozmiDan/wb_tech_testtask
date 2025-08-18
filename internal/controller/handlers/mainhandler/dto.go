package mainhandler

// ListGamesResponse — обёртка для GET /games
type ListGamesResponse struct {
	Data []entity.GameInList `json:"data"`
	Meta *Pagination         `json:"meta,omitempty"`
}

// --------------- ответы с ошибкой ---------------

// APIError — структура описания ошибки
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse — обёртка для не-200 ответов
type ErrorResponse struct {
	Error APIError `json:"error"`
}