# prom_wrapper_1

Push method...
metrics individually specified as separate var's

- Start Prometheus

docker run \
    -p 9090:9090 \
    -v /Users/george/Desktop/ProjectsCommon/prometheus/config:/etc/prometheus \
    prom/prometheus


- Start Grafana
docker run -p 3000:3000 grafana/grafana-enterprise