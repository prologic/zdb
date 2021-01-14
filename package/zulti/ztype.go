package zulti

import (
	// "fmt"
)

type TChannelResponse struct {
	Success string
	Message string
	Data []map[string]interface{}
}

type TCatchupResponse struct {
	// Success string
	Message string
	Data map[string]interface{}
}

type TConfig struct {
	ClickhouseType string
	ClickhouseURL  string

	ServerIPDev           string
	ServerPortDev					string

	ServerIP              string
	ServerPort            string	

	ServerStlPort         string
	ServerStlKeyFullchain string
	ServerStlKeySecret    string

	VodDomainAllow    []interface{}
	VodDomainExtGuard []interface{}

	VodDomain     string
	VodRootFolder string
	VodPathPrefix string

	VodRunGGCDN       string
	VodDomainCDNRetry string
	VodDomainCDN      string
	VodBrowserAgentGG string
	CustomerID        string

	VodRetryCdnOriginRate int64
	VodRequestTimeout int64
	SizeStartGGProxy  int64
	TimeInterval      int64

	// Need calculate
	VodDomainPrefixCDN   string
	VodDomainPrefixGG    string
	VodDomainPrefixGGCDN string

	// Non parse field
	TimeCount                int64          `json:"rxnonfield,omitempty"`
	MapClientFileRequestTime *ConcurrentMap `json:"rxnonfield,omitempty"`
	MapReferer               *ConcurrentMap `json:"rxnonfield,omitempty"`
	MapFFMPEG                *ConcurrentMap `json:"rxnonfield,omitempty"`
	MapM3U8                  *ConcurrentMap `json:"rxnonfield,omitempty"`	

	// Live
	LivePathPrefix        string
	LiveRootFolder        string
  APILiveChannels       string
  APILiveCatchups				string

	// Enums
	EVODLOG               string `json:"rxnonfield,omitempty"`
	EVODLOGRequestPrevent string `json:"rxnonfield,omitempty"`
	EVODLOGRequestOrigin  string `json:"rxnonfield,omitempty"`
	EVODLOGRequestAll     string `json:"rxnonfield,omitempty"`
}

func TConfigInit() TConfig {
	result := TConfig{
		TimeCount:                int64(0),
		MapClientFileRequestTime: NewConcurrentMap(),
		MapReferer:               NewConcurrentMap(),
		MapFFMPEG: 								NewConcurrentMap(),
		MapM3U8:                  NewConcurrentMap(),

		EVODLOG:               "VODREQUEST",
		EVODLOGRequestPrevent: "REQUEST_PREVENT",
		EVODLOGRequestOrigin:  "REQUEST_ORIGIN",
		EVODLOGRequestAll:     "REQUEST_ALL",
	}

	f := func() {
		result.MapM3U8.Load("MapM3U8.gob")		
		result.MapReferer.Load("MapReferer.gob")
	}
	go f()

	return result
}

func (config *TConfig) TConfigSave() {
	// Get path
	config.MapM3U8.Save("MapM3U8.gob")
	config.MapReferer.Save("MapReferer.gob")
}
