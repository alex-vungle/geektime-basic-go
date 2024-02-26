# 你可以直接执行 make 命令，也可以单独的命令复制到控制台。
# 注意，如果你是 Windows 并且不是在 WSL 下，
# 要注意文件分隔符使用 Windows 的分隔符。
.PHONY: generate
generate:
	@make mock
.PHONY: mock
mock:
	@go generate -tags=wireinject ./...
	@go mod tidy

.PHONY: grpc
grpc:
	@buf generate webook/api/proto

.PHONY: grpc_mock
grpc_mock:
	@mockgen -source=webook/api/proto/gen/article/v1/article_grpc.pb.go -package=artmocks -destination=webook/api/proto/gen/article/v1/mocks/article_grpc.mock.go
	@mockgen -source=webook/api/proto/gen/intr/v1/interactive_grpc.pb.go -package=intrmocks -destination=webook/api/proto/gen/intr/v1/mocks/interactive_grpc.mock.go
	@mockgen -source=webook/api/proto/gen/payment/v1/payment_grpc.pb.go -package=pmtmocks -destination=webook/api/proto/gen/payment/v1/mocks/payment_grpc.mock.go
	@mockgen -source=webook/api/proto/gen/follow/v1/follow_grpc.pb.go -package=followmocks -destination=webook/api/proto/gen/follow/v1/mocks/follow_grpc.mock.go


.PHONY: e2e
e2e:
	@docker compose -f webook/docker-compose.yaml down
	@docker compose -f webook/docker-compose.yaml up -d
	@go test -race ./webook/... -tags=e2e
	@docker compose -f webook/docker-compose.yaml down
.PHONY: e2e_up
e2e_up:
	@docker compose -f webook/docker-compose.yaml up -d
.PHONY: e2e_down
e2e_down:
	@docker compose -f webook/docker-compose.yaml down