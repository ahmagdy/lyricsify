## Lyricsify
A tool to load the liked songs in Spotify, scrape the lyrics and add the ability to search through all of them.

### Use case
You're a Spotify user, you have a lot of songs in your playlist, you remembered multiple words or a sentence in one of the songs in your playlist, you tried to google it and you weren't able to find the song.
This tool should solve your problem because it will load all of your lyrics and allow you to search in all of it in an easy and reliable way.

#### Tools used
- [Wire](https://github.com/google/wire)
- [goquery](https://github.com/PuerkitoBio/goquery)
- [elastic](https://github.com/olivere/elastic)
- [Viper](https://github.com/spf13/viper)

## How to use it
Consider checking [example](https://github.com/Ahmad-Magdy/lyricsify/blob/master/example/main.go) folder.

#### Config
Expected values to be set as Environment Variables or in a YAML file in the Documents folder.
```yaml
GENIUS_BASE_URL: 
GENIUS_TOKEN:
SPOTIFY_TOKEN:
LYRICS_INDEX_NAME: 
```

#### Sample
```go
ctx := context.Background()
// Initialize Instance of Lyricsify
lyricsify := lyricsify.InitializeLyricsify(ctx)
// Load all of your songs in Liked Songs section, the key is the song name and the value is the artist/s
songsMap, err := lyricsify.LoadSongs(ctx)
// ......
// To get song lyrics as a text
lyrics, err := lyricsify.FetchLyrics(ctx, song, artists)

// Save Lyrics in the datastore
err := lyricsify.SaveLyrics(ctx, song, lyrics)

```
#### Docker Compose
Elasticsearch is required to run the tool, docker compose is ready to be used.
```bash
docker-compose up
```


## License:
[The MIT License](https://github.com/Ahmad-Magdy/lyricsify/blob/master/LICENSE)
