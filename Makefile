

gen.service:
	mkdir services/$(service)
	mkdir services/$(service)/cmd
	mkdir services/$(service)/cmd/command_server
	mkdir services/$(service)/cmd/query_server
	mkdir services/$(service)/cmd/migrate
	mkdir services/$(service)/config
	mkdir services/$(service)/constant
	mkdir services/$(service)/docker
	mkdir services/$(service)/internal
	mkdir services/$(service)/internal/api
	mkdir services/$(service)/internal/api/grpc
	mkdir services/$(service)/internal/api/grpc/handler
	mkdir services/$(service)/internal/api/grpc/mapper
	mkdir services/$(service)/internal/api/messaging
	mkdir services/$(service)/internal/api/messaging/consumer
	mkdir services/$(service)/internal/api/messaging/dispatcher
	mkdir services/$(service)/internal/app
	mkdir services/$(service)/internal/app/command
	mkdir services/$(service)/internal/app/dto
	mkdir services/$(service)/internal/app/mapper
	mkdir services/$(service)/internal/app/query
	mkdir services/$(service)/internal/app/query/consumer
	mkdir services/$(service)/internal/app/service
	mkdir services/$(service)/internal/domain
	mkdir services/$(service)/internal/domain/entity
	mkdir services/$(service)/internal/domain/event
	mkdir services/$(service)/internal/domain/repository
	mkdir services/$(service)/internal/domain/valueobject
	mkdir services/$(service)/internal/infra
	mkdir services/$(service)/internal/infra/model
	mkdir services/$(service)/internal/infra/persistence
	mkdir services/$(service)/internal/infra/publisher
	mkdir services/$(service)/pkg
	mkdir services/$(service)/pkg/tracer
	mkdir services/$(service)/test

gen.go:
	buf generate