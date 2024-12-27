package dev

import (
	"fmt"
	"log"
	"strings"

	"math/rand"

	"github.com/fatih/color"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/groupssettings/v1"
)

var TICK string = ("[" + color.GreenString("+") + ("]"))
var HEADING string = (color.BlueString("#"))
var TICKERROR string = ("[" + color.RedString("!") + ("]"))
var TICKINPUT string = ("[" + color.YellowString(">") + ("]"))
var SEP string = (color.BlueString("------------------------------------------------------------------------------------------------------------------------------]"))

func CreateDemoGroups(srv *admin.Service, gsrv *groupssettings.Service, domainValue string) {
	fmt.Println(TICK, "Creating 75 demo groups...")

	for i := 0; i < 75; i++ {
		groupName := fmt.Sprintf("demo-group-%d-%s@%s", i, RandomString(5), domainValue)
		groupDescription := fmt.Sprintf("Demo group %d created by Gubble", i)

		// Create the group
		group := &admin.Group{
			Email:       groupName,
			Name:        groupName,
			Description: groupDescription,
		}
		_, err := srv.Groups.Insert(group).Do()
		if err != nil {
			fmt.Printf(TICKERROR, "Error creating group %s: %v\n", groupName, err)
			continue
		}

		// Configure group settings in a single Patch request
		groupSettings := &groupssettings.Groups{
			WhoCanJoin:               "ALL_IN_DOMAIN_CAN_JOIN",
			WhoCanViewMembership:     "ALL_IN_DOMAIN_CAN_VIEW",
			WhoCanViewGroup:          "ANYONE_CAN_VIEW",
			WhoCanPostMessage:        "ALL_IN_DOMAIN_CAN_POST",
			AllowExternalMembers:     "true",
			MembersCanPostAsTheGroup: "true",
			WhoCanLeaveGroup:         "ALL_MEMBERS_CAN_LEAVE",
			WhoCanContactOwner:       "ANYONE_CAN_CONTACT",
			WhoCanDiscoverGroup:      "ANYONE_CAN_DISCOVER",
		}

		_, err = gsrv.Groups.Patch(groupName, groupSettings).Do()
		if err != nil {
			fmt.Printf(TICKERROR, "Error configuring group settings for %s: %v\n", groupName, err)
		}

		fmt.Println(TICK, "Created group:", groupName)
	}
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		//#nosec:338. Just building a random group name. Not used for anything sensitive
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func DeleteDemoGroups(srv *admin.Service, domainValue string) {
	fmt.Println(TICK, "Deleting demo groups...")

	// Get all groups in the domain
	response, err := srv.Groups.List().Domain(domainValue).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve groups in domain: %v", err)
	}

	for _, g := range response.Groups {
		if strings.HasPrefix(g.Name, "demo-group-") {
			err := srv.Groups.Delete(g.Id).Do()
			if err != nil {
				fmt.Printf(TICKERROR, "Error deleting group %s: %v\n", g.Name, err)
			} else {
				fmt.Println(TICK, "Deleted group:", g.Name)
			}
		}
	}
}
