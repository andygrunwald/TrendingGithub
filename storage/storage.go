package storage

type Storage interface {
	NewPool(url, auth string) Pool
}

type Pool interface {
	Close() error
	Get() Connection
}

type Connection interface {
	// Close closes the connection.
	Close() error

	MarkRepositoryAsTweeted(projectName, score string) (bool, error)
	IsRepositoryAlreadyTweeted(projectName string) (bool, error)
}
