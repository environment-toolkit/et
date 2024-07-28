# CLI Commands

## Manage Beacons

## Bootstrap a Beacon

Initializes a new beacon with the necessary resources.

```console
et beacon new <BeaconName>
```

## Package a Beacon

Builds and packages the beacon for publishing.

```console
et beacon build
```

## Publish a Beacon

Uploads the beacon to the `hub` for others to discover.

> [!NOTE]
> Requires hub credentials to push. (Through Credential Chain).

```console
et beacon push
```

## Manage Hub auth

## Authenticate to the hub

> [!IMPORTANT]
> Required for `et beacon push`.

```console
et hub login
```

Temporary credentials are stored under `~/.et/auth` by default.

## Manage Stack

## Init

Init an empty spec file for a new stack

```console
et init
```

## Add a Beacon to the Stack

Downloads the beacon from the `hub`

```console
et add
```

## Stands up a Stack

use the fs to synthesizes the beacon into Terraform IaC, and applies it using Credential Chain for cloud provider.

```console
et up
```

## Tear Down a Stack

```console
et down
```
