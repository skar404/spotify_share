VERSION:=1.1.1
IMAGE:=cr.yandex/crpkmcbem8um7rd1gk5i/spotify_share

build:
	docker build . -t ${IMAGE}:${VERSION} -t ${IMAGE}

push:
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}

deploy: build push
