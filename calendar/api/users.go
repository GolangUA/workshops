package api

type User struct {
	Username string `json:"username"`
}

type UserPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserTimezone struct {
	Username string `json:"username"`
	Timezone string `json:"timezone"`
}
