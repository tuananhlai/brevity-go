#!/bin/bash

# Get current timestamp in nanoseconds
CURRENT_TIME=$(date +%s)000000000

echo "Testing OTLP endpoints..."

echo -n "Sending metrics... "
response=$(curl -s -w "%{http_code}" -X POST -H "Content-Type: application/json" http://localhost:4318/v1/metrics -d '{
  "resource_metrics": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": { "stringValue": "test-service" }
      }]
    },
    "scope_metrics": [{
      "metrics": [{
        "name": "test.metric",
        "gauge": {
          "dataPoints": [{
            "asDouble": 42,
            "timeUnixNano": "'$CURRENT_TIME'"
          }]
        }
      }]
    }]
  }]
}')
if [ "${response: -3}" = "200" ]; then
    echo "✅"
else
    echo "❌ (HTTP ${response: -3})"
fi

echo -n "Sending traces... "
response=$(curl -s -w "%{http_code}" -X POST -H "Content-Type: application/json" http://localhost:4318/v1/traces -d '{
  "resource_spans": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": { "stringValue": "test-service" }
      }]
    },
    "scope_spans": [{
      "spans": [{
        "trace_id": "5b8aa5a4d5acdeee6cc43eab4410e89d",
        "span_id": "8a4dd7c97f804bed",
        "name": "test-span",
        "kind": 1,
        "start_time_unix_nano": "'$CURRENT_TIME'",
        "end_time_unix_nano": "'$CURRENT_TIME'"
      }]
    }]
  }]
}')
if [ "${response: -3}" = "200" ]; then
    echo "✅"
else
    echo "❌ (HTTP ${response: -3})"
fi

echo -n "Sending logs... "
response=$(curl -s -w "%{http_code}" -X POST -H "Content-Type: application/json" http://localhost:4318/v1/logs -d '{
  "resource_logs": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": { "stringValue": "test-service" }
      }]
    },
    "scope_logs": [{
      "log_records": [{
        "time_unix_nano": "'$CURRENT_TIME'",
        "severity_text": "INFO",
        "body": {
          "stringValue": "This is a test log message"
        }
      }]
    }]
  }]
}')
if [ "${response: -3}" = "200" ]; then
    echo "✅"
else
    echo "❌ (HTTP ${response: -3})"
fi

echo -e "\nAll requests sent. Check Grafana at http://localhost:3000"
