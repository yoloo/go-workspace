package config
type AppConfig struct{
	Mycache struct{
		Password string `xml:"password,attr"`
		Db int `xml:"db,attr"`
		Address string `xml:"address,attr"`
	} `xml:"mycache"`
	Myredis struct{
		Password string `xml:"password,attr"`
		Db int `xml:"db,attr"`
		Address string `xml:"address,attr"`
	} `xml:"myredis"`
	ListenPort struct{
		Port int `xml:"port,attr"`
	} `xml:"listen_port"`
	Iplimit []struct{
		Ip string `xml:"ip,attr"`
	} `xml:"iplimit"`
}
