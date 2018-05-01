package crconfig

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
	"reflect"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func ExpectedMakeStats() tc.CRConfigStats {
	return tc.CRConfigStats{
		CDNName:   randStr(),
		TMHost:    randStr(),
		TMPath:    randStr(),
		TMUser:    randStr(),
		TMVersion: randStr(),
	}
}

func TestMakeStats(t *testing.T) {
	expected := ExpectedMakeStats()
	start := time.Now()
	actual := makeStats(*expected.CDNName, *expected.TMUser, *expected.TMHost, *expected.TMPath, *expected.TMVersion, nil)
	end := time.Now()
	expected.DateUnixSeconds = actual.DateUnixSeconds
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("makeStats expected: %+v, actual: %+v", expected, actual)
	}
	if actual.DateUnixSeconds == nil || *actual.DateUnixSeconds < start.Unix() || *actual.DateUnixSeconds > end.Unix() {
		t.Errorf("makeStats DateUniSeconds expected: < %+v > %+v, actual: %+v", start.Unix(), end.Unix(), actual.DateUnixSeconds)
	}
}
