# matrix-server-version-check
This is a check plugin for the [matrix server](https://matrix.org/) version written in go. This check compares the current versions that are available on the matrix server to a hardcoded array with versions numbers.

# Usage
- -u specify the URL for the GET request (is required, for example: "https://<your-homeserver-url/_matrix/client/versions")
