#
# Copyright (C) 2023 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

FROM --platform=linux/amd64 registry.access.redhat.com/ubi9/nodejs-18:1-80.1699550448 as node-builder

RUN npm install -g yarn

COPY ./podman-desktop-extension/src /opt/app-root/src/src
COPY ./podman-desktop-extension/package.json /opt/app-root/src/
COPY ./podman-desktop-extension/tsconfig.json /opt/app-root/src/
COPY ./podman-desktop-extension/vite.config.js /opt/app-root/src/
COPY ./podman-desktop-extension/vitest.config.js /opt/app-root/src/
COPY ./podman-desktop-extension/LICENSE /opt/app-root/src/
COPY ./podman-desktop-extension/icon.png /opt/app-root/src/
COPY ./podman-desktop-extension/README.md /opt/app-root/src/

RUN yarn && yarn build


FROM scratch as extension-builder
COPY --from=node-builder /opt/app-root/src/dist/ /extension/dist
COPY ./podman-desktop-extension/package.json /extension/
COPY ./podman-desktop-extension/LICENSE /extension/
COPY ./podman-desktop-extension/icon.png /extension/
COPY ./podman-desktop-extension/README.md /extension/


FROM --platform=linux/amd64 registry.access.redhat.com/ubi9/go-toolset:1.19.13-4.1697647145 as cli-builder
ARG PLATFORM_ARG

COPY ./go.mod /opt/app-root/src
COPY ./go.sum /opt/app-root/src
COPY ./main.go /opt/app-root/src
COPY ./pkg /opt/app-root/src/pkg/
COPY ./Makefile /opt/app-root/src

RUN make local-cross-${PLATFORM_ARG}


FROM --platform=$TARGETPLATFORM scratch
ARG PLATFORM_ARG

LABEL org.opencontainers.image.title="OpenShift Checker" \
        org.opencontainers.image.description="Analyze a Containerfile and highlight the directives and commands which could cause an unexpected behavior when running on an OpenShift cluster." \
        org.opencontainers.image.vendor="Red Hat" \
        io.podman-desktop.api.version=">= 1.5.3"

COPY --from=extension-builder /extension /extension

COPY --from=cli-builder /opt/app-root/src/bin/doa.cross.linux.${PLATFORM_ARG} /extension/doa.linux
COPY --from=cli-builder /opt/app-root/src/bin/doa.cross.windows.${PLATFORM_ARG} /extension/doa.exe
COPY --from=cli-builder /opt/app-root/src/bin/doa.cross.darwin.${PLATFORM_ARG} /extension/doa.darwin
