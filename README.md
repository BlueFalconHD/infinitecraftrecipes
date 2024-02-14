# Infinite Craft Recipe Brute Forcer
> Made for [Infinite Craft](https://neal.fun/infinite-craft/) by [Neal Argarwal](https://twitter.com/nealagarwal)

This go program can try every possible recipe for a certain number of iterations. Running for 3 hours gave me ~14mb of data (the format was a horrible one, this is rewritten with a much better format)

## Usage
You must modify the program to change the iteration count (soon to be a command line option, along with rate limiting delay time, etc.). Running `go run main.go` or building and running will start the process. The time it takes for each iteration isn't linear, the more iterations, the more time.

I **highly** recommend using a VPN when running this, so if you do get ratelimited, you can still use Infinite Craft when your turn your VPN off.

## Output format
After running, a file (`crafting_data.json`) will be outputted, with the following data structure.

### `items`
Each item in the `items` object is keyed with the item ID (just the name of the item), and formatted like so:

```json
"Allergy": {
  "emoji": "🤧",
  "name": "Allergy",
  "recipes": [
    "Pollen+Surf",
    "Swamp+Pollen",
    "Wind+Pollen",
    "Storm+Pollen",
    "Pollen+Dandelion",
    "Pollen+Fog",
    "Pollen+Dust"
  ]
}
```

### `recipes`
Each item in the `recipes` object is keyed with the two components of the recipe formatted like: `Item1+Item2`, the same format used in the recipes array of items. It contains a `first`, `second`, and `result` property, all being item IDs.

```json
"Ash+Dandelion": {
  "first": "Ash",
  "second": "Dandelion",
  "result": "Wish"
}
```

## Ideas
I really want to implement or add some of the following features:

- Weight function determining which recipes should be attempted first
- ~~Web UI/Recipe browser~~ https://github.com/BlueFalconHD/infinite-craft-recipes-web
- Proxy support to bypass ratelimiting
- CLI + options

If you know how, feel free to submit a PR implementing any of these.
