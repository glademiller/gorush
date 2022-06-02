OS := $(shell uname)

preview:
ifeq ($(OS),Darwin)
	sed -i "" -e "s/version:.*/version: $(DEVELOP_VERSION)/" Chart.yaml
	sed -i "" -e "s/version:.*/version: $(DEVELOP_VERSION)/" ../*/Chart.yaml
	sed -i "" -e "s/tag:.*/tag: $(DEVELOP_VERSION)/" values.yaml
else ifeq ($(OS),Linux)
	sed -i -e "s/version:.*/version: $(DEVELOP_VERSION)/" Chart.yaml
	sed -i -e "s/version:.*/version: $(DEVELOP_VERSION)/" ../*/Chart.yaml
	sed -i -e "s|repository:.*|repository: $(DOCKER_REGISTRY)\/netscale-technologies\/gorush|" values.yaml
	sed -i -e "s/tag:.*/tag: $(DEVELOP_VERSION)/" values.yaml
else
	echo "platfrom $(OS) not supported to release from"
	exit -1
endif
	echo "  version: $(DEVELOP_VERSION)" >> requirements.yaml
	jx step helm build --verbose
