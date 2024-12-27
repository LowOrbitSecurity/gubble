## manage-groups.go

This code implements functions for creating and deleting demo groups in Google Workspace.

**Functions:**

- **CreateDemoGroups:** Creates 75 demo groups with randomized names and specific settings.

- **RandomString:** Generates a random string of specified length for group name randomization.

- **DeleteDemoGroups:** Deletes all groups with names starting with "demo-group-" within the specified domain.

**Usage:**

These functions are intended for development and testing purposes. They can be used to quickly populate a Google Workspace domain with demo groups for testing gubble's functionality or other scenarios.

**Note:**

The `-demo` and `-delete-demo` flags in the main program utilize these functions to create and delete demo groups, respectively.
