# sshaft

An SSH authentication factor.

# Architecture

sshaft is designed for [tomiko](https://github.com/toshokan/tomiko) and its `identity` component. Users can add their SSH keys to `identity` and choose to enable sshaft as an authentication factor if desired.

When `identity` presents an authentication challenge to a user, it records that there is an active SSH challenge for that user, and informs the user that they should perform an SSH authentication against the server configured for sshaft.

sshaft's `sshaft-keys` component is hooked into the `sshd` process as an `AuthorizedKeysCommand`. Whenever an authentication request comes in at the SSH server, `sshd` calls `sshaft-keys` to acquire the list of public components of keys that may be used to successfully authenticate. 

sshaft uses the credentials and token endpoint in its config file to acquire an access token from `tomiko` that authorizes it to retrieve the list of SSH public keys for users with active SSH challenges from `identity`. It outputs these keys in the format of an `authorized_keys` file, alongside a command override that ensures users that successfully authenticate get directed to the next sshaft component rather than to a shell on the system.

The next component, `sshaft-login`, is only called (by the command override from earlier) when a user successfully authenticates with their SSH key. This component informs `identity` that the authentication challenge has been passed, which allows the user to continue their login flow.

# Usage

sshaft is meant to be used as a containerized service with the included Dockerfile.
An sshd config for openssh is included.
You will need to supply a config file and SSH host keys.

# Config file

```json
{
  "token_endpoint": "https://idp.example.com/oauth/v1/token",
  "client_id": "my-client-id",
  "client_secret": "shhhhh",
  "token_scope": "relevant_scope1 relevant_scope2",
  "mfa_list_endpoint": "https://authc.example.com/api/v1/mfa/ssh/keys",
  "mfa_accept_endpoint": "https://authc.example.com/api/v1/mfa/ssh/accept",
  "login_path": "/sshaft/login"
}

```

In a standard `tomiko` deployment, `idp.example.com` is replaced with `tomiko`'s hostname, `authc.example.com` is replaced with `identity`'s hostname, and the scope will be `tomiko::mfa:rw`.

# License

sshaft is distributed under the terms of the GPLv3 license. 
