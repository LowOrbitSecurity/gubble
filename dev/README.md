## manage-groups.go

This code implements functions for creating and deleting demo groups in Google Workspace.

**Functions:**

- **CreateDemoGroups:** Creates 75 demo groups with randomized names and specific settings.

![image](https://github.com/user-attachments/assets/bfb2b4a7-02e3-45fc-b7e7-e9cca21e6421)


- **RandomString:** Generates a random string of specified length for group name randomization.

- **DeleteDemoGroups:** Deletes all groups with names starting with "demo-group-" within the specified domain.

![image](https://github.com/user-attachments/assets/a80d8f7f-d2ee-4e7c-bf1b-84acb47206c0)



**Usage:**

These functions are intended for development and testing purposes. They can be used to quickly populate a Google Workspace domain with demo groups for testing gubble's functionality or other scenarios.
Create Demo groups. The `-demo` and `-delete-demo` flags in the main program utilize these functions to create and delete demo groups, respectively.

```bash
sudo go run main.go  -credentials ~/Downloads/gubble.json -domain <yourdomain>.com -demo
```

Delete all demo groups
```bash
 sudo go run main.go  -credentials ~/Downloads/gubble.json -domain <yourdomain>.com -delete-demo
```
