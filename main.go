package main

// https://developers.google.com/admin-sdk/groups-settings/v1/reference/groups#json
import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"gubble/dev"
	auth "gubble/utils"
	"log"
	"os"

	"github.com/fatih/color"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/groupssettings/v1"
	"google.golang.org/api/option"
)

var TICK string = ("[" + color.GreenString("+") + ("]"))
var HEADING string = (color.BlueString("#"))
var TICKERROR string = ("[" + color.RedString("!") + ("]"))
var TICKINPUT string = ("[" + color.YellowString(">") + ("]"))
var SEP string = (color.BlueString("------------------------------------------------------------------------------------------------------------------------------]"))
var credentialsFile string
var domainValue string
var logLocation string
var verbose bool
var demoMode bool
var deleteDemoMode bool
var banner = (`
⠀⠀⠀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⠀⠀⠀⠀⠀⠀⠀⠀⣠⣦⡀⠀⠀⠀
⠀⠀⠛⣿⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠛⣿⠛⠀⠀⠀⠀⠀⡀⠺⣿⣿⠟⢀⡀⠀
⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣦⠈⠁⣴⣿⣿⡦
⠀⠀⠀⠀⠀ gubble     ⠀⣠⣦⡈⠻⠟⢁⣴⣦⡈⠻⠋⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⡀⠺⣿⣿⠟⢀⡀⠻⣿⡿⠋⠀⠀⠀
⠀⣠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣶⡿⠿⣿⣦⡈⠁⣴⣿⣿⡦⠈⠀⠀⠀⠀⠀
⠲⣿⠷⠂⠀⠀⠀⠀⠀⠀⢀⣴⡿⠋⣠⣦⡈⠻⣿⣦⡈⠻⠋⠀⠀⠀⠀⠀⠀⠀
⠀⠈⠀⠀⠀⠀⠀⠀⠀⠰⣿⣿⡀⠺⣿⣿⣿⡦⠈⣻⣿⡦⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣠⣦⡈⠻⣿⣦⡈⠻⠋⣠⣾⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⡀⠺⣿⣿⠟⢀⡈⠻⣿⣶⣾⡿⠋⣠⣦⡀⠀⢀⣠⣤⣀⡀⠀⠀
⠀⠀⠀⠀⣠⣾⣿⣦⠈⠁⣴⣿⣿⡦⠈⠛⠋⠀⠀⠈⠛⢁⣴⣿⣿⡿⠋⠀⠀⠀
⠀⠀⣠⣦⡈⠻⠟⢁⣴⣦⡈⠻⠋⠀⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣏⠀⠀⠀⠀⠀
⠀⠺⣿⣿⠟⢀⡀⠻⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠰⣿⡿⠛⠁⠙⣷⣶⣦⠀⠀
⠀⠀⠈⠁⣴⣿⣿⡦⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠋⠀⠀⠀⠀⠻⠿⠟⠀⠀
⠀⠀⠀⠀⠈⠻⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`)

type GroupSettings struct {
	Name                     string `json:"name"`
	Description              string `json:"Description"`
	WhoCanJoin               string `json:"whoCanJoin"`
	WhoCanViewMembership     string `json:"whoCanViewMembership"`
	WhoCanViewGroup          string `json:"whoCanViewGroup"`
	WhoCanPostMessage        string `json:"whoCanPostMessage"`
	AllowExternalMembers     string `json:"allowExternalMembers"`
	MembersCanPostAsTheGroup string `json:"membersCanPostAsTheGroup"`
	WhoCanLeaveGroup         string `json:"WhoCanLeaveGroup"`
	WhoCanContactOwner       string `json:"whoCanContactOwner"`
	WhoCanDiscoverGroup      string `json:"whoCanDiscoverGroup"`
	DefaultSender            string `json:"defaultSender"`
}

