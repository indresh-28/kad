package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSyncApp(t *testing.T) {

	assert := require.New(t)

	var wantConfig AppConfig
	err := yaml.Unmarshal([]byte(content), &wantConfig)
	assert.Nil(err)

	store := cassandra.NewCassandraStore(logging.NewLogger(), nil)
	err = store.Connect([]string{"localhost:9042"}, "cassandra", "cassandra")
	assert.Nil(err)

	agent, err := NewAgent(logging.NewLogger(), nil, store)
	assert.Nil(err)

	err = agent.syncApp(context.TODO(), &agentpb.SyncAppRequest{Payload: []byte(content)})
	assert.Nil(err)

	gotConfig, err := getAppsByName(agent.store.Session(), "signoz")
	assert.Nil(err)

	reflect.DeepEqual(wantConfig, gotConfig)

}

var content = `
Name: "signoz"
ChartName: "signoz/signoz"
RepoName: "signoz"
RepoURL: "https://charts.signoz.io"
Namespace: "observability"
ReleaseName: "signoz"
Version: "0.14.0"
CreateNamespace: true
Override:
  Values:
    clickhouse:
      password": admin
    frontend:
      ingress:
        enabled": true
        hosts:
        - host: "signoz.{{.DomainName}}"
          paths:
          - path: /
            pathType: ImplementationSpecific
            port: 3301
        tls:
          hosts:
          - "signoz.{{.DomainName}}"
          secretName: cert-signoz
      annotations:
        cert-manager.io/cluster-issuer": letsencrypt-prod-cluster
        nginx.ingress.kubernetes.io/backend-protocol": HTTPS
        nginx.ingress.kubernetes.io/force-ssl-redirect": true
        nginx.ingress.kubernetes.io/ssl-redirect": true
`
