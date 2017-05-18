// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !notime

package collector

import (
	"bufio"
	"fmt"
	"io"
	//"net"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/common/log"
)

type peerconnCollector struct {
	rateEntries      *prometheus.Desc
	delayEntries     *prometheus.Desc
	lastconnEntries  *prometheus.Desc
	lastdelayEntries *prometheus.Desc
}

func init() {
	Factories["peerconnect"] = NewPeerConnCollector
}

//type PeerConnectResultStru struct {
//	timeDelay  float64
//	successNum int
//	failedNum  int
//}

func NewPeerConnCollector() (Collector, error) {
	return &peerconnCollector{
		rateEntries: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "peerconnect", "success_rate"),
			"Rate of Peer Connection Success(0~1)",
			[]string{"peerdevice"}, nil,
		),
		delayEntries: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "peerconnect", "timedelay_average"),
			"Avarage Time Delay of Peer Connection(ms)",
			[]string{"peerdevice"}, nil,
		),
		lastconnEntries: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "peerconnect", "last_connect"),
			"Last Connection Record(0:failed,1:success)",
			[]string{"peerdevice"}, nil,
		),
		lastdelayEntries: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "peerconnect", "last_delay"),
			"Last Connection Delay(ms)",
			[]string{"peerdevice"}, nil,
		),
	}, nil
}

func getPeerConnectEntries() (map[string]float64, map[string]int, map[string]int, map[string]float64, map[string]int, error) {
	file, err := os.Open("/etc/proberesultnew")
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()

	delayEntries, succeedEntries, failedEntries, lastdelayEntries, lastconnectEntries, err := parsePeerConnectEntries(file)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return delayEntries, succeedEntries, failedEntries, lastdelayEntries, lastconnectEntries, nil
}

func parsePeerConnectEntries(data io.Reader) (map[string]float64, map[string]int, map[string]int, map[string]float64, map[string]int, error) {
	scanner := bufio.NewScanner(data)
	//entries := make(map[string]uint32)

	delayEntries := make(map[string]float64)
	succeedEntries := make(map[string]int)
	failedEntries := make(map[string]int)
	lastdelayEntries := make(map[string]float64)
	lastconnectEntries := make(map[string]int)
	var tempColumns string

	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())
		if strings.Contains(columns[0], "IP") == false {
			//tempResultStru := entries[columns[0]]
			tempColumns = columns[0] + ":" + columns[2] + "(" + columns[1] + ")"
			if strings.Contains(columns[3], "---") == true {
				//	tempResultStru.failedNum = tempResultStru.failedNum + 1
				failedEntries[tempColumns] = failedEntries[tempColumns] + 1
				delayEntries[tempColumns] = delayEntries[tempColumns] + 0
				lastconnectEntries[tempColumns] = 0
				lastdelayEntries[tempColumns] = -1
			} else {
				tempTimeDelay, err := strconv.ParseFloat(columns[3], 64)
				if err == nil {
					//tempResultStru.successNum = tempResultStru.successNum + 1
					//tempResultStru.timeDelay = tempResultStru.timeDelay + tempTimeDelay
					delayEntries[tempColumns] = delayEntries[tempColumns] + tempTimeDelay
					succeedEntries[tempColumns] = succeedEntries[tempColumns] + 1
					lastconnectEntries[tempColumns] = 1
					lastdelayEntries[tempColumns] = tempTimeDelay
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to parse peer connection info: %s", err)
	}

	return delayEntries, succeedEntries, failedEntries, lastdelayEntries, lastconnectEntries, nil
}

func (c *peerconnCollector) Update(ch chan<- prometheus.Metric) error {
	var tempSuccessRate, tempTimeDelay float64
	//var tempSuccess, tempFail int
	delayEntries, succeedEntries, failedEntries, lastDelayEntries, lastConnectEntries, err := getPeerConnectEntries()
	if err != nil {
		return fmt.Errorf("could not get peer connection entries: %s", err)
	}

	for addrConn, peerConResult := range delayEntries {
		tempSuccess, ok := succeedEntries[addrConn]
		if ok {
			//		tempSuccess = tempSuccess
		} else {
			tempSuccess = 0
		}

		tempFail, ok := failedEntries[addrConn]
		if ok {
			//		tempFail = tempFail
		} else {
			tempFail = 0
		}

		templastDelay, ok := lastDelayEntries[addrConn]
		if ok {
			if templastDelay != -1 {
				templastDelay = templastDelay * 1000
			}
		} else {
			templastDelay = -1
		}

		templastConn, ok := lastConnectEntries[addrConn]
		if ok {
			//              tempSuccess = tempSuccess
		} else {
			templastConn = 0
		}

		tempSuccessRate = (float64)(tempSuccess) / ((float64)(tempSuccess) + (float64)(tempFail))
		if tempSuccess == 0 {
			tempTimeDelay = -1
		} else {
			tempTimeDelay = (float64)(peerConResult * 1000 / ((float64)(tempSuccess)))
		}

		//ch <- prometheus.MustNewConstMetric(c.rateEntries, prometheus.GaugeValue, float64(tempSuccessRate), addrConn)
		//ch <- prometheus.MustNewConstMetric(c.delayEntries, prometheus.GaugeValue, float64(tempTimeDelay), addrConn)
		ch <- prometheus.MustNewConstMetric(c.rateEntries, prometheus.GaugeValue, float64(tempSuccessRate), addrConn)
		ch <- prometheus.MustNewConstMetric(c.delayEntries, prometheus.GaugeValue, float64(tempTimeDelay), addrConn)
		ch <- prometheus.MustNewConstMetric(c.lastconnEntries, prometheus.GaugeValue, float64(templastConn), addrConn)
		ch <- prometheus.MustNewConstMetric(c.lastdelayEntries, prometheus.GaugeValue, float64(templastDelay), addrConn)
	}
	return nil
}
