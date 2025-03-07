# Tutorial
## Setup
Let's create our first load test
```bash
make start
```
Insert `GRAFANA_TOKEN` created in previous command
```bash
export LOKI_URL=http://localhost:3030/loki/api/v1/push
export GRAFANA_URL=http://localhost:3000
export GRAFANA_TOKEN=...
export DATA_SOURCE_NAME=Loki
export DASHBOARD_FOLDER=LoadTests
make dashboard
```
## Overview
General idea is to be able to compose load tests programmatically by combining different `Generators`

`Generator` is an entity that can execute some workload using some `Gun` or `Instance` definition, each `Generator` may have only one `Gun` or `Instance` implmenetation used

`Gun` can be an implementation of single or multiple sequential requests workload for stateless protocols

`Instance` is a stateful implementation that's more suitable for stateful protocols or when you client have some logic

Each `Generator` have a `Schedule` that control workload params throughout the test (increase/decrease RPS or Instances)

`Generators` can be combined to run multiple workload units in parallel or sequentially

`Profiles` are wrappers that allow you to run multiple generators with different `Schedules` and wait for all of them to finish

`AlertChecker` can be used in tests to check if any specific alerts with label and dashboardUUID was triggered and update test status

Load testing workflow can look like:
```mermaid
sequenceDiagram
    participant Product repo
    participant Runner
    participant K8s
    participant Loki
    participant Grafana
    participant Devs
    Product repo->>Product repo: Define NFR for different workloads<br/>Define application dashboard<br/>Define dashboard alerts<br/>Define load tests
    Product repo->>Grafana: Upload app dashboard<br/>Alerts has "requirement_name" label<br/>Each "requirement_name" groups is based on some NFR
    loop CI runs
    Product repo->>Runner: CI Runs small load test
    Runner->>Runner: Execute load test logic<br/>Run multiple generators
    Runner->>Loki: Stream load test data
    Runner->>Grafana: Checking "requirement_name": "baseline" alerts
    Grafana->>Devs: Notify devs (Dashboard URL/Alert groups)
    Product repo->>Runner: CI Runs huge load test
    Runner->>K8s: Split workload into multiple jobs<br/>Monitor jobs statuses
    K8s->>Loki: Stream load test data
    Runner->>Grafana: Checking "requirement_name": "stress" alerts
    Grafana->>Devs: Notify devs (Dashboard URL/Alert groups)
    end
```
## Examples

## RPS test
- [test](https://github.com/smartcontractkit/wasp/blob/master/examples/simple_rps/main.go#L9)
- [gun](https://github.com/smartcontractkit/wasp/blob/master/examples/simple_rps/gun.go#L23)
```
cd examples/simple_rps
go run .
```
Open [dashboard](http://localhost:3000/d/wasp/wasp-load-generator?orgId=1&var-test_group=generator_healthcheck&var-app=generator_healthcheck&var-cluster=generator_healthcheck&var-namespace=generator_healthcheck&var-branch=generator_healthcheck&var-commit=generator_healthcheck&from=now-5m&to=now&var-test_id=generator_healthcheck&var-gen_name=All&var-go_test_name=simple_rps&refresh=5s)

`Gun` must implement this [interface](https://github.com/smartcontractkit/wasp/blob/0dca04d432705472a8705ce473e175a77a3da9ed/wasp.go#L36)

## Instance test
- [test](https://github.com/smartcontractkit/wasp/blob/master/examples/simple_instances/main.go#L10)
- [instance](https://github.com/smartcontractkit/wasp/blob/master/examples/simple_instances/instance.go#L34)
```
cd examples/simple_instances
go run .
```
Open [dashboard](http://localhost:3000/d/wasp/wasp-load-generator?orgId=1&var-test_group=generator_healthcheck&var-app=generator_healthcheck&var-cluster=generator_healthcheck&var-namespace=generator_healthcheck&var-branch=generator_healthcheck&var-commit=generator_healthcheck&from=now-5m&to=now&var-test_id=generator_healthcheck&var-gen_name=All&var-go_test_name=simple_instances&refresh=5s)

`Instance` must implement this [interface](https://github.com/smartcontractkit/wasp/blob/2be83837defe2b1c7aa3aa187a34e698ff7fde69/wasp.go#L41)

## Profile test (group multiple generators in parallel)
- [test](https://github.com/smartcontractkit/wasp/blob/master/examples/profiles/main.go#L10)
- [gun](https://github.com/smartcontractkit/wasp/blob/master/examples/profiles/gun.go#L23)
- [instance](https://github.com/smartcontractkit/wasp/blob/master/examples/profiles/instance.go#L34)
```
cd examples/profiles
go run .
```
Open [dashboard](http://localhost:3000/d/wasp/wasp-load-generator?orgId=1&var-test_group=generator_healthcheck&var-app=generator_healthcheck&var-cluster=generator_healthcheck&var-namespace=generator_healthcheck&var-branch=generator_healthcheck&var-commit=generator_healthcheck&from=now-5m&to=now&var-test_id=generator_healthcheck&var-gen_name=All&var-go_test_name=my_test_ws&var-go_test_name=my_test&refresh=5s)

## Usage in tests
- [test](https://github.com/smartcontractkit/wasp/blob/master/examples/go_test/main_test.go#L15)
- [gun](https://github.com/smartcontractkit/wasp/blob/master/examples/go_test/gun.go#L23)
```
cd examples/go_test
go test -v -count 1 .
```
Open [dashboard](http://localhost:3000/d/wasp/wasp-load-generator?orgId=1&var-test_group=generator_healthcheck&var-app=generator_healthcheck&var-cluster=generator_healthcheck&var-namespace=generator_healthcheck&var-branch=generator_healthcheck&var-commit=generator_healthcheck&from=now-5m&to=now&var-test_id=generator_healthcheck&var-gen_name=All&var-go_test_name=TestProfile&refresh=5s)

## Checking alerts
- [test](https://github.com/smartcontractkit/wasp/blob/alerts/examples/alerts/main_test.go#L11)
- [gun](https://github.com/smartcontractkit/wasp/blob/alerts/examples/alerts/gun.go#L23)
```
cd examples/alerts
go test -v -count 1 .
```
Open [alert groups](http://localhost:3000/alerting/groups)