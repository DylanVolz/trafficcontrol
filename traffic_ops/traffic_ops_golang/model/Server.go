package model

import "github.com/jmoiron/sqlx"

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

type ServersResponse struct {
	Version  string   `json:"version"`
	Response []Server `json:"response"`
}

type Server struct {
	Cachegroup     string `json:"cachegroup"`
	CachegroupId   int    `json:"cachegroupId"`
	CdnId          int    `json:"cdnId"`
	CdnName        string `json:"cdnName"`
	DomainName     string `json:"domainName"`
	Guid           string `json:"guid"`
	HostName       string `json:"hostName"`
	HttpsPort      int    `json:"httpsPort"`
	Id             int    `json:"id"`
	IloIpAddress   string `json:"iloIpAddress"`
	IloIpGateway   string `json:"iloIpGateway"`
	IloIpNetmask   string `json:"iloIpNetmask"`
	IloPassword    string `json:"iloPassword"`
	IloUsername    string `json:"iloUsername"`
	InterfaceMtu   string `json:"interfaceMtu"`
	InterfaceName  string `json:"interfaceName"`
	Ip6Address     string `json:"ip6Address"`
	Ip6Gateway     string `json:"ip6Gateway"`
	IpAddress      string `json:"ipAddress"`
	IpGateway      string `json:"ipGateway"`
	IpNetmask      string `json:"ipNetmask"`
	LastUpdated    string `json:"lastUpdated"`
	MgmtIpAddress  string `json:"mgmtIpAddress"`
	MgmtIpGateway  string `json:"mgmtIpGateway"`
	MgmtIpNetmask  string `json:"mgmtIpNetmask"`
	OfflineReason  string `json:"offlineReason"`
	PhysLocation   string `json:"physLocation"`
	PhysLocationId int    `json:"physLocationId"`
	Profile        string `json:"profile"`
	ProfileDesc    string `json:"profileDesc"`
	ProfileId      int    `json:"profileId"`
	Rack           string `json:"rack"`
	RouterHostName string `json:"routerHostName"`
	RouterPortName string `json:"routerPortName"`
	Status         string `json:"status"`
	StatusId       int    `json:"statusId"`
	TcpPort        int    `json:"tcpPort"`
	ServerType     string `json:"type"`
	ServerTypeId   int    `json:"typeId"`
	UpdPending     bool   `json:"updPending"`
	XmppId         string `json:"xmppId"`
	XmppPasswd     string `json:"xmppPasswd"`
}

var queryServers = `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroupId,
s.cdn_id as cdnId,
cdn.name as cdnName,
s.domain_name as domainName,
COALESCE(s.guid, '') as guid,
s.host_name as hostName,
COALESCE(s.https_port, 0) as httpsPort,
s.id as id,
COALESCE(s.ilo_ip_address, '') as iloIpAddress,
COALESCE(s.ilo_ip_gateway, '') as iloIpGateway,
COALESCE(s.ilo_ip_netmask, '') as iloIpNetmask,
COALESCE(s.ilo_password, '') as iloPassword,
COALESCE(s.ilo_username, '') as iloUsername,
s.interface_mtu as interfaceMtu,
s.interface_name as interfaceName,
s.ip6_address as ip6Address,
s.ip6_gateway as ip6Gateway,
s.ip_address as ipAddress,
s.ip_gateway as ipGateway,
s.ip_netmask as ipNetmask,
s.last_updated as lastUpdated,
s.mgmt_ip_address as mgmtIpAddress,
s.mgmt_ip_gateway as mgmtIpGateway,
s.mgmt_ip_netmask as mgmtIpNetmask,
s.offline_reason as offlineReason,
pl.name as physLocation,
s.phys_location as physLocationId,
p.name as profile,
p.description as profileDesc,
s.profile as profileId,
s.rack as rack,
s.router_host_name as routerHostName,
s.router_port_name as routerPortName,
st.name as status,
s.status as statusId,
s.tcp_port as tcpPort,
t.name as serverType,
s.type as serverTypeId,
s.upd_pending as updPending,
s.xmpp_id as xmppId,
s.xmpp_passwd as xmppPasswd
FROM server s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`
var queryServerById = `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroupId,
s.cdn_id as cdnId,
cdn.name as cdnName,
s.domain_name as domainName,
COALESCE(s.guid, '') as guid,
s.host_name as hostName,
COALESCE(s.https_port, 0) as httpsPort,
s.id as id,
COALESCE(s.ilo_ip_address, '') as iloIpAddress,
COALESCE(s.ilo_ip_gateway, '') as iloIpGateway,
COALESCE(s.ilo_ip_netmask, '') as iloIpNetmask,
COALESCE(s.ilo_password, '') as iloPassword,
COALESCE(s.ilo_username, '') as iloUsername,
s.interface_mtu as interfaceMtu,
s.interface_name as interfaceName,
s.ip6_address as ip6Address,
s.ip6_gateway as ip6Gateway,
s.ip_address as ipAddress,
s.ip_gateway as ipGateway,
s.ip_netmask as ipNetmask,
s.last_updated as lastUpdated,
s.mgmt_ip_address as mgmtIpAddress,
s.mgmt_ip_gateway as mgmtIpGateway,
s.mgmt_ip_netmask as mgmtIpNetmask,
s.offline_reason as offlineReason,
pl.name as physLocation,
s.phys_location as physLocationId,
p.name as profile,
p.description as profileDesc,
s.profile as profileId,
s.rack as rack,
s.router_host_name as routerHostName,
s.router_port_name as routerPortName,
st.name as status,
s.status as statusId,
s.tcp_port as tcpPort,
t.name as serverType,
s.type as serverTypeId,
s.upd_pending as updPending,
s.xmpp_id as xmppId,
s.xmpp_passwd as xmppPasswd
FROM server s
WHERE s.id =:serverId
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`

func QueryServers(privLevel int) db.Queryx {
	return db.Queryx{
		Query:  queryServers,
		Params: []interface{}{privLevel},
	}
}

func QueryServerByID(serverId int) db.Queryx {
	return db.Queryx{
		Query: queryServerById,
		Params: []interface{}{
			serverId,
		},
	}
}

func Select(*sqlx.Tx, db.Queryx, ...interface{}) error {
	query
	rows, err := db.Queryx(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	servers := []Server{}
	for rows.Next() {
		var s Server
		err = rows.StructScan(&s)
		if privLevel < PrivLevelAdmin {
			s.IloPassword = HiddenField
			s.XmppPasswd = HiddenField
		}
		servers = append(servers, s)
	}
	return servers, nil
}