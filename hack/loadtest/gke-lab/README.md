# GKE lab perftest

One-time load validation for the **kollect-lab** GKE cluster lives in the sibling
`kollect-lab` repo:

```
kollect-repos/kollect-lab/tests/perftest/
```

- **`run-once.sh`** — guarded single run (prep → deploy → soak → snapshot)
- Default **10k** Deployments; scale env vars for **100k** design proof
- **`values-perftest.yaml`** — `resourcesProfile: large`

Upstream multi-cluster 100k bundle (reference): `./100k/`

See `kollect-lab/tests/perftest/README.md` and
[docs/operator-manual/load-test-runbook.md](../../../docs/operator-manual/load-test-runbook.md).
