# Week 5 Alert Response Runbook

This runbook explains how to investigate and respond to alerts created during Week 5.

## General Alert Response

1. Read the alert name, severity, description, namespace, pod, and container.
2. Confirm whether the alert is still firing in Prometheus or Grafana.
3. Check affected pods and recent Kubernetes events.
4. Inspect application and infrastructure logs.
5. Apply the safest corrective action.
6. Confirm the metric returns to normal and the alert resolves.
7. Record the cause and corrective action.

Useful commands:

    kubectl get pods -A
    kubectl get events -A --sort-by=.lastTimestamp
    kubectl logs -n <namespace> <pod-name> --tail=100
    kubectl describe pod -n <namespace> <pod-name>

## HighApplicationLatency

### Trigger

P95 application latency remains above 1000 milliseconds for more than two minutes.

### Investigation

1. Check the Grafana P95 Request Latency panel.
2. Identify recent application changes.
3. Check pod CPU, memory, readiness, and restart status.
4. Inspect backend and Istio proxy logs.
5. Check downstream services and network errors.

Useful commands:

    kubectl get pods -A
    kubectl top pods -A
    kubectl logs -n week1 <backend-pod> --tail=100
    kubectl describe pod -n week1 <backend-pod>

### Response

- Roll back a faulty deployment if latency began after a release.
- Increase replicas if traffic exceeds available capacity.
- Fix slow application code or downstream services.
- Confirm P95 latency returns below 1000 milliseconds.

## HighApplicationErrorRate

### Trigger

More than 10 percent of application requests return HTTP 5xx errors for over two minutes.

### Investigation

1. Check the Grafana Error Rate and Request Rate panels.
2. Inspect backend, frontend, Kong, and Istio logs.
3. Check service endpoints and pod readiness.
4. Identify recent configuration or deployment changes.

Useful commands:

    kubectl get pods -A
    kubectl get endpoints -n week1
    kubectl logs -n week1 <backend-pod> --tail=100
    kubectl logs -n kong <kong-pod> --tail=100

### Response

- Roll back a faulty deployment or configuration.
- Restart only the affected unhealthy pod when appropriate.
- Restore unavailable dependencies.
- Confirm the error rate falls below 10 percent.

## PodCrashLooping

### Trigger

A container restarts more than twice within five minutes.

### Investigation

1. Identify the affected namespace, pod, and container.
2. Inspect pod status and Kubernetes events.
3. Check current logs and previous crash logs.
4. Check configuration, secrets, resource limits, and health probes.

Useful commands:

    kubectl describe pod -n <namespace> <pod-name>
    kubectl logs -n <namespace> <pod-name> -c <container>
    kubectl logs -n <namespace> <pod-name> -c <container> --previous

### Response

- Correct invalid configuration or missing secrets.
- Fix the failing startup command or application error.
- Adjust resource limits only when evidence shows resource exhaustion.
- Roll back a faulty image.
- Confirm the pod remains Running without further restarts.

## Week5DeliberateTestAlertV3

### Purpose

This test alert deliberately uses vector(1) to verify the complete alert path:

Prometheus rule -> Alertmanager -> webhook notification receiver.

### Test Procedure

    kubectl apply -f monitoring/test-alert.yaml
    kubectl logs -n monitoring deployment/week5-alert-receiver --since=5m

### Expected Result

The receiver logs contain:

- Week5DeliberateTestAlertV3
- Status firing
- Receiver week5-webhook

### Cleanup

    kubectl delete -f monitoring/test-alert.yaml
