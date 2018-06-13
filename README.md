# BoshSpecs

[![Build Status](https://travis-ci.org/ahelal/boshspecs.svg?branch=master)](https://travis-ci.org/ahelal/boshspecs)

**WARNING** *this is a POC, miles away to be production ready.*

BoshSpecs is go CLI tool that can run various testing frameworks (InSpec, GOSS, BATS, ...) targeting one or many bosh directors with deployed instances. 
Designed to work with Concourse to verify compliance.

## Overview 

Bosh releases are usually well tested and contain smoke test, but operators who deploy bosh releases either reuse or create a manifest file with little or no testing within their infrastructure. how to verify compliance with internal specification. 

A bosh deployment might be successful, but that does not mean it is working as needed Some example use cases.

* Cloudfoundry should have X number of diago cells, logging configure or sshing to apps is  disabled.
* Application is accessible and firewall is configured correctly.
* Security policies states all instances disks should be encrypted or all communication should be HTTPS, ...
* When teams get bigger specification and intent are lost. Describing them as code will help maintain consistency and avoid regression. 

Since most installation and configuration are described as code it makes sense to test our code. The goal is to try to assert specification are being met within the bosh deployments and infrastructure level.


## Supported verifiers 

Currently BoshSpecs supports three times of test frame workers "test verifiers" 
* [Inspec](https://www.inspec.io/) testing framework for infrastructure with a human- and machine-readable language for specifying compliance, security and policy requirements. 
* [GOSS](https://github.com/aelsabbahy/goss)
* shell you friendly normal bash. You can elevate test verifiers to run bats, shunit2 or bunit.

More to come if needed.

## Prerequisite

* A bosh director installed
* Bosh CLI  >= 2 installed


## simple example

You will need to create `.boshspecs.yml`  

```yaml
bosh:
  - name: "boshGCP"
    environment: x.x.x.x
    client: admin
    client-secret: pass
    ca-cert: test/deployments/ca_ingore.txt
# at least one deployment is required
deployments:
  - name: cf
# at least one spec is required
specs:
    - name: simpleTest
      type: shell
```

Boshspecs will look in current working directory `test/simpleTest/` for all *.sh files and execute them against all instances groups in deployment cf

running 

```sh
boshspecs verify
```

## complex example

```yaml
# no bosh section use environment variables instead
deployments:
  - name: cf
    specs: 
      - name: cf_api
        type: inspec
        filters:
            instance_group: "api" # only run on instance_group api

      - name: cf_cell
        type: inspec
        filters:
            instance_group: "diego-cell" # only run on instance_group diego-cell

  - name: concourse
    specs:
      - name: webTests
        type: inspec
        filters:
            instance_group: "web" # only run on instance_group webb

      - name: workerTests
        type: inspec
        filters:
            instance_group: "worker" # only run on instance_group worker

specs:
    - name: diskEncryption
      type: inspec

    - name: firewall
      type: shell
      local_exec: true # run locally instead of instance
      path: different/Test/Location
```
This configuration will generate the following combination of tests

```
BOSH  Deployment  Spec            Instance Group  Instance  ID
      cf          diskencryption  *               0         /cf/diskencryption
      cf          firewall        *               0         /cf/firewall
      cf          cf_api          api             0         /cf/cf_api
      cf          cf_cell         diego-cell      0         /cf/cf_cell
      concourse   diskencryption  *               0         /concourse/diskencryption
      concourse   firewall        *               0         /concourse/firewall
      concourse   webtests        web             0         /concourse/webtests
      concourse   workertests     worker          0         /concourse/workertests
```

## Commands

For help run `boshspecs --help`
```
COMMANDS:
     ping, p    Ping a bosh director
     verify, v  verify a deployment
     list, l    list specs
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -v             Be more verbose.
   --no-color, --nc          don't use color.
   --debug, -d               Enable debug mode.
   --config value, -c value  Config file. (default: ".boshspecs.yml")
   --help, -h                show help
   --version                 print BoshSpecs version
```

## Development

```
# Getting deps 
make deps

# run tests 
make test

# Build
make build
```
