# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
############################################################
# Dockerfile to build Traffic Ops 1.6.0 container images
# Based on CentOS 7.2
############################################################

# Example Build and Run:
# docker network create --ipv6=true cdnetv6 --subnet="2001:db8:1::/64"
# docker build --rm --tag traffic_ops:2.1 --build-arg RPM=traffic_ops-2.1.0-6485.2303ddd4.el7.x86_64.rpm traffic_ops
#
# docker run --name my-traffic-ops-psql --hostname my-traffic-ops-psql --net cdnetv6 --detach postgres
#

#example replacing defaults file with your own defaults:
# docker run -v /yourpath/to/config.json:/opt/traffic_ops/app/defaults.json --name my-traffic-ops --hostname my-traffic-ops --net cdnetv6 --publish 443:443 --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnetv6 --detach traffic_ops:2.1
#example for dev use mounting the lib folder
# docker run -v /Users/dvolz200/workspace/my_trafficcontrol/traffic_ops/app/lib:/opt/traffic_ops/app/lib --name my-traffic-ops-2 --hostname my-traffic-ops --net cdnetv6 --publish 443:443 --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnetv6 --detach traffic_ops:2.1
# basic example:
# docker run --name my-traffic-ops --hostname my-traffic-ops --net cdnetv6 --publish 443:443 --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnetv6 --detach traffic_ops:2.1
# for strace capability add  --cap-add SYS_PTRACE

FROM centos
MAINTAINER Dan Kirkwood

RUN yum install -y perl cpanminus perl-Test-CPAN-Meta wget
RUN wget https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-redhat96-9.6-3.noarch.rpm
RUN rpm -ihv pgdg-redhat96-9.6-3.noarch.rpm
RUN cpanm Carton

# MUST provide RPM arg using --build-arg RPM=...  Can be local file or http://...
ARG RPM=traffic_ops-2.1.0-6486.3d55fbad.el7.x86_64.rpm
ADD $RPM /

RUN yum install -y /$(basename $RPM) postgresql96 postgresql96-devel git

# once installed, remove rpm to lower image size
RUN rm /$(basename $RPM)

RUN cd /opt/traffic_ops/app && carton # this is run again by postinstall but running it at build time speeds up image start up time (at the cost of build time...).

EXPOSE 443
WORKDIR /opt/traffic_ops/app
ENV MOJO_MODE production
ADD profile.origin.traffic_ops /
ADD run.sh /
ADD defaults.json /
RUN cp /defaults.json /opt/traffic_ops/app/defaults.json
CMD /run.sh
