#!/bin/bash
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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.

# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# DOMAIN
# TRAFFIC_VAULT_PASS
envvars=( DOMAIN TRAFFIC_VAULT_PASS)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

start() {
		/etc/rc.d/init.d/traffic_ops start
		exec tail -f /var/log/traffic_ops/traffic_ops.log
}

init() {
		echo "{\"user\": \"riakuser\",\"password\": \"$TRAFFIC_VAULT_PASS\"}" > /opt/traffic_ops/app/conf/production/riak.conf
		#pull needed values from file: /opt/traffic_ops/app/defaults.json (this json file is a dictionary of lists)
		# key = "/opt/traffic_ops/app/conf/production/database.conf"
		PSQL_IP="$( cat /opt/traffic_ops/app/defaults.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["Database server hostname IP or FQDN"] for x in obj["/opt/traffic_ops/app/conf/production/database.conf"] if "Database server hostname IP or FQDN" in x]; print match[0]')"
		PSQL_PORT="$( cat /opt/traffic_ops/app/defaults.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["Database port number"] for x in obj["/opt/traffic_ops/app/conf/production/database.conf"] if "Database port number" in x]; print match[0]')"
		PSQL_TRAFFIC_OPS_PASS="$( cat /opt/traffic_ops/app/defaults.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["Password for Traffic Ops database user"] for x in obj["/opt/traffic_ops/app/conf/production/database.conf"] if "Password for Traffic Ops database user" in x]; print match[0]')"
		# key = "/opt/traffic_ops/install/data/json/users.json"
		ADMIN_USER="$( cat /opt/traffic_ops/app/defaults.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["Administration username for Traffic Ops"] for x in obj["/opt/traffic_ops/install/data/json/users.json"] if "Administration username for Traffic Ops" in x]; print match[0]')"
		ADMIN_PASS="$( cat /opt/traffic_ops/app/defaults.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["Password for the admin user"] for x in obj["/opt/traffic_ops/install/data/json/users.json"] if "Password for the admin user" in x]; print match[0]')"  	
     	
		TODB_HOST=$PSQL_IP TODB_PORT=$PSQL_PORT TODB_USERNAME_PASSWORD=$PSQL_TRAFFIC_OPS_PASS /opt/traffic_ops/install/bin/todb_bootstrap.sh

		export TERM=xterm && export USER=root && /opt/traffic_ops/install/bin/postinstall -cfile=/opt/traffic_ops/app/defaults.json

		TRAFFIC_OPS_URI="https://localhost"

		TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$ADMIN_USER"'", "p":"'"$ADMIN_PASS"'" }' $TRAFFIC_OPS_URI/api/1.2/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
		echo "Got cookie: $TMP_TO_COOKIE"

		TMP_DOMAIN=$DOMAIN
		sed -i -- "s/{{.Domain}}/$TMP_DOMAIN/g" /profile.origin.traffic_ops
		echo "Got domain: $TMP_DOMAIN"

		echo "Importing origin"
		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" -F "filename=profile.origin.traffic_ops" -F "profile_to_import=@/profile.origin.traffic_ops" $TRAFFIC_OPS_URI/profile/doImport

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "division.name=East" $TRAFFIC_OPS_URI/division/create
		TMP_DIVISION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/region/add | grep --color=never -oE "<option value=\"[0-9]+\">East</option>" | grep --color=never -oE "[0-9]+")"
		echo "Got division ID: $TMP_DIVISION_ID"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "region.name=Eastish" --data-urlencode "region.division_id=$TMP_DIVISION_ID" $TRAFFIC_OPS_URI/region/create
		TMP_REGION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/regions.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="Eastish"]; print match[0]')"
		echo "Got region ID: $TMP_REGION_ID"

		TMP_CACHEGROUP_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="MID_LOC"]; print match[0]')"
		echo "Got cachegroup type ID: $TMP_CACHEGROUP_TYPE"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "cg_data.name=mid-east" --data-urlencode "cg_data.short_name=east" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=-1" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
		TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
		echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"

		TMP_CACHEGROUP_EDGE_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="EDGE_LOC"]; print match[0]')"
		echo "Got cachegroup type ID: $TMP_CACHEGROUP_EDGE_TYPE"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "cg_data.name=edge-east" --data-urlencode "cg_data.short_name=eeast" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=$TMP_CACHEGROUP_ID" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_EDGE_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
		TMP_CACHEGROUP_EDGE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="edge-east"]; print match[0]')"
		echo "Got cachegroup edge ID: $TMP_CACHEGROUP_EDGE_ID"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "location.name=plocation-nyc-1" --data-urlencode "location.short_name=nyc" --data-urlencode "location.address=1 Main Street" --data-urlencode "location.city=nyc" --data-urlencode "location.state=NY" --data-urlencode "location.zip=12345" --data-urlencode "location.poc=" --data-urlencode "location.phone=" --data-urlencode "location.email=no@no.no" --data-urlencode "location.comments=" --data-urlencode "location.region=$TMP_REGION_ID" $TRAFFIC_OPS_URI/phys_location/create

		echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
