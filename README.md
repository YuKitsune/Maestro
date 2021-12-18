
# Todo
- [X] Tidy up the frontend
  - [X] Back to Home button
  - [X] Artist images
  - [X] Album images
  - [X] Track images (use album image if necessary)
  - [X] ~~Remix-ify~~ Wanted to do a nested route, but remix doesn't support nested index route
  - [x] ~~Theme~~ Too complicated, UI is simple enough
  - [X] Proper tidy up
- [X] Update docker compose file
  - [X] include frontend
  - [X] ~~specify config files for API~~ Just using environment variables
- [ ] Services
  - [X] Fix spotify token expiring every now and then
  - [X] Fix deezer share link
  - [X] Fix Apple Music track artwork
  - [ ] Deep search
    - [ ] When a track has been found
      - [ ] Query other streaming services for the same ISRC code
    - [ ] When an album has been found
      - [ ] Query other streaming services for the same album and artist name
      - [ ] Ensure all track ISRCs match
      - [ ] Create entries for each track too
    - [ ] When an artist has been found
      - [ ] Query other streaming services for the same artist based on name
      - Don't bother creating links for albums and tracks, that would take way too long...
  - [X] Dedicated services API
    - [X] List all available services with logos
  - [X] Use Apple Music private key, team ID and key ID in config instead of token
  - [X] Tidy up regex for better handling

- [X] Error handling
  - [X] Return proper error responses
  - [X] Handle errors on the frontend
  - [X] Back button to home page
  - [X] Report button for a quick GitHub issue

- [X] Logos
  - [X] Figure out how we're gonna do logos

- [X] Logging
- [ ] Metrics
- [ ] Analytics

- [ ] Refactoring
  - [ ] Clean up context usages
  - [ ] Clear up names
  - [ ] Make streaming service configs a bit nicer
  - [ ] Optimize API docker file
  - [ ] Clean up logging

- [ ] More services
  - [ ] Tidal
  - [ ] iHeartRadio
  - [ ] Pandora
  - [ ] Amazon Music
  - [ ] YouTube music
  - [ ] Bandcamp (Maybe)
  - [ ] SoundCloud (Maybe)
