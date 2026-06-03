# Releasing

This repo ships **two independent artifacts**, each with its own version line
and its own release workflow. They are deliberately decoupled, so releasing one
does not release the other.

| Artifact | Version line | Trigger | Workflow | Output |
|----------|-------------|---------|----------|--------|
| **App container image** | `vX.Y.Z` (e.g. `v0.3.3`) | push a `v*` git tag | `.github/workflows/release-image.yaml` | `ghcr.io/spirl/spiffe-demo-app:vX.Y.Z` |
| **Helm chart** | `X.Y.Z` (e.g. `0.4.3`) | push to `main` touching `charts/**` | `.github/workflows/release.yml` | chart-releaser publishes the packaged chart |

> The two version lines are **not** in lockstep. As of this writing the image is
> on the `0.3.x` line and the chart on the `0.4.x` line. That is expected.

## Why the two are decoupled

`chart-releaser` creates a git tag named `spiffe-demo-app-<chartversion>` when it
publishes a chart. The image workflow triggers on `v*` specifically so that a
chart release (`spiffe-demo-app-*`) does **not** also fire an image build. Keep
these tag namespaces distinct. Do not point the image workflow at
`spiffe-demo-app-*`.

---

## Releasing the app image

Use this when the application source (`main.go`, `api/`, `go.mod`) has changed.

1. Make the source change, open a PR, review, merge to `main`.
2. Create a GitHub Release with a tag of the form `vX.Y.Z` (e.g. `v0.3.3`),
   targeting `main`.
   - A **draft** release does *not* push the tag, so the build does **not** fire
     until you click **Publish**. Use a draft if you want to review notes first.
   - Publishing pushes the tag, which fires `release-image.yaml`.
3. The workflow runs `ko build` and pushes `ghcr.io/spirl/spiffe-demo-app:vX.Y.Z`
   (multi-arch per `.ko.yaml`: linux/arm64 + linux/amd64).
4. Confirm the image tag appears in GHCR under the org `spirl` package
   `spiffe-demo-app`.

The image tag equals the git tag verbatim (e.g. tag `v0.3.3` gives image
`v0.3.3`), so it matches the chart's `values.yaml` convention (`image.tag:
vX.Y.Z`).

### Notes
- The image workflow uses Go `1.25.x` (`go.mod` requires `go 1.25.9`).
- GHCR auth is automatic: `ko-build/setup-ko` logs in with the workflow
  `GITHUB_TOKEN` (the job grants `packages: write`).

---

## Releasing the Helm chart

Use this to publish a new chart version, typically to adopt a newly released
image, or to change chart templates/values.

1. Edit `charts/spiffe-demo-app/`:
   - `values.yaml`: set `image.tag` to the image you want deployed by default
     (e.g. `v0.3.3`).
   - `Chart.yaml`: bump `version` (e.g. `0.4.2` to `0.4.3`).
     **Required**, because chart-releaser dedupes by version and will silently
     republish nothing if the version is unchanged.
   - `README.md`: regenerate the values table if `values.yaml` changed
     (`helm-docs`), or edit the affected rows by hand if `helm-docs` isn't
     installed.
2. Open a PR, review, merge to `main`.
3. The push to `main` (path `charts/**`) fires `release.yml`, which runs
   `helm/chart-releaser-action` to package and publish chart `X.Y.Z`.

### Sequencing
If a chart release is meant to adopt a new image, **publish the image first**
(so `ghcr.io/spirl/spiffe-demo-app:vX.Y.Z` exists), then point the chart at it.
Otherwise the chart references a tag that isn't in the registry yet.

---

## End-to-end: shipping an app change

1. App source change, PR, merge to `main`.
2. Publish image release `vX.Y.Z` (fires image build, pushes to GHCR).
3. Verify the `vX.Y.Z` image is in GHCR.
4. Bump chart `values.yaml` (`image.tag`) and `Chart.yaml` (`version`), PR, merge.
5. Chart `X.Y.Z` is published by chart-releaser.
6. Re-point the deployment at the new chart version (lives in the
   `defakto-deploy` repo, outside this one).
