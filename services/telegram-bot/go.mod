module github.com/timeb30/techstreamshop/services/telegram-bot

go 1.25.3

//replace github.com/timeb30/techsreamshop/pgk/kafkaclient => ../pkg/kafkaclient

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/timeb30/techstreamshop/pkg/kafkaclient v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.13.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)
