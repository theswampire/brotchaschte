package geocode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/hbollon/go-edlib"
)

type EPSG_SpatialReference uint16

const (
	LV95  EPSG_SpatialReference = 2056
	WGS84 EPSG_SpatialReference = 4326
)

const API_ENDPOINT = "https://api3.geo.admin.ch/rest/services/api/SearchServer"

type GeoCode struct {
	X     float64
	Y     float64
	Label string
	Id    string
}

type AddressToCoordsResponse struct {
	Results []struct {
		Id     int `json:"id"`
		Weight int `json:"weight"`
		Attrs  struct {
			Detail string  `json:"detail"`
			Label  string  `json:"label"`
			X      float64 `json:"x"`
			Y      float64 `json:"y"`
		} `json:"attrs"`
	} `json:"results"`
}

type Address struct {
	id      string
	address string
	zipcode string
	city    string
}

func (a *Address) combinedWithoutCity() string {
	return fmt.Sprintf("%s %s", a.address, a.zipcode)
}

type MaybeGeoCode struct {
	GeoCode
	ok            bool
	hasSuggestion bool
	error_message string
}

type GeocodeError struct {
	HasSuggested      bool
	SuggestedEasting  float64
	SuggestedNorthing float64
	SuggestedLabel    string
	Message           string
}

func (ge *GeocodeError) Error() string {
	return fmt.Sprintf("geocode error: %s", ge.Message)
}

func newGeocodeErrorWithSuggestion(easting, northing float64, label, msg string) *GeocodeError {
	return &GeocodeError{HasSuggested: true,
		SuggestedEasting:  easting,
		SuggestedNorthing: northing,
		SuggestedLabel:    label,
		Message:           msg,
	}
}
func newGeocodeError(msg string) *GeocodeError {
	return &GeocodeError{HasSuggested: false, Message: msg}
}

func AddressToCoords(address Address, id string, sr EPSG_SpatialReference) (GeoCode, error) {
	u, err := url.Parse(API_ENDPOINT)
	if err != nil {
		return GeoCode{}, newGeocodeError("invalid base api")
	}

	// expand "str" to "strasse" as no street in switzerland ends on str
	street_expanded := strings.ReplaceAll(address.address, "str ", "strasse ")
	street_expanded = strings.ReplaceAll(street_expanded, "str. ", "strasse ")
	street_expanded = strings.TrimSpace(strings.ToLower(street_expanded))
	street_expanded = fmt.Sprintf("%s %s", street_expanded, address.zipcode)

	values := u.Query()

	values.Add("origins", "address")
	values.Add("type", "locations")
	values.Add("searchText", street_expanded+" "+address.city)
	values.Add("sr", fmt.Sprintf("%d", sr))

	u.RawQuery = values.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return GeoCode{Id: id}, newGeocodeError("failed to lookup")
	}

	var response AddressToCoordsResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return GeoCode{Id: id}, newGeocodeError("failed to decode JSON")
	}

	if len(response.Results) == 0 {
		return GeoCode{Id: id}, newGeocodeError("Keine Adressen wie diese gefunden")
	}

	candidates := make([]string, len(response.Results))
	for i, res := range response.Results {
		// label has format: "streetname number <b>zipcode city</b>"
		// transform to "streetname number zipcode"
		label_parts := strings.Split(res.Attrs.Label, "<b>")
		zipcode := strings.Split(label_parts[1], " ")[0]
		candidates[i] = strings.TrimSpace(strings.ToLower(label_parts[0] + zipcode))
	}

	match, _ := edlib.FuzzySearch(street_expanded, candidates, edlib.JaroWinkler)

	bestmatch := response.Results[slices.Index(candidates, match)]
	bestmatch.Attrs.Label = removeBTag(bestmatch.Attrs.Label)

	// exact match
	if match == street_expanded {
		return GeoCode{Id: id, X: bestmatch.Attrs.X, Y: bestmatch.Attrs.Y, Label: bestmatch.Attrs.Label}, nil
	}
	fmt.Printf("'%s' '%s'\n", match, street_expanded)

	return GeoCode{Id: id}, newGeocodeErrorWithSuggestion(
		bestmatch.Attrs.Y,
		bestmatch.Attrs.X,
		bestmatch.Attrs.Label,
		"bester Match: "+bestmatch.Attrs.Label)

}

func removeBTag(label string) string {
	return strings.ReplaceAll(strings.ReplaceAll(label, "<b>", ""), "</b>", "")

}

func AddressesToCoords(addresses []Address, sr EPSG_SpatialReference) []MaybeGeoCode {
	var wg sync.WaitGroup
	var coords []MaybeGeoCode
	ch := make(chan MaybeGeoCode, len(addresses))

	for _, address := range addresses {
		wg.Add(1)
		go func() {
			defer wg.Done()
			geocode, err := AddressToCoords(address, address.id, sr)
			var geocodeErr *GeocodeError
			if errors.As(err, &geocodeErr) {

				if geocodeErr.HasSuggested {
					geocode.Label = geocodeErr.SuggestedLabel
					geocode.X = geocodeErr.SuggestedEasting
					geocode.Y = geocodeErr.SuggestedNorthing
				}

				ch <- MaybeGeoCode{
					GeoCode:       geocode,
					ok:            false,
					hasSuggestion: geocodeErr.HasSuggested,
					error_message: geocodeErr.Message,
				}
				return
			}

			ch <- MaybeGeoCode{
				GeoCode: geocode,
				ok:      true,
			}
		}()
	}

	wg.Wait()
	close(ch)

	for c := range ch {
		coords = append(coords, c)
	}
	return coords
}
