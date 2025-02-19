// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"
)

func TestAccIBMAtrackerRouteBasic(t *testing.T) {
	var conf atrackerv1.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	receiveGlobalEvents := "false"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	receiveGlobalEventsUpdate := "true"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMAtrackerRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMAtrackerRouteConfigBasic(name, receiveGlobalEvents),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMAtrackerRouteExists("ibm_atracker_route.atracker_route", conf),
					resource.TestCheckResourceAttr("ibm_atracker_route.atracker_route", "name", name),
					resource.TestCheckResourceAttr("ibm_atracker_route.atracker_route", "receive_global_events", receiveGlobalEvents),
				),
			},
			{
				Config: testAccCheckIBMAtrackerRouteConfigBasic(nameUpdate, receiveGlobalEventsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_atracker_route.atracker_route", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_atracker_route.atracker_route", "receive_global_events", receiveGlobalEventsUpdate),
				),
			},
			{
				ResourceName:      "ibm_atracker_route.atracker_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMAtrackerRouteConfigBasic(name string, receiveGlobalEvents string) string {
	return fmt.Sprintf(`


		resource "ibm_atracker_target" "atracker_target" {
			name = "my-cos-target"
			target_type = "cloud_object_storage"
			cos_endpoint {
				endpoint = "s3.private.us-east.cloud-object-storage.appdomain.cloud"
				target_crn = "crn:v1:bluemix:public:cloud-object-storage:global:a/11111111111111111111111111111111:22222222-2222-2222-2222-222222222222::"
				bucket = "my-atracker-bucket"
				api_key = "xxxxxxxxxxxxxx"
			}
		}

		resource "ibm_atracker_route" "atracker_route" {
			name = "%s"
			receive_global_events = %s
			rules {
				target_ids = [ ibm_atracker_target.atracker_target.id ]
			}
		}

	`, name, receiveGlobalEvents)
}

func testAccCheckIBMAtrackerRouteExists(n string, obj atrackerv1.Route) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		atrackerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).AtrackerV1()
		if err != nil {
			return err
		}

		getRouteOptions := &atrackerv1.GetRouteOptions{}

		getRouteOptions.SetID(rs.Primary.ID)

		route, _, err := atrackerClient.GetRoute(getRouteOptions)
		if err != nil {
			return err
		}

		obj = *route
		return nil
	}
}

func testAccCheckIBMAtrackerRouteDestroy(s *terraform.State) error {
	atrackerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).AtrackerV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_atracker_route" {
			continue
		}

		getRouteOptions := &atrackerv1.GetRouteOptions{}

		getRouteOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := atrackerClient.GetRoute(getRouteOptions)

		if err == nil {
			return fmt.Errorf("Activity Tracker Route still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("[ERROR] Error checking for Activity Tracker Route (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
