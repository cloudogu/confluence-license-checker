MAKEFILES_VERSION=9.0.4

# Set these to the desired values
ARTIFACT_ID=confluence-license-checker
VERSION=0.1.0
GOTAG=1.12.7-1

GO_ENVIRONMENT=GO111MODULE=on
.DEFAULT_GOAL:=compile

include build/make/variables.mk
include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/package-tar.mk
include build/make/deploy-debian.mk
include build/make/digital-signature.mk
include build/make/self-update.mk
include build/make/release.mk
