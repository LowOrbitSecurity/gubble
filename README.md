# gubble

**gubble** is a tool designed to audit Google Workspace group settings. It helps identify potential security risks associated with group configurations by analyzing settings such as who can join, view membership, post messages, view conversations, and more.

![image](https://github.com/user-attachments/assets/67a8c56f-d7e9-498a-8494-5d35f98789a3)


Google Groups is a service that allows users to create and manage online discussion groups and email lists. It provides a platform for communities and discussions, with features like email lists, web forums, Q&A forums, and collaborative inboxes. However, With great power comes great responsbility. A group with misconfigured permissions can lead to excessive data exposure and privielge escalation if the risks are not understood. 



During penetration tests, testers often aim to identify groups that can be joined, groups that allow for many users to read privileged conversations, or groups configured in such a way that would make internal phishing easier. Gubble is a tool that aims to automate that process.


|  METHOD NAME           |          Risky Permission            |                         Notes                                                                    |
|------------------------|--------------------------------------|--------------------------------------------------------------------------------------------------|
|whoCanJoin              |ANYONE_CAN_JOIN  ALL_IN_DOMAIN_CAN_JOIN|Anyone in the domain can join the group. This could be used for privielge escalation.             |
|whoCanViewMembership    |ALL_IN_DOMAIN_CAN_VIEW                |Only bad if you have "secret" projects                                                            |
|whoCanViewGroup         |ANYONE_CAN_VIEW  ALL_IN_DOMAIN_CAN_VIEW|This means who can read conversations                                                             |
|allowExternalMembers    |TRUE                                  |External Identities can be added to the group.                                                    |
|whoCanPostMessage       |ALL_IN_DOMAIN_CAN_POST  ANYONE_CAN_POST|This can be utilized for phishing.                                                                |
|membersCanPostAsTheGroup|TRUE                                  |This can be abused for phishing.                                                                  |
|whoCanLeaveGroup        |NONE_CAN_LEAVE                        |This can be used as a honeypot. Make a juicy group name and alert on users joining it since they can't leave|
|whoCanContactOwner      |ANYONE_CAN_CONTACT                    |                                                                                                  |
|whoCanDiscoverGroup     |ANYONE_CAN_DISCOVER                   |                                                                                                  |
|defaultSender           |GROUP                                 |                                                                                                  |



## Setup

gubble requires the following permissions to function correctly:

- **Admin SDK API:**
  - `admin.directory.group.readonly`:  This scope allows gubble to read the list of groups in your Google Workspace domain and their basic information.
- **Groups Settings API:**
  - `apps.groups.settings`: This scope allows gubble to read the detailed settings of each group, such as who can join, post messages, and view membership.

These permissions are granted during the OAuth client configuration and consent screen setup. Make sure to enable the necessary APIs and add the required scopes to your OAuth client.

### 1. Create a GCP Project
Create a GCP project that can be used to configure the correction APIs and permissions needed for gubble to run.

![image](https://github.com/user-attachments/assets/b2a5dcf2-1a67-4e82-9484-65c385876b03)

After creation, switch to the newly created project 

### 2. Enable APIs & Services

This section enables the necessary APIs (Admin SDK and Groups Settings API) in your Google Cloud console to allow gubble to access and analyze your Google Workspace group settings.

1. Go to **APIs & Services -> Enabled APIs & Services** in your Google Cloud console.
2. Click **Enable APIS and Services**.
3. Search for and enable the following APIs:
    - **Admin SDK API**
    - **Groups Settings API**
  
  ![image](https://github.com/user-attachments/assets/98fadf37-60d7-41ca-99cc-81b08a736b53)


### 3. Create OAuth Consent Screen

This section sets up the OAuth consent screen, which is what users in your Google Workspace domain will see when authorizing gubble to access their group information.

1. Go to **APIs & Services -> OAuth Consent Screen** in your Google Cloud console. (If you see "Go To New Experience", click it.)
2. Select **Get Started**
3. App Information:
  - App Name: `gubble`
  - User Support Email: `your_email@example.com`
5. Contact Information:
  - `your_email@example.com`
4. App Information:
  - Internal
5. Click **Create**

### 4. Add Scopes to OAuth Client

This section describes the process of setting up the OAuth consent screen for your application. This screen is displayed to users when they are asked to authorize your application to access their Google Workspace group information. The setup involves providing basic application details and defining the necessary permissions (scopes) that your application requires.

1. Edit the OAuth client you just created.
2. Go to the **Scopes** tab.
3. Add the following scopes:

    | API                | Scope                                                        | Description                                  |
    | ------------------ | ----------------------------------------------------------- | -------------------------------------------- |
    | Admin SDK API      | `https://www.googleapis.com/auth/admin.directory.group.readonly` | View groups on your domain                   |
    | Groups Settings API | `https://www.googleapis.com/auth/apps.groups.settings`      | View and manage the settings of a G Suite group |


    > If you do not see these scopes, please ensure the Admin SDK API and Groups Settings API have been enabled.

![image](https://github.com/user-attachments/assets/8730ee80-37fc-4316-80b3-278bdb9fe0fb)

4. Click **Save**


### 5. Create OAuth Client

This section guides you through creating an OAuth client ID for gubble, which allows the tool to securely access your Google Workspace data with your authorization. This involves creating credentials in your Google Cloud console, downloading the configuration file, and renaming it to gubble.json for gubble to use.

Navigate to **APIs & Services** -> **OAuth Consent Screen** (Click **Go To New Experience** if prompted)
1. Click **Clients** -> **Create Client** (You may have to refresh the page)
2. Application type:
  - Desktop app
3. Name:
  - `gubble`
4. Click **Create**
5. Click the **Download Arrow**
6. Click **Download JSON**
7. Rename the file to gubble.json. You will need to specify this credential with the `-credentials ~/Downloads/gubble.json` later.
    ```bash
    mv ~/Downloads/client_secret_* ~/Downloads/gubble.json
    ``` 

    > ⚠️ `gubble.json` contains sensitive data, KEEP IT SAFE. ⚠️

## Usage

**Flags:**

- `-credentials`: Path to the `gubble.json` credentials file (required).
- `-domain`: Your Google Workspace domain (required).
- `-log`: Location to save the log file (optional).
- `-verbose`: Verbose mode. Prints all group settings, even those not considered a risk (optional)# Usage

```bash
sudo gubble -credentials /path/to/credentials.json -domain yourdomain.com
```

- Click to Authenticate to GCP.

![image](https://github.com/user-attachments/assets/df60b98d-e3df-46c9-9813-3f0bec78ee11)

- Allow gubble access (Notice gubble is only getting access to the scopes provided during setup).

![image](https://github.com/user-attachments/assets/3dd11f7b-d4af-4a34-b0ca-1aea85d95453)

- After authentication, you may close the window and return to your terminal.

## Example Output

![image](https://github.com/user-attachments/assets/219ac2c9-db2a-4e39-aaf4-dbc1bea9fc02)

### Dev Tools

**Dev Flags:** (You probably don't need these)
These flags require additional setup. Please see [the dev readme](https://github.com/LowOrbitSecurity/gubble/tree/main/dev)

- `-demo`: Demo mode. This will create 75 google groups.
- `-delete-demo`: Deletes the groups made in the demo.

If you're doing research with google groups, I've added tooling I've created to help test gubble. Please see [the dev readme](https://github.com/LowOrbitSecurity/gubble/tree/main/dev)
