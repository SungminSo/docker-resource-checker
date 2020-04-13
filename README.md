# Docker Container Resource Monitoring
## Usage
to make binary: ``` make build ```

move to ``` ./build ```

``` ./checker -t|--time <time> -c|--container <docker container name> ```

## Help
```
Usage:
    ./checker -t|--time <time> -c|--container <docker container name>

Flags:
    -h, --help          help for docker resource monitor
    -t, --time          set monitoring time
    -c, --container     name of docker container for monitoring

Returns:
    - Print resource usage summary to console
    - Save the JSON file named resource_<timestamp>.json
```

## resource
#### CPU
- avg
- min
- max
#### Memory
- avg
- min
- max
#### Network
- in
- out
#### Disk
- write
- read