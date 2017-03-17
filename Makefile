# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR:
.SUFFIXES:

# Constants, these can be overwritten in your Makefile.local
BUILD_SERVER := magneticio/buildserver

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

# Targets
.PHONY: all
all: default

# Using our buildserver which contains all the necessary dependencies
.PHONY: default
default:
	docker pull $(BUILD_SERVER)
	docker run \
		--rm \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume $(shell command -v docker):/usr/bin/docker \
		--volume $(CURDIR):/srv/src \
		--workdir=/srv/src \
		$(BUILD_SERVER) \
			./build.sh --build

