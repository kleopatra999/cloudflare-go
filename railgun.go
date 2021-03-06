package cloudflare

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Railgun represents a Railgun's properties.
type Railgun struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	Enabled        bool      `json:"enabled"`
	ZonesConnected int       `json:"zones_connected"`
	Build          string    `json:"build"`
	Version        string    `json:"version"`
	Revision       string    `json:"revision"`
	ActivationKey  string    `json:"activation_key"`
	ActivatedOn    string    `json:"activated_on"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
	UpgradeInfo    struct {
		LatestVersion string `json:"latest_version"`
		DownloadLink  string `json:"download_link"`
	} `json:"upgrade_info"`
}

// RailgunListOptions represents the parameters used to list railguns.
type RailgunListOptions struct {
	Direction string
}

// railgunResponse represents the response from the Create Railgun and the Railgun Details endpoints.
type railgunResponse struct {
	Response
	Result Railgun `json:"result"`
}

// railgunsResponse represents the response from the List Railguns endpoint.
type railgunsResponse struct {
	Response
	Result []Railgun `json:"result"`
}

// CreateRailgun creates a new Railgun.
// API reference:
// 	https://api.cloudflare.com/#railgun-create-railgun
// 	POST /railguns
func (api *API) CreateRailgun(name string, org Organization) (Railgun, error) {
	uri := "/railguns"
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}

	params := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	res, err := api.makeRequest("POST", uri, params)
	if err != nil {
		return Railgun{}, errors.Wrap(err, errMakeRequestError)
	}
	var r railgunResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return Railgun{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// ListRailguns lists Railguns connected to an account.
// API reference:
//  https://api.cloudflare.com/#railgun-list-railguns
//  GET /railguns
func (api *API) ListRailguns(options RailgunListOptions, org Organization) ([]Railgun, error) {
	v := url.Values{}
	if options.Direction != "" {
		v.Set("direction", options.Direction)
	}
	uri := "/railguns" + "?" + v.Encode()
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, errMakeRequestError)
	}
	var r railgunsResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// RailgunDetails returns the details for a Railgun.
// API reference:
// 	https://api.cloudflare.com/#railgun-railgun-details
// 	GET /railguns/:identifier
func (api *API) RailgunDetails(railgunID string, org Organization) (Railgun, error) {
	uri := "/railguns/" + railgunID
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}

	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return Railgun{}, errors.Wrap(err, errMakeRequestError)
	}
	var r railgunResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return Railgun{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// RailgunZones returns the zones that are currently using a Railgun.
// API reference:
// 	https://api.cloudflare.com/#railgun-get-zones-connected-to-a-railgun
// 	GET /railguns/:identifier/zones
func (api *API) RailgunZones(railgunID string, org Organization) ([]Zone, error) {
	uri := "/railguns/" + railgunID + "/zones"
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, errMakeRequestError)
	}
	var r ZonesResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// enableRailgun enables (true) or disables (false) a Railgun for all zones connected to it.
// API reference:
//  https://api.cloudflare.com/#railgun-enable-or-disable-a-railgun
//  PATCH /railguns/:identifier
func (api *API) enableRailgun(railgunID string, org Organization, enable bool) (Railgun, error) {
	uri := "/railguns/" + railgunID
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}
	params := struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: enable,
	}
	res, err := api.makeRequest("PATCH", uri, params)
	if err != nil {
		return Railgun{}, errors.Wrap(err, errMakeRequestError)
	}
	var r railgunResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return Railgun{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// EnableRailgun enables a Railgun for all zones connected to it.
// API reference:
//  https://api.cloudflare.com/#railgun-enable-or-disable-a-railgun
//  PATCH /railguns/:identifier
func (api *API) EnableRailgun(railgunID string, org Organization) (Railgun, error) {
	return api.enableRailgun(railgunID, org, true)
}

// DisableRailgun enables a Railgun for all zones connected to it.
// API reference:
//  https://api.cloudflare.com/#railgun-enable-or-disable-a-railgun
//  PATCH /railguns/:identifier
func (api *API) DisableRailgun(railgunID string, org Organization) (Railgun, error) {
	return api.enableRailgun(railgunID, org, false)
}

// DeleteRailgun disables and deletes a Railgun.
// API reference:
// 	https://api.cloudflare.com/#railgun-delete-railgun
// 	DELETE /railguns/:identifier
func (api *API) DeleteRailgun(railgunID string, org Organization) error {
	uri := "/railguns/" + railgunID
	if org.ID != "" {
		uri = "/organizations/" + org.ID + uri
	}
	if _, err := api.makeRequest("DELETE", uri, nil); err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}
	return nil
}

// ZoneRailgun represents the status of a Railgun on a zone.
type ZoneRailgun struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Connected bool   `json:"connected"`
}

// zoneRailgunResponse represents the response from the Zone Railgun Details endpoint.
type zoneRailgunResponse struct {
	Response
	Result ZoneRailgun `json:"result"`
}

// zoneRailgunsResponse represents the response from the Zone Railgun endpoint.
type zoneRailgunsResponse struct {
	Response
	Result []ZoneRailgun `json:"result"`
}

// RailgunDiagnosis represents the test results from testing railgun connections
// to a zone.
type RailgunDiagnosis struct {
	Method          string `json:"method"`
	HostName        string `json:"host_name"`
	HTTPStatus      int    `json:"http_status"`
	Railgun         string `json:"railgun"`
	URL             string `json:"url"`
	ResponseStatus  string `json:"response_status"`
	Protocol        string `json:"protocol"`
	ElapsedTime     string `json:"elapsed_time"`
	BodySize        string `json:"body_size"`
	BodyHash        string `json:"body_hash"`
	MissingHeaders  string `json:"missing_headers"`
	ConnectionClose bool   `json:"connection_close"`
	Cloudflare      string `json:"cloudflare"`
	CFRay           string `json:"cf-ray"`
	// NOTE: CloudFlare's online API documentation does not yet have definitions
	// for the following fields. See: https://api.cloudflare.com/#railgun-connections-for-a-zone-test-railgun-connection/
	CFWANError    string `json:"cf-wan-error"`
	CFCacheStatus string `json:"cf-cache-status"`
}

// railgunDiagnosisResponse represents the response from the Test Railgun Connection enpoint.
type railgunDiagnosisResponse struct {
	Response
	Result RailgunDiagnosis `json:"result"`
}

// ZoneRailguns returns the available Railguns for a zone.
// API reference:
// 	https://api.cloudflare.com/#railguns-for-a-zone-get-available-railguns
// 	GET /zones/:zone_identifier/railguns
func (api *API) ZoneRailguns(zoneID string) ([]ZoneRailgun, error) {
	uri := "/zones/" + zoneID + "/railguns"
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, errMakeRequestError)
	}
	var r zoneRailgunsResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// Railgun returns the configuration for a given Railgun.
// API reference:
// 	https://api.cloudflare.com/#railguns-for-a-zone-get-railgun-details
// 	GET /zones/:zone_identifier/railguns/:identifier
func (api *API) ZoneRailgunDetails(zoneID, railgunID string) (ZoneRailgun, error) {
	uri := "/zones/" + zoneID + "/railguns/" + railgunID
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return ZoneRailgun{}, errors.Wrap(err, errMakeRequestError)
	}
	var r zoneRailgunResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return ZoneRailgun{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// TestRailgunResponse tests a Railgun connection for a given zone.
// API reference:
//  https://api.cloudflare.com/#railgun-connections-for-a-zone-test-railgun-connection
//  GET /zones/:zone_identifier/railguns/:identifier/diagnose
func (api *API) TestRailgunConnection(zoneID, railgunID string) (RailgunDiagnosis, error) {
	uri := "/zones/" + zoneID + "/railguns/" + railgunID + "/diagnose"
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return RailgunDiagnosis{}, errors.Wrap(err, errMakeRequestError)
	}
	var r railgunDiagnosisResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return RailgunDiagnosis{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// connectZoneRailgun connects (true) or disconnects (false) a Railgun for a given zone.
// API reference:
//  https://api.cloudflare.com/#railguns-for-a-zone-connect-or-disconnect-a-railgun
//  PATCH /zones/:zone_identifier/railguns/:identifier
func (api *API) connectZoneRailgun(zoneID, railgunID string, connect bool) (ZoneRailgun, error) {
	uri := "/zones/" + zoneID + "/railguns/" + railgunID
	params := struct {
		Connected bool `json:"connected"`
	}{
		Connected: connect,
	}
	res, err := api.makeRequest("PATCH", uri, params)
	if err != nil {
		return ZoneRailgun{}, errors.Wrap(err, errMakeRequestError)
	}
	var r zoneRailgunResponse
	if err := json.Unmarshal(res, &r); err != nil {
		return ZoneRailgun{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// ZoneRailgun connects a Railgun for a given zone.
// API reference:
// 	https://api.cloudflare.com/#railguns-for-a-zone-connect-or-disconnect-a-railgun
// 	PATCH /zones/:zone_identifier/railguns/:identifier
func (api *API) ConnectZoneRailgun(zoneID, railgunID string) (ZoneRailgun, error) {
	return api.connectZoneRailgun(zoneID, railgunID, true)
}

// ZoneRailgun disconnects a Railgun for a given zone.
// API reference:
//  https://api.cloudflare.com/#railguns-for-a-zone-connect-or-disconnect-a-railgun
//  PATCH /zones/:zone_identifier/railguns/:identifier
func (api *API) DisconnectZoneRailgun(zoneID, railgunID string) (ZoneRailgun, error) {
	return api.connectZoneRailgun(zoneID, railgunID, false)
}
