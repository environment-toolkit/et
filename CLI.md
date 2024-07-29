# CLI Spec

Proposal of `et` cli features and commands.

## Manage Formations

The Environment Toolkit uses a workspace to keep track of the states it manages. Each `et` invocation is within the context of such workspace and referred to as the active `formation`.

> [!NOTE]
> Formations tie all the states managed by `et` together and are a crucial indexing mechanism.

The metadata tracked about states in a `formation` includes things such as:

- A list of Beacon libraries with their version constraint + their `locks`.
- A list of references to States provisioned. `et` tracks "revision" information related to each state. <!-- allows propagation of Spec -->

Formations are used by the toolkit for cross state lookups, dependency cycle prevention and execution orchestration.

### Activate a Formation

To change the active formation.

```console
et formation use ./examples/formations/myorg.yml
```

> Alternatively: use the global CLI `-f/--formation` flag for execution context.

The formation spec provides an initial mapping of environments to Cloud Provider accounts.

<!--
### Bootstrap Formation

Provision the necessary Cloud resources for the Environment Toolkit to manage environments in your public cloud provider.

```console
et formation bootstrap
```

### Orchestrate Formation

Orchestrate/Refresh a deployment across all states in a formation?

``console
et formation rollout
```

### Visualize environments

Export a dotgraph viz of one or more environments and all its states within a formation.

```console
et formation graph [env1,env2]
```
-->

### Register a Beacon library

Make all beacons of a target Beacon library available for `et init` and track the version constraint in the `formation`.

```console
et add @envtio/base[@version-constraint]
```

>[!NOTE]
> The version [constraint](https://docs.npmjs.com/about-semantic-versioning#using-semantic-versioning-to-specify-update-types-your-package-can-accept) controls library updates.

<!-- TODO: Future feature of managing private Beacon pkges and credentials to fetch them -->

## Manage Beacon Libraries

Beacons are published as libraries. Currently, each beacon library is bootstrapped as a separate `git` repository.

<!-- projen by default initializes an empty git directory, in the future we may want to support monorepos better -->

### Bootstrap a Beacon Library

Initializes a new beacon Library in a directory of similar name.

```console
et lib new <library-name> [directory-name]
```

<!-- Implementation details - First Iteration

As a first iteration, this largely depends on the [JSII projen manifest](https://github.com/projen/projen/blob/v0.84.10/src/cli/util.ts#L62) and runs `projen new` in the target directory under the hood.

-->

### Package a Beacon Library

Builds and packages the beacon for publishing to the `hub`.

```console
et lib build
```

<!-- Implementation details - First Iteration

Projen is the task runner and re-uses [JSII pacmak](https://github.com/aws/jsii/blob/main/packages/jsii-pacmak/README.md) tasks to build and publish JSII packages for beacon libraries.

See the Projen JSII Project - [Packaging Tasks](https://github.com/projen/projen/blob/v0.84.10/src/cdk/jsii-project.ts#L507)

-->

### Publish a Beacon Library

Uploads the beacon to `hub.envt.io` and makes the Beacons within easily discoverable across the hub.

> [!NOTE]
> Requires hub credentials to push. (Through Credential Chain).

```console
et lib push
```

<!-- Implementation details - First Iteration

Requires Hub SaaS to handle authentication, backing npm registry and auto generated documentation functionality.

Future SaaS feature: private hubs
-->

## Manage Hub auth

### Authenticate to the hub

> [!IMPORTANT]
> Currently only required for `et lib push`.

```console
et hub login
```

Temporary credentials are stored under `~/.et/auth` by default.

## Manage State

### Init

Init a spec file for a beacon in the available `formation` libraries.

```console
et init <beacon-type>
```

<!-- init should prompt through the beacon props -->

### Stand up a Beacon

```console
et up [-s spec.yml] [environment] [region]
```

Evaluate and resolve the Spec properties across target environment(s), region(s) and its dependencies within the Formation.

Target `environment`/`region` is optional, when omitted all environments and regions are sequentially targetted.

<!-- CLI will sequentially synth, plan and apply environments/regions, SaaS offers event orchestrated state management? -->

<!-- technical implementation detail: 

Uses the fs to synthesizes the beacon into Terraform IaC with resolved references from the "Formation" (workspace), and apply it using Credential Chain for terraform provider.

-->

### Tear Down a Beacon

>[!IMPORTANT]
> This is `--dry-run` by default.

```console
et down [-s spec.yml] [environment] [region] [--no-dry-run]
```
