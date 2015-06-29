package profitbricks

import "fmt"

func mkdcid() string {
	dc := CreateDatacenter([]byte(`{
    "properties": {
        "name": "GOSDK",
        "description": "datacenter-description",
        "location": "us/lasdev"
    }
	}`))

	return dc.Id
}

func mklocid() string {
	resp := ListLocations()

	locid := resp.Items[0].Id
	return locid
}

func mksrvid(srv_dcid string) string {
	var jason = []byte(`{"properties":{
						"name":"GoServer",
						"cores":4,
						"ram": 4096}
					}`)
	srv := CreateServer(srv_dcid, jason)

	return srv.Id
}

func mknic(lbal_dcid, serverid string) string {
	resp := CreateNic(lbal_dcid, serverid, []byte(`{"properties": {"name":"Original Nic","lan":1}}`))
	fmt.Println("===========================")
	fmt.Println("created a nic with id " + resp.Id)
	fmt.Println("===========================")
	return resp.Id
}
