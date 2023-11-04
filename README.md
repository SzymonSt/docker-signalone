# docker-signalone

> [!WARNING]
> Extension is still work in progress. It is not published as official SignalOne Docker extension yet.

## Overview
Signal0ne docker extension is a tool for debugging and monitoring containerized apps, which enables automated insights about failed containers and containers affected by resource usage anomalies.

## How to use

To build extension run following command:
```bash
docker buildx build -t 322456/signalone-ext:dev  .
```

To run extension:
```bash
docker extension install 322456/signalone-ext:dev
```

Optionally you can build every underlying component separetly using:
```bash
cd vm
docker buildx build -t 322456/signalonemlworker:dev -f ./mlworker/Dockerfile .
docker buildx build -t 322456/signaloneapi:dev -f ./api/Dockerfile .
cd ../signal-front
docker buildx build -t 322456/signaloneui:dev  .
```

## Reporting issues

Please report issues using "Issues" github repository tab. Do not duplicate issues.

## Contributing
To contribute to this project start by browsing through open issues. If you find any issue you can help with do a fork and create a pull request.

## License
MIT