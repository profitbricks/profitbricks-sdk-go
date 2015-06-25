package profitbricks

import (
	"fmt"
	"testing"
)

var dcID string

func TestListDatacenters(t *testing.T) {
	shouldbe := "collection"
	want := 200
	resp := ListDatacenters()

	if resp.Type != shouldbe {
		t.Errorf(bad_type(shouldbe, resp.Type))
	}
	if resp.Resp.StatusCode != want {
		t.Errorf(bad_status(want, resp.Resp.StatusCode))
	}
}
func TestCreateDatacenter(t *testing.T) {
	want := 202
	var jason = []byte(`{
    "properties": {
        "name": "datacenter-name",
        "description": "datacenter-description",
        "location": "us/lasdev"
    }
	}`)
	resp := CreateDatacenter(jason)
	dcID = resp.Id
	if resp.Resp.StatusCode != want {
		t.Errorf(bad_status(want, resp.Resp.StatusCode))
	}
}
func TestGetDatacenter(t *testing.T) {
	shouldbe := "datacenter"
	want := 200
	
	fmt.Println(dcID)
	resp := GetDatacenter(dcID)
	if resp.Type != shouldbe {
		t.Errorf(bad_type(shouldbe, resp.Type))
	}
	if resp.Resp.StatusCode != want {
		t.Errorf(bad_status(want, resp.Resp.StatusCode))
	}
}

func TestPatchDatacenter(t *testing.T) {
	want := 202
	jason_patch := []byte(`{"name":"Renamed DC"}`)
	resp := PatchDatacenter(dcID, jason_patch)
	if resp.Resp.StatusCode != want {
		t.Errorf(bad_status(want, resp.Resp.StatusCode))
	}

}

func TestDeleteDatacenter(t *testing.T) {
	want := 202
	resp := DeleteDatacenter(dcID)
	if resp.StatusCode != want {
		t.Errorf(bad_status(want, resp.StatusCode))
	}
}
