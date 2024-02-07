# docker-signalone

> [!WARNING]
> Extension is still work in progress. It is not published as official SignalOne Docker extension yet.

## Overview
Signal0ne docker extension is a tool for debugging and monitoring containerized apps, which enables automated insights about failed containers and containers affected by resource usage anomalies.
![Alt text](image.png)


## How to install release

```
docker extension install 322456/signalone-extension:latest
```

## How to use locally

### Prerequisites
- Docker with compose
- Docker Desktop client
- Make

### Env variables
```
cp backend/.env.template backend/.default.env
# In backend/.default.env replace _APIKEY_ with your Huggingface API key 
# Adjust other variables if needed(optional)

cp ext/agent/.env.template ext/agent/.default.env
```

### Backend

```
make --directory=./backend build-backend

```
```
make --directory=./backend start-backend
```
OR
```
make --directory=./backend start-backend-with-init # to start backend with init sample development data
```

### Extension
```
#Build extension(both agent and frontend)
make --directory=./ext build-extension

#Run extension on your local docker desktop environment
make --directory=./ext install-extension-local
```

### Simulated development environment

```
make --directory=./ext start-devenv
```

## Reporting issues

Please report issues using "Issues" github repository tab. Do not duplicate issues.

## Contributing
To contribute to this project start by browsing through open issues. If you find any issue you can help with do a fork and create a pull request.

## License
MIT