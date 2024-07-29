# et

<div style="text-align:center">
    <img src="images/logo.png" alt="logo" style="width:150px;"/>
</div>

The environment toolkit provides a set of tools for building Infrastructure as Code (IaC) configurations on public cloud providers. It consists of a CLI called `et` to bootstrap, build, publish, and execute `beacons`, which are collections of Provider API Objects. These `beacons` are indexed and made discoverable through the `hub`.

This way, the environment toolkit brings composition to infrastructure, allowing a set of `beacons` to be composed together into a `formation` simplifying product feature development.

## Architecture Terms

1. `et` - the cli
1. Hub - the default hub for beacons `hub.envt.io`
1. Beacon - CDK-TF higher level constructs created by the community and discoverable through the `hub`. <!-- Beacons are namespaced by class and expose composition contracts. -->
1. Formation - A workspace for `et` to track and resolve dependencies between Beacons.
1. Spec - the spec file which contains a Beacon type (`type`), identifier (`name`), Primpary Properties (`service.props`) and may include secondary supporting beacons
1. Class - a term we took from CMDB, which helps validate the specs, their properties and dependencies.
1. State - the result of running a Spec with Environment Toolkit. All Beacon states are tracked in a Formation
1. Resource - each Beacon + Props results in Resources within a State
1. Reference - A reference will connect Resources across States

## CLI

Refer to the [CLI spec](./CLI.md).

## Beacon Library

A CDKTF library of Constructs built on top of the base Environment Toolkit library of Beacon Constructs. The example above uses the `queue`, `database` and `container` Beacons.

## Class

A way to group resources together to form a common function. The master list of Classes can be found here: `{{TBA}}`. The list is community driven, and classes are grouped by cloud provider.

In the above example the class is a `compute` type. A `compute` class has the primary resource of `function | container` and secondary resources of `database`, `queue` and `bucket`.
