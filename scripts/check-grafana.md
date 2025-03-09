# Checking Telemetry Data in Grafana

1. Open http://localhost:3000 (default credentials: admin/admin)
2. Click on "Explore" in the left sidebar

## Check Metrics
1. Select "Prometheus" from the data source dropdown
2. Enter this query: `test_metric`
3. Click "Run query"
4. You should see a data point with value 42

## Check Traces
1. Select "Tempo" from the data source dropdown
2. In "Search" tab:
   - Service Name: `test-service`
   - Min Duration: `0ms`
   - Max Duration: `1s`
3. Click "Run query"
4. You should see a trace named "test-span"

## Check Logs
1. Select "Loki" from the data source dropdown
2. Enter this query: `{service="test-service"}`
3. Click "Run query"
4. You should see the message "This is a test log message"

Note: If data doesn't appear immediately, wait a few seconds and try again. The OTLP collector might take a moment to process and forward the data.
