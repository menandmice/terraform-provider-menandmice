TEST		?=$$(go list ./... | grep -v 'vendor')
HOSTNAME	:=registry.terraform.io
NAMESPACE	:=local
NAME		:=menandmice
BINARY		:=terraform-provider-${NAME}
VERSION		:=0.2.1
OS			:=$(shell uname|tr A-Z a-z)
TF_VERSION  := "1.2.3"

ifeq ($(shell uname -m),x86_64)
  ARCH   ?= amd64
endif
ifeq ($(shell uname -m),i686)
  ARCH   ?= 386
endif
ifeq ($(shell uname -m),arm64)
  ARCH   ?= arm64
endif
ifeq ($(shell uname -m),aarch64)
  ARCH   ?= arm
endif

OS_ARCH		:=${OS}_${ARCH}
$(info ${OS_ARCH})

TERRAFORM_BIN_LOCATION	:= $(shell command -v terraform 2> /dev/null)

TERRAFORM_VERSION	:= "$(shell terraform version | awk 'NR == 1 {split ($$2 ,version, "v"); print version[2]}')"

TERRAFORM_PLATFORM	:= "$(shell terraform version | awk 'NR == 2 {split ($$2 ,platform_arch, " "); print platform_arch[1]}')"

$(info Version check)


default: install

install: build

build:
	go build -o ${BINARY}

##
# VERSION CHECK
# See if the installed Terraform version is greater than 0.12
#
ifeq ($(shell expr ${TERRAFORM_VERSION} \>= 0.12), 1)
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
else
	mkdir -p ~/.terraform.d/plugins/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${OS_ARCH}
endif
	rm examples/.terraform.lock.hcl ||true

generate_doc:
	tfplugindocs  generate # https://github.com/hashicorp/terraform-plugin-docs

example: init
	cd examples && terraform init && terraform apply -auto-approve && terraform destroy -auto-approve

init : install
	cd examples && terraform init

apply: init
	cd examples && terraform apply -auto-approve

plan: init
	cd examples&& terraform plan

destroy : init
	cd examples && terraform destroy -auto-approve

test:
	@[ "${TERRAFORM_BIN_LOCATION}" ] || ( echo ">> Terraform not installed"; exit 1 )
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
