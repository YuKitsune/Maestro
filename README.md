
# TODO
## Current issues
- Naming differences across services makes it hard to search
  - Example:
    - Hyperlandia (feat. Foster the People) [Vocal Mix]
    - Hyperlandia - Vocal Mix
  - Idea:
    1. Fuzzy searching
       1. Match on artist name
       2. Fuzzy search all albums for the closest match
       3. Match on ISRC code for tracks
          1. Fuzzy search all tracks for the closest match in lieu of an ISRC code

## Test cases
- Artists
  - Unique Artist
  - Artists with similar name
  - Artists with "and" or "&"
- Albums
  - Names should be normalized
  - Unique Album
  - Album with multiple artists
  - Albums with similar name
  - Albums with "and" or "&"
- Track
  - Unique Track
    - Track with multiple artists
    - Track with multiple albums
    - Tracks with similar names
    - Tracks with "and" or "&"
