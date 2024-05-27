# ghloc (GitHub Lines Of Code)

ghloc is a project for counting the number of non-empty lines of code in a project.

## About This Fork

It was made for my [GitHub LOC](https://github.com/DervexDev/github-loc) Chrome extension and contains following changes:

- Supports private repositories
- Requires additional authorization
- Optimized for Vercel deployment

## API

The idea is simple: you make a request to the API in the format `/<username>/<repository>/<branch>` and you get the response with human-readable JSON.

You can see only some files using `match` URL parameter, e.g. with `/someuser/somerepo?match=js` only paths containing `js` will be considered. Examples of more powerful usage:

- `match=.js$` will show only paths ending with `.js`.
- `match=^src/` will show only paths starting with `src/` (i.e. placed in the `src` folder).
- `match=!test` will filter out paths containing `test`.
- `match=!test,!.sum` will filter out paths containing `test` or `.sum`.
- `match=.json$,!^package-lock.json$` will show only json files except for `package-lock.json` file.

There is also `filter` URL parameter, which has the opposite behavior to `match` parameter. `filter` has the same syntax but it declares which files must be filtered out.

To make the response more compact (removing spaces from the json) you can use `pretty=false`, e.g. `/someuser/somerepo?pretty=false`.
