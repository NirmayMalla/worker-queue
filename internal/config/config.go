package config

type Config struct {
	Port 				string
	WorkerCount int
	QueueSize		int
}

func Setup() Config {
	return Config {
		Port: ":8080",
		WorkerCount: 5,
		QueueSize: 100,
	}
}
