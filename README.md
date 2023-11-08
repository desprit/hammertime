# Local

```sh
make run
```

# Deployment

## Create

```sh
cd server
fly launch --name hammertime --region ams --build-target final --copy-config --no-deploy
fly volume create sqlite --region ams -s 2
fly deploy
```

## Update

```sh
fly deploy
```

# Destroy

```sh
fly apps destroy hammertime
```
