package dto

type RedisSecrets struct {
	Address    string `json:"redis_address"`
	Password   string `json:"redis_password"`
	User       string `json:"redis_user"`
	DatabaseId int    `json:"redis_database_id"`
}
