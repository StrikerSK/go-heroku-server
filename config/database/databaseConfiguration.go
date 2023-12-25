package database

type DatabaseConfiguration interface {
	GetType() string
}

type PostgresDatabaseConfiguration struct {
	databaseType     string
	databaseHost     string
	databasePort     string
	databaseName     string
	databaseUsername string
	databasePassword string
}

func (dc *PostgresDatabaseConfiguration) GetType() string {
	return dc.databaseType
}
