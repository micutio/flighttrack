# Flight Tracker

Testing flight tracking APIs and http/get requests in Go.

## Links for further investigation

### Flight Tracking APIs

Free ADSB APIs:

 - https://opendata.adsb.fi/api/v2/ (https://github.com/adsbfi/opendata)
 - https://github.com/ADSB-One/api
 - https://airplanes.live/api-guide/
   - might shift to feeder only
 - https://api.adsb.lol/docs
 - https://openskynetwork.github.io/opensky-api/rest.html
   - limited to [now] time, but access to bounding boxes and airport ARR/DEP
 - example for antarctica bounding box: https://opensky-network.org/api/states/all?lamin=-90&lomin=-180&lamax=-50&lomax=180

### Airline code to name mappings

- https://github.com/tbouron/MMM-FlightTracker

### Various ICAO code mappings

- https://github.com/rikgale/ICAOList

### Raw ADS-B data parsing in Go

- https://github.com/cjkreklow/go-adsb
