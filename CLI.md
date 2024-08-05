# CLI Spec

Proposal of `et` cli features and commands.

## Grid

The Environment Toolkit uses the `Grid` as a workspace to keep track of the states for the `Specs` it manages. Each `et` invocation must be within an "active" `grid`. Committing the `grid` to git improves reproducibility and collaboration.

A reference to the active grid is kept in `et`'s default config (`~/.et/config`).

> [!IMPORTANT]
> The `grid` ties all the states managed by `et` together and provides a crucial mechanism for lookups.

The metadata tracked about `States` includes things such as:

- The relationship between `Spec` files and their `States` across environments + regions.
- A list of Beacon libraries with the desired version constraint as well as the actual `lock`'ed version used.

The `grid` is required by the toolkit to allow cross `State` lookups, dependency cycle prevention, state refactoring and execution orchestration.

### Set Current grid

Minimal configuration is required to activate a grid such as an initial mapping of environments to Cloud Provider accounts. This information is defined in a `grid` manifest. see [examples/grids/my-org.yml](./examples/grids/my-org.yml)

```console
et use ./examples/grids/my-org.yml
```

> Alternatively: use the global CLI `-g / --grid` flag for executions across multiple grids.

<!--
### Bootstrap grid

Provision the conventional Cloud resources for the Environment Toolkit to manage environments within the Cloud Provider.

```console
et bootstrap ./examples/grids/my-org.yml
```
-->

### Register a Beacon library

Add a Beacon library version constraint into the current `grid` for reproducibility.

```console
et add @envtio/base[@version-constraint]
```

> [!NOTE]
> The version [constraint](https://docs.npmjs.com/about-semantic-versioning#using-semantic-versioning-to-specify-update-types-your-package-can-accept) controls library updates.

<!-- TODO: Future feature of managing private Beacon pkges auth mechanisms and facility the init command for available beacons -->

## Manage State

### Init Spec

Init a spec file for a beacon.

```console
et init [<library-ref>/]<beacon-type>
```

Init will automatically add the beacon library to the current `grid`.

### Stand up Spec

```console
et up [-f spec.yml] <environment> <region>
```

Evaluate and resolve the Spec properties and dependencies across target environment and region within the current `grid`.

<!-- CLI will:

- evaluate the spec, resolving resource references through the grid
- unresolved referenced properties halt the process
- resolved referenced properties are templated out
- stack synthesis and plan is executed using the Terraform Provider credential chain (i.e assume role arn)
- on confirmation stack is applied
- resulting State is indexed by UUID for env/region (attributes) in the grid and used for future spec resolving.

SaaS offers advanced orchestration mechanisms over a `grid` of states.

-->

### Tear Down Spec

> [!IMPORTANT]
> This is `--dry-run` by default.

```console
et down [-f spec.yml] <environment> <region> [--no-dry-run]
```

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
