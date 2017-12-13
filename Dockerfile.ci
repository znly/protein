FROM google/cloud-sdk:alpine as gcloud

FROM golang:1.9.2

ARG GS_PATH=undefined
ENV GS_PATH ${GS_PATH}

COPY --from=gcloud /google-cloud-sdk /
COPY . /go/src/github.com/znly/protein

WORKDIR /go/src/github.com/znly/protein

RUN make -f Makefile.ci deps_external

RUN cp id_rsa /id_rsa \
      && chmod 600 /id_rsa \
      && echo "IdentityFile /id_rsa" >> /etc/ssh/ssh_config \
      && echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config \
      &&  git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN .buildkite/scripts/vendors.sh
RUN go generate . ./protoscan

RUN go install $(go list ./... | grep -v vendor) || true
