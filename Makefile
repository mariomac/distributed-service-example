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

.PHONY: push-backend-kind push-frontend-kind push-worker-kind push-loadgen-kind push-all-kind
push-backend-kind:
	kind load docker-image $(IMG_BASE):backend
push-frontend-kind:
	kind load docker-image $(IMG_BASE):frontend
push-worker-kind:
	kind load docker-image $(IMG_BASE):worker
push-loadgen-kind:
	kind load docker-image $(IMG_BASE):loadgen
push-all-kind: push-backend-kind push-frontend-kind push-worker-kind push-loadgen-kind

.PHONY: push-backend-k3d push-frontend-k3d push-worker-k3d push-loadgen-k3d push-all-k3d
push-backend-k3d:
	k3d image import $(IMG_BASE):backend
push-frontend-k3d:
	k3d image import $(IMG_BASE):frontend
push-worker-k3d:
	k3d image import $(IMG_BASE):worker
push-loadgen-k3d:
	k3d image import $(IMG_BASE):loadgen
push-all-k3d: push-backend-k3d push-frontend-k3d push-worker-k3d push-loadgen-k3d

.PHONY: hub-buildx-push-all hub-buildx-push-backend hub-buildx-push-frontend hub-buildx-push-worker hub-buildx-push-loadgen
hub-buildx-push-backend:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(IMG_BASE):backend backend/.
hub-buildx-push-frontend:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(IMG_BASE):frontend frontend/.
hub-buildx-push-worker:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(IMG_BASE):worker worker/.
hub-buildx-push-loadgen:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(IMG_BASE):loadgen loadgen/.
hub-buildx-push-all: hub-buildx-push-backend hub-buildx-push-frontend hub-buildx-push-worker hub-buildx-push-loadgen