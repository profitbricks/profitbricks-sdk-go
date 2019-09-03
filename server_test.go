package profitbricks

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/suite"
)

type TestClientServer struct {
	ClientBaseSuite
}

func TestClient_Server(t *testing.T) {
	suite.Run(t, new(TestClientServer))
}
func (s *TestClientServer) TestClient_GetServer() {
	mResp := makeJsonResponse(http.StatusOK, []byte(`
{
  "id" : "{serverId}",
  "type" : "server",
  "href" : "https://api.ionos.com/cloudapi/v5/datacenters/{dataCenterId}/servers/{serverId}",
  "metadata" : {
    "createdDate" : "2000-01-01T01:00:00Z",
    "createdBy" : "[user]",
    "createdByUserId" : "[userId]",
    "etag" : "[etag]",
    "lastModifiedDate" : "2000-01-01T01:00:00Z",
    "lastModifiedBy" : "[user]",
    "lastModifiedByUserId" : "[userId]",
    "state" : "AVAILABLE"
  },
  "properties" : {
    "name" : "Server001",
    "cores" : 1,
    "ram" : 512,
    "availabilityZone" : "ZONE_1",
    "vmState" : "RUNNING",
    "bootCdrom" : null,
    "bootVolume" : {
      "id" : "{volumeId}",
      "type" : "volume",
      "href" : "https://api.ionos.com/cloudapi/v5/datacenters/{dataCenterId}/volumes/{volumeId}"
    },
    "cpuFamily" : "AMD_OPTERON"
  },
  "entities" : {
    "cdroms" : {
      "id" : "{serverId}/cdroms",
      "type" : "collection",
      "href" : "https://api.ionos.com/cloudapi/v5/datacenters/{dataCenterId}/servers/{serverId}/cdroms"
    },
    "volumes" : {
      "id" : "{serverId}/volumes",
      "type" : "collection",
      "href" : "https://api.ionos.com/cloudapi/v5/datacenters/{dataCenterId}/servers/{serverId}/volumes"
    },
    "nics" : {
      "id" : "{serverId}/nics",
      "type" : "collection",
      "href" : "https://api.ionos.com/cloudapi/v5/datacenters/{dataCenterId}/servers/{serverId}/nics"
    }
  }
}`))
	httpmock.RegisterResponder(http.MethodGet, `=~/datacenters/1/servers/2\?`,
		httpmock.ResponderFromResponse(mResp))
	srv, err := s.c.GetServer("1", "2")
	s.NoError(err)
	s.Equal("Server001", srv.Properties.Name)
}
