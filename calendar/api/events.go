package api

type Event struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	TimeZone    string   `json:"timezone"`
	Duration    int      `json:"duration"`
	Notes       []string `json:"notes"`
}
