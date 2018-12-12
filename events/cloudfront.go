package events

type CloudFrontConfig struct {
	DistributionDomainName string `json:"distributionDomainName"`
	DistributionID         string `json:"distributionId"`
	EventType              string `json:"eventType"`
	RequestID              string `json:"requestId"`
}

type CloudFrontRequest struct {
	Records []CloudFrontEventRequestRecord `json:"Records"`
}

type CloudFrontEventRequestRecord struct {
	Cf struct {
		Config  CloudFrontConfig `json:"config"`
		Request struct {
			Body struct {
				Action         string `json:"action"`
				Data           string `json:"data"`
				Encoding       string `json:"encoding"`
				InputTruncated bool   `json:"inputTruncated"`
			} `json:"body"`
			ClientIP    string `json:"clientIp"`
			QueryString string `json:"querystring"`
			URI         string `json:"uri"`
			Method      string `json:"method"`
			Headers     map[string][]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
			Origin map[string]struct {
				AuthMethod    string `json:"authMethod"`
				CustomHeaders map[string][]struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"customHeaders"`
				DomainName       string   `json:"domainName"`
				KeepaliveTimeout int      `json:"keepaliveTimeout"`
				Path             string   `json:"path"`
				Region           string   `json:"region"`
				Port             int      `json:"port"`
				Protocol         string   `json:"protocol"`
				ReadTimeout      int      `json:"readTimeout"`
				SslProtocols     []string `json:"sslProtocols"`
			} `json:"origin"`
		} `json:"request"`
	} `json:"cf"`
}

type CloudFrontResponse struct {
	Records []CloudFrontEventResponseRecord `json:"records"`
}

type CloudFrontEventResponseRecord struct {
	Cf struct {
		Config  CloudFrontConfig `json:"config"`
		Request struct {
			ClientIP    string `json:"clientIp"`
			QueryString string `json:"querystring"`
			URI         string `json:"uri"`
			Method      string `json:"method"`
			Headers     map[string][]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"request"`
		Response struct {
			Status            string `json:"status"`
			StatusDescription string `json:"statusDescription"`
			Headers           map[string][]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"response"`
	} `json:"cf"`
}
