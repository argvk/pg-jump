/*
Copyright 2017 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package proxy

import (
	"github.com/crunchydata/crunchy-proxy/config"
	"github.com/golang/glog"
	"net"
)

func ReturnConnection(ch chan int, connIndex int) {
	glog.V(2).Infof("returning poolIndex %d\n", connIndex)
	ch <- connIndex
}

func SetupPools() {
	if !config.Cfg.Pool.Enabled {
		glog.Errorln("[pool] pooling not enabled")
		return
	}

	glog.V(2).Infoln("[pool] pooling enabled")

	for i := 0; i < len(config.Cfg.Replicas); i++ {
		setupPoolForNode(&config.Cfg.Replicas[i])
	}

	setupPoolForNode(&config.Cfg.Master)

}

func setupPoolForNode(node *config.Node) {
	var err error

	node.Pool.Channel = make(chan int, config.Cfg.Pool.Capacity)
	node.Pool.Connections = make([]*net.TCPConn, config.Cfg.Pool.Capacity)
	for j := 0; j < config.Cfg.Pool.Capacity; j++ {
		node.Pool.Channel <- j
		//add a connection to the node pool
		glog.V(2).Infoln("[pool] adding conn to node %s pool\n", node.HostPort)
		node.Pool.Connections[j], err = node.GetConnection()
		if err != nil {
			glog.Errorln("error in getting pool conn for node " + err.Error())
		}
		Authenticate(node, node.Pool.Connections[j])
	}
}
