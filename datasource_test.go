package gapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gobs/pretty"
)

const (
	createdDataSourceJSON = `{"id":1,"message":"Datasource added", "name": "test_datasource"}`
	getDataSourcesJSON    = `
		[
		  {
		    "id":1,
		    "orgId":1,
		    "name":"datasource_elastic",
		    "type":"elasticsearch",
		    "access":"proxy",
		    "url":"http://mydatasource.com",
		    "password":"",
		    "user":"",
		    "database":"grafana-dash",
		    "basicAuth":false,
		    "basicAuthUser":"",
		    "basicAuthPassword":"",
		    "isDefault":false,
		    "jsonData":null
		  }
		]
		`
)

func gapiTestTools(code int, body string) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: tr}

	url := url.URL{
		Scheme: "http",
		Host:   "my-grafana.com",
	}

	client := &Client{"my-key", url, httpClient}

	return server, client
}

func TestNewDataSource(t *testing.T) {
	server, client := gapiTestTools(200, createdDataSourceJSON)
	defer server.Close()

	ds := &DataSource{
		Name:      "foo",
		Type:      "cloudwatch",
		URL:       "http://some-url.com",
		Access:    "access",
		IsDefault: true,
		JSONData: JSONData{
			AssumeRoleArn:           "arn:aws:iam::123:role/some-role",
			AuthType:                "keys",
			CustomMetricsNamespaces: "SomeNamespace",
			DefaultRegion:           "us-east-1",
			TlsSkipVerify:           true,
		},
		SecureJSONData: SecureJSONData{
			AccessKey: "123",
			SecretKey: "456",
		},
	}

	created, err := client.NewDataSource(ds)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(created))

	if created != 1 {
		t.Error("datasource creation response should return the created datasource ID")
	}
}

func TestDataSources(t *testing.T) {
	server, client := gapiTestTools(200, getDataSourcesJSON)
	defer server.Close()

	list, err := client.DataSources()
	if err != nil {
		t.Error(err)
	}
	if len(list) == 0 {
		t.Error("Datasources parse error")
	}
}
