# gubble

**gubble** is a tool designed to audit Google Workspace group settings. It helps identify potential security risks associated with group configurations by analyzing settings such as who can join, view membership, post messages, view conversations, and more.

![image](https://github.com/user-attachments/assets/b8a4c02b-91d9-42c9-9fdd-ac33adeb5cf3)

## Setup

gubble requires the following permissions to function correctly:

- **Admin SDK API:**
  - `admin.directory.group.readonly`:  This scope allows gubble to read the list of groups in your Google Workspace domain and their basic information.
- **Groups Settings API:**
  - `apps.groups.settings`: This scope allows gubble to read the detailed settings of each group, such as who can join, post messages, and view membership.

These permissions are granted during the OAuth client configuration and consent screen setup. Make sure to enable the necessary APIs and add the required scopes to your OAuth client.

### 1. Enable APIs & Services

This section enables the necessary APIs (Admin SDK and Groups Settings API) in your Google Cloud console to allow gubble to access and analyze your Google Workspace group settings.

1. Go to **APIs & Services -> Enabled APIs & Services** in your Google Cloud console.
2. Click **Enable APIS and Services**.
3. Search for and enable the following APIs:
    - **Admin SDK API**
    - **Groups Settings API**

### 2. Create OAuth Client

This section guides you through creating an OAuth client ID for gubble, which allows the tool to securely access your Google Workspace data with your authorization. This involves creating credentials in your Google Cloud console, downloading the configuration file, and renaming it to gubble.json for gubble to use.

1. Go to **APIs & Services -> Credentials** in your Google Cloud console.
2. Click **Create Credentials -> OAuth client ID**.
3. Select **Application Type:** Desktop App.
4. Enter **Name:** gubble.
5. Click **Create**.
6. Download the generated JSON file.
7. Move the downloaded JSON file to a suitable location and rename it to `gubble.json`:

    ```bash
    mv ~/Downloads/client_secret_* ~/Downloads/gubble.json
    ``` 

    > ⚠️ `gubble.json` contains sensitive data, KEEP IT SAFE. ⚠️

### 3. OAuth Consent Screen Setup

This section sets up the OAuth consent screen, which is what users in your Google Workspace domain will see when authorizing gubble to access their group information. This involves configuring the consent screen with basic app information and specifying the required permissions (scopes).

1. Go to **APIs & Services -> OAuth Consent Screen** in your Google Cloud console.
2. Select **User Type:** Internal.
3. Click **Create**.
4. Fill out the required information (App name, logo, etc.).
5. In the **Scopes** section, add the same scopes as in the previous step:

    ```markdown
    https://www.googleapis.com/auth/admin.directory.group.readonly
    https://www.googleapis.com/auth/apps.groups.settings
    ```

    > If you do not see these scopes, please ensure the Admin SDK API and Groups Settings API have been enabled.

### 4. Add Scopes to OAuth Client

This section describes the process of setting up the OAuth consent screen for your application. This screen is displayed to users when they are asked to authorize your application to access their Google Workspace group information. The setup involves providing basic application details and defining the necessary permissions (scopes) that your application requires.

1. Edit the OAuth client you just created.
2. Go to the **Scopes** tab.
3. Add the following scopes:

    | API                | Scope                                                        | Description                                  |
    | ------------------ | ----------------------------------------------------------- | -------------------------------------------- |
    | Admin SDK API      | `https://www.googleapis.com/auth/admin.directory.group.readonly` | View groups on your domain                   |
    | Groups Settings API | `https://www.googleapis.com/auth/apps.groups.settings`      | View and manage the settings of a G Suite group |

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

![image](https://github.com/user-attachments/assets/afcd3b2f-ed03-410a-a4fa-aef584efea72)

- Allow gubble access (Notice gubble is only getting access to the scopes provided during setup).

![image](https://github.com/user-attachments/assets/183b318e-99ce-47fb-97ca-8ba00e6d21c2)

- After authentication, you may close the window and return to your terminal.

## Example Output
![image](https://github.com/user-attachments/assets/e5113464-fd36-4c09-936b-2e884b8a7013)
