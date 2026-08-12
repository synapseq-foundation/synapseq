# SynapSeq Remote

SynapSeq Remote provides ready-to-use sequences through a local catalog index.

## Custom Catalogs

Advanced users can sync a self-hosted catalog by passing its base URL:

```bash
synapseq -sync-url https://my-sequences.com
```

Custom catalogs serve their index at `/index.json`, unlike the official
catalog's `/free/index.json`.

Set `SYNAPSEQ_REMOTE_BASE_URL` to use the custom catalog with `-list`,
`-search`, `-info`, `-download`, and `-get` in later commands. Custom URLs must
be root HTTP(S) URLs without a path, credentials, query, or fragment.

SynapSeq stores each custom catalog in its own cache directory and validates
the catalog JSON and every downloaded SPSQ file before use.
