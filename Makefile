VERSION:=1.0.0-alpha
IMAGE:=cr.yandex/crpkmcbem8um7rd1gk5i/spotify_share

build:
	docker build . -t ${IMAGE}:${VERSION} -t ${IMAGE}

push:
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}

deploy: build push
