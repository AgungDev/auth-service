SWAG ?= swag

.PHONY: swagger
swagger:
	$(SWAG) init -g cmd/server/main.go
