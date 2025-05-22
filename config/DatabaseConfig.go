package config

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "point_calculator",
	}
}
