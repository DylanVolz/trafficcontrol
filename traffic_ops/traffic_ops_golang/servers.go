package main

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

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"./db"
	"./model"
	"database/sql"
)

const privLevel = 20
const HiddenField = "********"

func serversXHandler(db *db.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getServersXResponse(q, db, privLevel)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getServersX(q url.Values, dbRef *db.DB, privLevel int) ([]Server, error) {

	tx := dbRef.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	servers := []Server{}
	if err := tx.Selectx(&servers, model.QueryServers(privLevel), db.Limit(0, 50)); err == sql.ErrNoRows {
	} else if err == nil {

	} else {
		return nil, err
	}

	return servers, nil
}

	//rows, err := db.Queryx(DbQuery())
	//defer rows.Close()
	//if err != nil {
	//	return nil, err
	//}
	//servers := []Server{}
	//for rows.Next() {
	//	var s Server
	//	err = rows.StructScan(&s)
	//	if privLevel < PrivLevelAdmin {
	//		s.IloPassword = HiddenFieldx
	//		s.XmppPasswd = HiddenField
	//	}
	//	servers = append(servers, s)
	//}
	//return servers, nil
//}

func getServersXResponse(q url.Values, db *db.DB, privLevel int) (*ServersResponse, error) {
	servers, err := getServersX(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}

	resp := ServersResponse{
		Response: servers,
	}
	return &resp, nil
}