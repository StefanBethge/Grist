# Examples

Sample recipes and expected inputs/outputs live here. Each example should be
self-contained so it can be used as a regression fixture for `grist run`.

## minimal

End-to-end smoke test for the Phase 1 walking skeleton.

```sh
cd examples
grist run minimal.yaml
cat minimal.out.csv
```

Input `minimal.in.csv` has padded cells and a legacy `E-Mail` header; the
recipe trims `name`/`E-Mail` and renames the header to `email`.
