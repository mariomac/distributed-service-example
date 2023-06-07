IMG_BASE = mariomac/ds-example

.PHONY: build-backend build-frontend build-worker build-loadgen build-all
build-backend:
	docker build -t $(IMG_BASE):backend backend/.
build-frontend:
	docker build -t $(IMG_BASE):frontend frontend/.
build-worker:
	docker build -t $(IMG_BASE):worker worker/.
build-loadgen:
	docker build -t $(IMG_BASE):loadgen loadgen/.
build-all: build-backend build-frontend build-worker build-loadgen

.PHONY: push-backend push-frontend push-worker push-loadgen push-all
push-backend:
	kind load docker-image $(IMG_BASE):backend
push-frontend:
	kind load docker-image $(IMG_BASE):frontend
push-worker:
	kind load docker-image $(IMG_BASE):worker
push-loadgen:
	kind load docker-image $(IMG_BASE):loadgen
push-all: push-backend push-frontend push-worker push-loadgen

