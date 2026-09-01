# Engines

An *engine* is the program responsible for executing run commands like `plan` and `apply`. OTF provides support for two engines:

* `terraform`
* `tofu`

The default engine is `terraform`. This can be overridden with the `otfd` flag [`--default-engine`](config/flags.md#-default-engine).

!!! warning
    If you're running more than one instance of `otfd`, take care to set this flag to the same value on each instance. Doing otherwise will lead to unpredictable results.

When you create a workspace, it'll use the default engine. You can override the engine for a workspace in its settings.

## Selecting an engine over the API

The workspaces API has no `engine` attribute — the TFE API it implements predates engines. To select one anyway, append a hint to the `terraform-version` attribute:

```
"terraform-version": "1.6.0 (tofu)"
```

The hint sets the workspace's engine alongside its version, in a single request. It works wherever `terraform-version` is accepted: creating and updating a workspace, and creating a run (where it overrides the workspace's engine for that run alone). Omit the hint and the engine is left as-is, so existing clients are unaffected.

The hint is only an input. It isn't stored or echoed back: reading a workspace returns a bare `terraform-version` plus a separate read-only `engine` attribute.

This is the supported path for migrating a workspace from terraform to tofu without changing `--default-engine`, which would affect every newly-created workspace.

OTF rejects a version the engine never published, so a workspace can't be left pinned to, say, tofu 1.5.7 (tofu's earliest release is 1.6.0). The check applies whether the engine changes via the API, the UI, or the hint.

When you create a run OTF will download the workspace's engine if it hasn't already been downloaded. The engine binaries are downloaded to the directory specified by the flag [`--engine-bins-dir`](config/flags.md#-engine-bins-dir).
