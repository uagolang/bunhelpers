# Changesets

This directory contains changeset files that describe changes and their impact on versioning.

## How to add a changeset

Run the following command and follow the prompts:

```bash
make changeset
```

This will create a new changeset file describing your changes and the type of version bump needed (major, minor, or patch).

## Changeset file format

Changeset files are markdown files with the following format:

```markdown
---
bump: minor
---

Description of your changes goes here.
```

The `bump` field can be:
- `major` - Breaking changes (1.0.0 -> 2.0.0)
- `minor` - New features (1.0.0 -> 1.1.0)
- `patch` - Bug fixes (1.0.0 -> 1.0.1)

