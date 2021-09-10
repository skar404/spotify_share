VERSION:=1.1.7
IMAGE:=cr.yandex/crpkmcbem8um7rd1gk5i/spotify_share

build:
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 docker build . -t ${IMAGE}:${VERSION} -t ${IMAGE}

push:
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}

deploy: build push
