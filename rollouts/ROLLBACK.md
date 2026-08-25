# Automated Canary Rollback Procedure

## Purpose

The backend uses Argo Rollouts and Prometheus analysis to prevent an unhealthy revision from replacing the stable application.

## Canary Stages

1. Shift 20% capacity to the canary.
2. Pause for 30 seconds.
3. Run three Prometheus measurements at 15-second intervals.
4. Shift 50% capacity to the canary.
5. Pause for 30 seconds.
6. Repeat Prometheus analysis.
7. Promote the revision to 100%.

## Metric and Failure Condition

The `backend-success-rate` AnalysisTemplate calculates the proportion of destination requests returning HTTP 5xx responses.

The healthy condition is `result[0] <= 0.05`, meaning the error proportion must remain at or below 5%.

The rollout uses `failureLimit: 1`. A failed metric assessment immediately aborts the canary update.

## Automatic Rollback Behaviour

When analysis fails:

1. The new revision becomes `Degraded`.
2. The failed canary ReplicaSet is scaled down.
3. The previous stable ReplicaSet remains active.
4. The Service continues routing to healthy stable pods.
5. Argo CD remains degraded until Git contains a corrected configuration.

## Controlled Rollback Test

A controlled test temporarily changed the success condition to `result[0] < 0` and changed `APP_VERSION` to `rollback-test-v3`.

Because Prometheus returned a non-negative value, the metric failed. Argo Rollouts automatically aborted revision 3, scaled down its canary ReplicaSet, and kept all five revision 2 stable pods available.

The safe condition `result[0] <= 0.05` and `APP_VERSION=canary-v2` were restored afterward.

Evidence is stored in `screenshots/week6/02-automated-rollback.png`.

## Monitoring Commands

Watch rollout progress:

    kubectl argo rollouts get rollout backend -n week1 --watch

Inspect analysis results:

    kubectl get analysisrun -n week1
    kubectl describe analysisrun <analysis-run-name> -n week1

Retry after correcting Git configuration:

    kubectl argo rollouts retry rollout backend -n week1

Perform an emergency manual undo:

    kubectl argo rollouts undo backend -n week1

Request an immediate Argo CD refresh:

    kubectl annotate application week4-microservices -n argocd argocd.argoproj.io/refresh=hard --overwrite

## Recovery Verification

Run:

    kubectl argo rollouts get rollout backend -n week1
    kubectl get pods -n week1

Expected state:

- Status: `Healthy`
- Step: `7/7`
- SetWeight: `100`
- Desired, ready, and available replicas: `5`
