# et

<div style="text-align:center">
    <img src="images/logo.png" alt="logo" style="width:150px;"/>
</div>

Environment Toolkit (`et`) is an Infrastructure as Code tool created by Platform Engineers at No_Ops.

We bring composition to infrastructure, the community has built a set of Blueprints which can be composed together to build a stack.

## Architecture Terms

1. `et` - the cli
2. Blueprints - CDK-TF modules created by the community and registered in the Environment Toolkit registry
3. Stack - a set of Blueprints composed together using Environment Tookit
4. Classes - a term we took from CMDB, which helps validate a list of Blueprints
5. Resources - a Blueprint + Data = a resource
6. Primary Resources - Every class has a set of resources that are primary and atleast one of those primary services needs to be defined
7. Secondary Resources - Every class has a set of resrouces that are secondary, these are resources that support the primary service and are optional

## Start with an example

Below is an et.yml file, we'll be using this example to explain how everything works.

```
class: compute
code: outfra

data:
  aws:
    role: arn:aws:iam::211125614781:role/et-deployment
    state:
        bucket: et-example-prod
        role: arn:aws:iam::211125614781:role/et-bucket
  environment: prod
  region: us-east1

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
  - code: outfra
    type: container
    data:
      image: {{env:image}}
      port: 8080
      desired: 2
      cpu: 256
      memory: 512
      env_vars:
        SERVICE: bond
        VERSION: v1.0.0
        STREAM_TYPE: apub
        STREAM_AWS_QUEUENAME: {{this:events:name}}
        STREAM_AWS_TOPICARN: {{this:events:topic-arn}}
        DATA_TYPE: pg
        DATA_PG_HOST: {{this:db:host}}
        DATA_PG_PORT: {{this:db:port}}
        DATA_PG_DATABASE: {{this:db:name}}
        DATA_PG_SSLMODE: require
        API_IDENTITYURL: {{outbound:identity:url}}
        API_CONFIGURL: {{outbound:config:url}}
      secrets:
        DATA_PG_USERNAME: {{this:db:userpass:username}}
        DATA_PG_PASSWORD: {{this:db:userpass:password}}
        SECURITY_SIGNKEY: {{secret:signkey}}
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
              AWS_ACCOUNTID: "211125614781"
              UGLYDOMAIN_DOMAIN: n-cc.net
              UGLYDOMAIN_HOSTEDZONE: Z059167936DYP5QD4B3PQ
              UGLYDOMAIN_ROLEEXTERNALID: 9a2777d2-757f-4b80-b63c-a698d0826cf9
              UGLYDOMAIN_ROLEARN: arn:aws:iam::211125614781:role/PublicDomainExternal

access:
  inbound:
    - api # not really needed since the edge references the api
  outbound:
    - config # not really needed since the container references the config
    - identity # not really needed since the container references the identity
```

## Class

A way to group resources together to form a common function. The master list of Classes can be found here: {{TBA}}. The list is community driven, and classes are grouped by cloud provider.

In the above example the class is a 'compute' type.
