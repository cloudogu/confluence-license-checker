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

---
## What is the Cloudogu EcoSystem?
The Cloudogu EcoSystem is an open platform, which lets you choose how and where your team creates great software. Each service or tool is delivered as a Dogu, a Docker container. Each Dogu can easily be integrated in your environment just by pulling it from our registry.

We have a growing number of ready-to-use Dogus, e.g. SCM-Manager, Jenkins, Nexus Repository, SonarQube, Redmine and many more. Every Dogu can be tailored to your specific needs. Take advantage of a central authentication service, a dynamic navigation, that lets you easily switch between the web UIs and a smart configuration magic, which automatically detects and responds to dependencies between Dogus.

The Cloudogu EcoSystem is open source and it runs either on-premises or in the cloud. The Cloudogu EcoSystem is developed by Cloudogu GmbH under [AGPL-3.0-only](https://spdx.org/licenses/AGPL-3.0-only.html).

## License
Copyright © 2020 - present Cloudogu GmbH
This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
You should have received a copy of the GNU Affero General Public License along with this program. If not, see https://www.gnu.org/licenses/.
See [LICENSE](LICENSE) for details.


---
MADE WITH :heart:&nbsp;FOR DEV ADDICTS. [Legal notice / Imprint](https://cloudogu.com/en/imprint/?mtm_campaign=ecosystem&mtm_kwd=imprint&mtm_source=github&mtm_medium=link)
