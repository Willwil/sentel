//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"
	uuid "github.com/satori/go.uuid"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

// MQTT service declaration
type mqttService struct {
	core.ServiceBase
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	localAddrs []string
	storage    Storage
	protocol   string
	stats      *base.Stats
	metrics    *base.Metrics
	consumer   sarama.Consumer
}

const (
	MqttProtocolTcp = "tcp"
	MqttProtocolWs  = "ws"
	MqttProtocolTls = "tls"
)

// MqttFactory
type MqttFactory struct {
	Protocol string
}

// New create mqtt service factory
func (this *MqttFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// Get all local ip address
	localAddrs := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Errorf("Failed to get local interface:%s", err)
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			return nil, errors.New("Failed to get local address")
		}
	}

	// Create storage
	s, err := NewStorage(c)
	if err != nil {
		return nil, errors.New("Failed to create storage in mqtt")
	}
	// kafka
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	t := &mqttService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   this.Protocol,
		localAddrs: localAddrs,
		storage:    s,
		stats:      base.NewStats(true),
		metrics:    base.NewMetrics(true),
		consumer:   consumer,
	}
	return t, nil
}

// MQTT Service

// Name
func (this *mqttService) Name() string {
	return "mqtt:" + this.protocol
}

func (this *mqttService) NewSession(conn net.Conn) (base.Session, error) {
	id := this.createSessionId()
	s, err := newMqttSession(this, conn, id)
	return s, err
}

// CreateSessionId create id for new session
func (this *mqttService) createSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (this *mqttService) getSessionTotalCount() int64 {
	return int64(len(this.sessions))
}

func (this *mqttService) removeSession(s base.Session) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessions[s.Identifier()] = nil
}
func (this *mqttService) registerSession(s base.Session) {
	this.mutex.Lock()
	this.sessions[s.Identifier()] = s
	this.mutex.Unlock()
}

// Info
func (this *mqttService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "mqtt",
	}
}

// Stats and Metrics
func (this *mqttService) GetStats() *base.Stats     { return this.stats }
func (this *mqttService) GetMetrics() *base.Metrics { return this.metrics }

// Client
func (this *mqttService) GetClients() []*base.ClientInfo       { return nil }
func (this *mqttService) GetClient(id string) *base.ClientInfo { return nil }
func (this *mqttService) KickoffClient(id string) error        { return nil }

// Session Info
func (this *mqttService) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (this *mqttService) GetSession(id string) *base.SessionInfo                     { return nil }

// Route Info
func (this *mqttService) GetRoutes() []*base.RouteInfo { return nil }
func (this *mqttService) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (this *mqttService) GetTopics() []*base.TopicInfo       { return nil }
func (this *mqttService) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (this *mqttService) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (this *mqttService) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (this *mqttService) GetServiceInfo() *base.ServiceInfo { return nil }

// Start is mainloop for mqtt service
func (this *mqttService) Start() error {
	this.subscribeTopic("session")

	host, _ := this.Config.String(this.protocol, "listen")
	listen, err := net.Listen("tcp", host)
	if err != nil {
		glog.Errorf("Mqtt listen failed:%s", err)
		return err
	}
	glog.Infof("Mqtt server is listening on '%s'...", host)
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		session, err := this.NewSession(conn)
		if err != nil {
			glog.Errorf("Mqtt create session failed:%s", err)
			return err
		}
		this.registerSession(session)
		go func(s base.Session) {
			err := s.Handle()
			if err != nil {
				conn.Close()
				glog.Error(err)
			}
		}(session)
	}
	// notify main
	// this.quit <- 1
	return nil
}

// Stop
func (this *mqttService) Stop() {
}

// launchMqttMonitor
func (this *mqttService) launchMqttMonitor() error {
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (this *mqttService) subscribeTopic(topic string) error {
	partitionList, err := this.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := this.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		this.WaitGroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer this.WaitGroup.Done()
			for msg := range pc.Messages() {
				this.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}

// handleNotifications handle notification from kafka
func (this *mqttService) handleNotifications(topic string, value []byte) error {
	switch topic {
	case TopicNameSession:
		return this.handleSessionNotifications(value)
	}
	return nil
}

// handleSessionNotifications handle session notification  from kafka
func (this *mqttService) handleSessionNotifications(value []byte) error {
	// Decode value received form other mqtt node
	var topics []SessionTopic
	if err := json.Unmarshal(value, &topics); err != nil {
		glog.Errorf("Mqtt session notifications failure:%s", err)
		return err
	}
	// Get local ip address
	for _, topic := range topics {
		switch topic.Action {
		case ObjectActionUpdate:
			// Only deal with notification that is not  launched by myself
			for _, addr := range this.localAddrs {
				if addr != topic.Launcher {
					s, err := this.storage.FindSession(topic.SessionId)
					if err != nil {
						s.state = topic.State
					}
					//this.storage.UpdateSession(&StorageSession{Id: topic.SessionId, State: topic.State})
				}
			}
		case ObjectActionDelete:

		case ObjectActionRegister:
		default:
		}
	}
	return nil
}
