# Confluence license-checker

`license-checker` is tool that simplifies the task of restarting Confluence after the Dogu setup. The tool supplies to commands:

- `test-setup`
- `watch`

The confluence setup uses an expired license for implementing the setup routine. When the dogu finishes setting up, Confluence greets the user with a screen to replace the outdated license with a valid one, following with a reboot. In highly guarded CES instances it may be a tough job to trigger a Dogu restart.

In that way, `license-checker` is a tool that belongs to the setup phase. It does not play any role during the production-use of the Confluence dogu.

The idea behind this tool is rather trivial and embeds into the setup workflow:

1. Confluence must be set up (usually with an invalid setup license)
1. Set up finishes
1. `license-checker test-setup` positively tests for a setup license
    1. `license-checker watch` is started in the background
1. Confluence starts up and urges to replace the invalid license
    1. An administrator adds a valid production license
1. `license-checker watch` detects a license change and restarts Confluence
1. `license-checker test-setup` negatively tests for a setup license
    1. `license-checker watch` will not be started
    1. Confluence starts in regular fashion 

`license-checker test-setup` checks if an expired license is currently configured. If so, `license-checker watch` should be started to watch for license changes. This command must be provided with an shell command to be executed.

Now, when the administrator adds a valid production license, the license-checker recognizes the change (compared to the earlier setup license) and executes the provided shell command which restarts Confluence.

Even when the dogu is restarted, `license-checker test-setup` will recognize the production license, avoiding to start the `license-checker watch` routine.