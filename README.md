# Iter

Proof of concept RAG server/website for travel recommendatinos

*Iter* comes from the Latin for "travel". The goal is to produce a LLM backed travel itinerary creator, frontend and backend, then see what we can link it to.

### Progress and TODOs:

#### TODOs
* Backend
  * User Login
  * Contextual Chat via embeddings
  * Integration with booking service (if I stick with this)
  * Geolocation - Mapbox Geocoding API gives me 100k tokens per month
  * Map generation
  * Saving the itinerary to a doc of some sort / Gannt chart.
 
* Frontend
  * Home Page
  * Log In Page
  * Itinerary Creator (given destination, wants, time)
  * Destination suggestion - limit to those we have context on.
  * Images
  * Itinerary rating / feedback
