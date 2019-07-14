package mongodbatlas

import (
	"reflect"
	"testing"
)

func testResourceVpcPeeringConnectionResourceV0_aws() map[string]interface{} {
	return map[string]interface{}{
		"aws_account_id":         "123456789",
		"connection_id":          "pcx-123wefds43erg34ter",
		"container_id":           "284nf7ek37g73jr7tj4jr8ei",
		"error_state_name":       "",
		"group":                  "812nf72jf82j72j72hejw8yr",
		"id":                     "7s6b3r7tdsgf7t3igsdft3fu",
		"identifier":             "7s6b3r7tdsgf7t3igsdft3fu",
		"route_table_cidr_block": "172.20.0.0/16",
		"status_name":            "AVAILABLE",
		"vpc_id":                 "vpc-12345678",
	}
}

func testResourceVpcPeeringConnectionResourceV1_aws() map[string]interface{} {
	v0 := testResourceVpcPeeringConnectionResourceV0_aws()
	return map[string]interface{}{
		"aws_account_id":         v0["aws_account_id"],
		"connection_id":          v0["connection_id"],
		"container_id":           v0["container_id"],
		"error_state_name":       v0["error_state_name"],
		"group":                  v0["group"],
		"id":                     v0["id"],
		"identifier":             v0["identifier"],
		"route_table_cidr_block": v0["route_table_cidr_block"],
		"status_name":            v0["status_name"],
		"vpc_id":                 v0["vpc_id"],
		"provider_name":          "AWS",
	}
}

func testResourceVpcPeeringConnectionResourceV0_gce() map[string]interface{} {
	return map[string]interface{}{
		"container_id":   "507f1f77bcf86cd799439011",
		"error_message":  "",
		"status":         "ADDING_PEER",
		"gcp_project_id": "my-sample-project-191923",
		"network_name":   "test1",
		"provider_name":  "GCP",
	}
}

func TestResourceAwsKinesisStreamStateUpgradeV0_aws(t *testing.T) {
	expected := testResourceVpcPeeringConnectionResourceV1_aws()
	actual, err := resourceVpcPeeringConnectionStateUpgradeV0(testResourceVpcPeeringConnectionResourceV0_aws(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestResourceAwsKinesisStreamStateUpgradeV0_gce(t *testing.T) {
	expected := testResourceVpcPeeringConnectionResourceV0_gce()
	actual, err := resourceVpcPeeringConnectionStateUpgradeV0(testResourceVpcPeeringConnectionResourceV0_gce(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
