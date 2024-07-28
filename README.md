# et

<div style="text-align:center">
    <img src="images/logo.png" alt="logo" style="width:150px;"/>
</div>

The environment toolkit provides a set of tools for building Infrastructure as Code (IaC) configurations on public cloud providers. It consists of a CLI called `et` to bootstrap, build, publish, and execute `beacons`, which are collections of resources. These `beacons` are indexed and made discoverable through the `hub`.

This way, the environment toolkit brings composition to infrastructure, through a set of `beacons` composed together to build a `stack`.

## Architecture Terms

1. `et` - the cli
1. Hub - the default hub for beacons `hub.envt.io`
1. Beacon - CDK-TF modules created by the community and registered in the Environment Toolkit `hub`.
1. Spec - the spec file which contains a `class`, `code`, `data` and a list of resources
1. Stack - the output of running a spec with Environment Toolkit
1. Class - a term we took from CMDB, which helps validate a list of Beacons
1. Resource - a Beacon + Data = a resource
1. Primary Resource - Every class has a set of resources that are primary and atleast one of those primary services needs to be defined
1. Secondary Resource - Every class has a set of resrouces that are secondary, these are resources that support the primary service and are optional
1. Reference - A reference will connect a resource to an external stack

## Start with an example

Below is an `et.yml` file, we'll be using this example to explain how everything works.

```yaml
class: compute
code: outfra

data:
  aws:
    role: arn:aws:iam::211125614781:role/et-deployment
    state:
      bucket: et-example-prod
      role: arn:aws:iam::211125614781:role/et-bucket
  environment: prod
  region: us-east-1

service:
  type: container
  data:
    image: ${{ env:image }}
    port: 8080
    desired: 2
    cpu: 256
    memory: 512
    env_vars:
      SERVICE: bond
      VERSION: v1.0.0
      STREAM_TYPE: apub
      STREAM_AWS_QUEUENAME: ${{ this:events:name }}
      STREAM_AWS_TOPICARN: ${{ this:events:topic-arn }}
      DATA_TYPE: pg
      DATA_PG_HOST: ${{ this:db:host }}
      DATA_PG_PORT: ${{ this:db:port }}
      DATA_PG_DATABASE: ${{ this:db:name }}
      DATA_PG_SSLMODE: require
      API_IDENTITYURL: ${{ outbound:identity:url }}
      API_CONFIGURL: ${{ outbound:config:url }}
    secrets:
      DATA_PG_USERNAME: ${{ this:db:userpass:username }}
      DATA_PG_PASSWORD: ${{ this:db:userpass:password }}
      SECURITY_SIGNKEY: ${{ secret:signkey }}
    edge:
      edges:
        - edge: api
          environment: prod
      health:
        path: /live
        port: 8082
        timeout: 25
        interval: 30
        unhealthy_threshold: 5
        healthy_threshold: 5
      rules:
        - path: /stacks/import
          priority: 100
  overrides:
    environments:
      - environment: prod
        data:
          env_vars:
            UGLYDOMAIN_DOMAIN: n-cc.net

resources:
  - code: events
    type: queue
    data:
      is_fifo: true
      subscription: events
  - code: db
    type: database
    data:
      cluster: database

access:
  inbound:
    - api # not really needed since the edge references the api
  outbound:
    - config # not really needed since the container references the config
    - identity # not really needed since the container references the identity
```

## CLI

The cli to manage Stacks and Beacons

## Beacon

A CDKTF library of Constructs built on top of the Environment Toolkit library of Constructs. The example above uses the `queue`, `database` and `container` Beacons.

## Spec

See above for an example.

## Class

A way to group resources together to form a common function. The master list of Classes can be found here: `{{TBA}}`. The list is community driven, and classes are grouped by cloud provider.

In the above example the class is a `compute` type. A `compute` class has the primary resource of `function | container` and secondary resources of `database`, `queue` and `bucket`.
