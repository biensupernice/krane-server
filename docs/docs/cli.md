# CLI

The Krane [CLI](https://github.com/krane/cli) allows you to interact with Krane to create container resources.

## Installing

```
npm i -g @krane/cli
```

## Authenticating

Krane uses [private and public key authentication](https://en.wikipedia.org/wiki/Public-key_cryptography). Both keys are used for ensuring authenticity of incoming requests.

1. Create the public and private key

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM' -f $HOME/.ssh/krane

-t type
-b bytes
-C comments
-m key format
-f output file
```

This will generate 2 different keys, a `private` & `public (.pub)` key.

2. Place the `public key` on the server where Krane is running, appended to `~/.ssh/authorized_keys`.

The `private key` is kept on the user's machine.

Now try authenticating. The CLI will prompt you to select the `private key` you just created. This will be used for authenticating with the `public key` located on the Krane server.

```
krane login
```

## Commands

### login

Authenticate with a Krane instance.

```
krane login
```

### remove

Remove a deployment.

```
krane remove <deployment>
```

### status

Returns information related to a deployment.

```
krane status <deployment>
```

or get the status of a single container.

```
krane status <deployment> <container>
```

### deploy

Create or run a deployment.

```
krane deploy -f </path/to/deployment.json>
```

Flags:

- `--file`(`-f`) Path to deployment configuration

  - required: `false`
  - default: `./deployment.json`

- `--tag`(`-t`) Image tag to apply to the deployment

  - required: `false`
  - default: `latest`

- `--scale`(`-s`) Number of containers to create
  - required: `false`
  - default: `1`

### list

List all deployments.

```
krane list
```

### logs

Read realtime container logs for a deployment.

```
krane logs <container>
```

### history

Get recent activity for a deployment.

```
krane history <deployment>
```

or get additional details on a specific deployment job.

```
krane history <deployment> <job>
```

### help

Display CLI usage and commands.

```
krane help
```

### start

Start all containers for a deployment. This command will not _create_ containers, only start any stopped containers.

```
krane start <deployment>
```

### stop

Stop all containers for a deployment.

```
krane stop <deployment>
```

### restart

Restart a deployment _re-creating_ container resources.

```
krane restart <deployment>
```

### secrets

List all deployment secrets.

```
krane secrets list <deployment>
```

Add a deployment secret.

```
krane secrets add <deployment> -k <key> -v <value>
```

Delete a deployment secret.

```
krane secrets delete <deployment> -k <key>
```