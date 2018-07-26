# Autoupdater

### Task List first release
- [x] Fetch updates passively
- [x] Fetch updates actively 
- [x] Actively look for updates on github
- [x] Actively look for updates on dockerhub 
- [x] Update docker services
- [x] Update standalone applications
- [x] Cli configuration
- [x] File configuration
- [x] Save status
- [x] Refactor configuration adding updaters as commands instead of as flags
- [ ] Write actual documentation
- [ ] Add .travis.yml for CI testing
- [ ] Add releases to github
- [x] Test vendoring is correct: Have trouble downloading with dep through VPN, the connection drops before completion 90% of the times.
- [ ] CircuitBreaker with incremental backoff when contacting external services

### Task List future releases
- [ ] Unit testing for docker services (complicated)
- [ ] Add integration testing
- [ ] Create API that allows to retrieve information about current versions, as well as to tweak parameters.
- [ ] Force update of a service through API
- [ ] Who updates the updater?