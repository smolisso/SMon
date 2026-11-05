# Nagios compatibility baseline

Read-only inspection on 2026-11-05 of the user-provided Ansible tree found:

- 16 active `cfg_file` includes, 283 registered host blocks, 393 registered service
  blocks, 13 hostgroups, 72 command definitions and 46 directly referenced check
  command names. All those command names resolve to definitions.
- 2,189 direct host/service assignments, before deduplication and complete
  inheritance resolution; this is not a validated runtime service count.
- 29 locally present plugin files referenced by check definitions, in addition
  to standard plugin references. The plugin directory also contains backups,
  source and support artifacts and must not be imported wholesale.

These are source-tree observations, not a runtime Nagios validation. No existing
plugins were run, no credentials were copied, and neither source repository was
changed.

| Requirement | Evidence / scope | SMon status |
| --- | --- | --- |
| Exit status 0/1/2/3 | Nagios service plugin protocol | Planned |
| Short and multiline output | Nagios plugin protocol | Planned |
| Raw and parsed performance data | Output separated with `\|` | Planned |
| Command macros | `$USER1$`, `$HOSTADDRESS$`, `$ARG1$` through `$ARG5$` | Planned |
| Script interpreters | POSIX shell, Bash, Perl, Python; one Python 3.11 shebang | External dependencies |
| Perl dependencies | Includes `Net::SNMP`, `Date::Parse`, `XML::Simple`, `IO::Socket::SSL` | Inventory only |
| External programs | SNMP tools, SSH, curl, OpenSSL, smbclient, jq, bc | Inventory only |
| Shell command syntax | Notifications/perfdata use pipes or redirection | Explicit future adapter |
| Configuration inheritance | `use`, templates, direct host lists, groups | Planned |
| Operational behavior | Check/retry intervals, periods, attempts, icons | Planned |
| Full Nagios replacement | Import, host semantics, notifications, dependencies, passive checks | Not claimed |

Initial runner verification must use local fixtures, with no network calls:

1. Each valid exit code and an invalid code.
2. Single-line and multiline text, raw performance data, empty output.
3. Timeout, child-process cleanup and output truncation.
4. Exact argument boundaries, including spaces and metacharacters.
5. Missing executable/interpreter and permission errors.

Only after those tests should selected real scripts be compared against Nagios
in an explicitly configured test environment.

Protocol references:

- https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/4/en/pluginapi.html
- https://www.monitoring-plugins.org/doc/guidelines.html
