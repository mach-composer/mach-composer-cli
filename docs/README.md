# MACH composer docs

This repo holds the source files for the official [MACH composer documentation](https://...).

## Installation

`make install`

## Preview

Run `make preview`, then check http://localhost:8000/.

## Editing

Make the required changes in Markdown format. Do not forget to add new files to `mkdocs.yml`. 

The following Markdown extensions have been enabled:

- [admonition](https://squidfunk.github.io/mkdocs-material/extensions/admonition/)
- [codehilite](https://squidfunk.github.io/mkdocs-material/extensions/codehilite/)
- [footnotes](https://squidfunk.github.io/mkdocs-material/extensions/footnotes/)
- [permalinks](https://squidfunk.github.io/mkdocs-material/extensions/permalinks/)

## Deployment

Push the changes. Any changes to master will trigger a new deployment on https://...
