/*
	各种配置
*/

package manager

import (
	"encoding/json"
)

// FTP代理

// HTTP代理
//type HttpProxyConfig struct {
//
//}

// HTTP 反向代理配置
type HttpReproxyConfig struct {
	LocalPort  uint   `json:"local_port"`
	ServerHost string `json:"server_host"`
	ServerPort uint   `json:"server_port"`
	Name       string `json:"name"`
}

// TCP/UDP 代理配置
type StreamProxyConfig struct {
	LocalNet   string `json:"local_net"`
	LocalPort  uint   `json:"local_port"`
	ServerHost string `json:"server_host"`
	ServerPort uint   `json:"server_port"`
}

// Shadowsocks 客户端配置
type SsClientConfig struct {
	LocalNet   string `json:"local_net"`
	LocalPort  uint   `json:"local_port"`
	ServerHost string `json:"server_host"`
	ServerPort uint   `json:"server_port"`
	ChannelNet string `json:"channel_net"`
	Password   string `json:"password"`
	Crypt      string `json:"crypt"`
}

// Shadowsocks 服务器端配置
type SsServerConfig struct {
	Port       uint   `json:"port"`
	ChannelNet string `json:"channel_net"`
	Crypt      string `json:"crypt"`
	Password   string `json:"password"`
}

type ProxyConfig struct {
	Api struct {
		Password string
	}
	Stream      []StreamProxyConfig `json:"stream"`
	HttpReproxy []HttpReproxyConfig `json:"http_reproxy"`

	SsClient *SsClientConfig  `json:"ss_client"`
	SsServer []SsServerConfig `json:"ss_server"`
}

//	将配置格式化为json字符串
//
func (this ProxyConfig) toBytes() (b []byte, err error) {
	b, err = json.Marshal(this)
	return
}

var Config ProxyConfig