func main() {
	fmt.Println(banner)
	fmt.Println(SEP)
	// Parse arguments
	flag.StringVar(&credentialsFile, "credentials", "", "Path to the credentials JSON file")
	flag.StringVar(&domainValue, "domain", "", "Domain. IE: yourcompany.com")
	flag.StringVar(&logLocation, "log", "", "Location to save logfile to. IE: ./groups.log")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode. Prints information even if it's not a risk")
	flag.BoolVar(&demoMode, "demo", false, "Demo mode. Creates 75 random groups for testing.")
	flag.BoolVar(&deleteDemoMode, "delete-demo", false, "Delete demo mode. Deletes all groups created by the -demo flag.")
	flag.Parse()

	// If credential file is not provided, exit and print usage
	if credentialsFile == "" || domainValue == "" {
		fmt.Println("Error: -credentials and -domain flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// #nosec:G304 -- Working as intended. This would allow for arbitrary file read if gubble were
	// allowed to be run by untrused users froma webapp wrapper.
	cred, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	var config *oauth2.Config
	if demoMode || deleteDemoMode {
		// Request both read and write scopes for demo mode and deleteDemoMode
		config, err = google.ConfigFromJSON(cred, admin.AdminDirectoryGroupScope, groupssettings.AppsGroupsSettingsScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
	} else {
		// Request only read-only scope for normal mode
		config, err = google.ConfigFromJSON(cred, admin.AdminDirectoryGroupReadonlyScope, groupssettings.AppsGroupsSettingsScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
	}
	client := auth.GetClient(config)
	fmt.Println(TICK, "Authentication Successful... Gathering groups. This may take a moment depending on the amount of groups...")
	// Init service client object for use with the API
	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	// Client a new client for interacting with admin API
	gsrv, err := groupssettings.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve groupssettings Client %v", err)
	}

	// Request the list of groups from the Admin SDK API. Recieve a list of groups back.
	response, err := srv.Groups.List().Domain(domainValue).Do()
	if err != nil {
		log.Println(color.RedString("Did you specify the correct domain with the -d flag?"))
		log.Fatalf("Unable to retrieve groups in domain: %v", err)
	}
	allGroupSettings := []GroupSettings{}

	if len(response.Groups) == 0 {
		fmt.Print("No groups found.\n")
	} else {

		for _, g := range response.Groups {
			gs, err := gsrv.Groups.Get(g.Email).Do()
			if err != nil {
				fmt.Printf("Error fetching settings for %s: %v\n", g.Email, err)
			} else {
				// Store the settings in a GroupSettings struct
				groupSettings := GroupSettings{
					Name:                     gs.Name,
					Description:              gs.Description,
					WhoCanJoin:               gs.WhoCanJoin,
					WhoCanViewMembership:     gs.WhoCanViewMembership,
					WhoCanViewGroup:          gs.WhoCanViewGroup,
					WhoCanPostMessage:        gs.WhoCanPostMessage,
					WhoCanLeaveGroup:         gs.WhoCanLeaveGroup,
					AllowExternalMembers:     gs.AllowExternalMembers,
					MembersCanPostAsTheGroup: gs.MembersCanPostAsTheGroup,
					WhoCanContactOwner:       gs.WhoCanContactOwner,
					WhoCanDiscoverGroup:      gs.WhoCanDiscoverGroup,
					//                    DefaultSender:    gs.defaultSender,
				}

				allGroupSettings = append(allGroupSettings, groupSettings)
			}
		}

		// Print the collected group settings
		for _, settings := range allGroupSettings {
			printGroupSettings(settings)
		}
	}

	if logLocation != "" {
		writelog(allGroupSettings)
	}
	if demoMode {
		dev.CreateDemoGroups(srv, gsrv, domainValue)
	}
	if deleteDemoMode {
		dev.DeleteDemoGroups(srv, domainValue)
	}
}
func writelog(allGroupSettings []GroupSettings) {
	fmt.Println(TICK, "Creating log file at", logLocation)
	// Create CSV file for logging
	// #nosec:G304 -- Working as intended. This would allow for arbitrary file overwrite if gubble were
	// allowed to be run by untrused users froma webapp wrapper.
	file, err := os.Create(logLocation)
	if err != nil {
		log.Fatalf("Unable to create CSV file: %v", err)
	}
	defer file.Close()

	// Write CSV header
	// BUG: Not sure why but the csv this creates can't be imported to google sheets.
	writer := csv.NewWriter(file)
	header := []string{
		"Name", "Description", "WhoCanJoin", "WhoCanViewMembership", "WhoCanViewGroup",
		"WhoCanPostMessage", "AllowExternalMembers", "MembersCanPostAsTheGroup",
		"WhoCanLeaveGroup", "WhoCanContactOwner", "WhoCanDiscoverGroup",
	}
	err = writer.Write(header)
	writer.Comma = '\t'   // Set the delimiter to comma
	writer.UseCRLF = true // Use CRLF line endings
	if err != nil {
		log.Fatalf("Unable to write CSV header: %v", err)
	}

	// Log group settings to CSV
	for _, settings := range allGroupSettings {
		record := []string{
			settings.Name, settings.Description, settings.WhoCanJoin, settings.WhoCanViewMembership,
			settings.WhoCanViewGroup, settings.WhoCanPostMessage, settings.AllowExternalMembers,
			settings.MembersCanPostAsTheGroup, settings.WhoCanLeaveGroup, settings.WhoCanContactOwner,
			settings.WhoCanDiscoverGroup,
		}
		err = writer.Write(record)
		if err != nil {
			log.Fatalf("Unable to write to CSV: %v", err)
		}
	}
	writer.Flush() // Flush any buffered data to the file
}

func riskDescription(label, riskLevel string, riskDescriptions map[string]string) string {
	if description, ok := riskDescriptions[label]; ok && riskLevel != "" {
		return description
	}
	return ""
}

func printGroupSettings(settings GroupSettings) {
	// Define colors for different risk levels because we're fancy
	var (
		highRiskColor   = color.RedString
		mediumRiskColor = color.YellowString
	)
	const riskPadding = 60

	riskDescriptions := map[string]string{
		"Join":                 "Risk: Anyone in the domain can join this group.",
		"View Membership":      "Potential Risk: All members of the domain can view who is in this group.",
		"View Conversations":   "Risk: Anyone can view conversations in this group.",
		"External Members":     "Risk: External users can be added to this group.",
		"Post Messages":        "Phishing Risk: Anyone in the domain can post to this group.",
		"Member Post As Group": "Phishing Risk: Users can send emails as the group (impersonation).",
		"Leave":                "Red Team Risk: This may be a honeypot, you can't leave once you join.",
		"Contact Owner":        "Phishing Risk: Anyone can email the owner of this group.",
		"Discover Group":       "Risk: TBD. Allowing anyone to discover the group sounds bad but idk what it does yet.",
	}

	fmt.Println(HEADING, color.BlueString(settings.Name), ":", settings.Description)
	fmt.Println(SEP)

	// Format and print group settings
	printSetting := func(label, value, riskLevel string) {
		switch riskLevel {
		case "high":
			fmt.Printf("\n%-*s %s", riskPadding, label+": "+highRiskColor(value), riskDescription(label, riskLevel, riskDescriptions))
		case "medium":
			fmt.Printf("\n%-*s %s", riskPadding, label+": "+mediumRiskColor(value), riskDescription(label, riskLevel, riskDescriptions))
		default:
			if verbose {
				fmt.Printf("\n%-*s %s", riskPadding, label+": "+value, riskDescription(label, riskLevel, riskDescriptions))
			}
		}
	}

	// Who can join
	joinRisk := ""
	if settings.WhoCanJoin == "ANYONE_CAN_JOIN" || settings.WhoCanJoin == "ALL_IN_DOMAIN_CAN_JOIN" {
		joinRisk = "high"
	}
	printSetting("Join", settings.WhoCanJoin, joinRisk)

	// View Membership
	viewMembershipRisk := ""
	if settings.WhoCanViewMembership == "ALL_IN_DOMAIN_CAN_VIEW" {
		viewMembershipRisk = "medium"
	}
	printSetting("View Membership", settings.WhoCanViewMembership, viewMembershipRisk)

	// Who Can View Group (Conversations)
	viewGroupRisk := ""
	if settings.WhoCanViewGroup == "ANYONE_CAN_VIEW" || settings.WhoCanViewGroup == "ALL_IN_DOMAIN_CAN_VIEW" {
		viewGroupRisk = "high"
	}
	printSetting("View Conversations", settings.WhoCanViewGroup, viewGroupRisk)

	// Allow External Members
	externalMembersRisk := ""
	if settings.AllowExternalMembers == "true" {
		externalMembersRisk = "high"
	}
	printSetting("External Members", settings.AllowExternalMembers, externalMembersRisk)

	// Posting Messages
	postMessageRisk := ""
	if settings.WhoCanPostMessage == "ALL_IN_DOMAIN_CAN_POST" {
		postMessageRisk = "medium"
	} else if settings.WhoCanPostMessage == "ANYONE_CAN_POST" {
		postMessageRisk = "medium"
	}
	printSetting("Post Messages", settings.WhoCanPostMessage, postMessageRisk)

	// Post As Group
	postAsGroupRisk := ""
	if settings.MembersCanPostAsTheGroup == "true" {
		postAsGroupRisk = "medium"
	}
	printSetting("Member Post As Group", settings.MembersCanPostAsTheGroup, postAsGroupRisk)

	// Can you leave?
	leaveRisk := ""
	if settings.WhoCanLeaveGroup == "NONE_CAN_LEAVE" {
		leaveRisk = "medium"
	}
	printSetting("Leave", settings.WhoCanLeaveGroup, leaveRisk)

	// Can you contact owner?
	contactOwnerRisk := ""
	if settings.WhoCanContactOwner == "ANYONE_CAN_CONTACT" {
		contactOwnerRisk = "medium"
	}
	printSetting("Contact Owner", settings.WhoCanContactOwner, contactOwnerRisk)

	// Who can discover?
	discoverGroupRisk := ""
	if settings.WhoCanDiscoverGroup == "ANYONE_CAN_DISCOVER" {
		discoverGroupRisk = "high"
	}
	printSetting("Discover Group", settings.WhoCanDiscoverGroup, discoverGroupRisk)

	fmt.Printf("\n\n")
}
