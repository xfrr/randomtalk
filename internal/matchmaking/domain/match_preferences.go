package matchdomain

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/location"
)

const (
	// default values for the ms preferences
	defaultMaxWaitSeconds = 10
	minAllowedAge         = 18
	maxAllowedAge         = 99
)

func DefaultPreferences() *MatchPreferences {
	return &MatchPreferences{
		state: &matchSessionPreferencesState{
			MaxAgeCriteria:     &maxAgeCriteria{maxAllowedAge},
			MinAgeCriteria:     &minAgeCriteria{minAllowedAge},
			GenderCriteria:     &genderCriteria{},
			InterestsCriteria:  &interestsCriteria{},
			DistanceCriteria:   &maxDistanceCriteria{},
			MaxWaitTimeSeconds: defaultMaxWaitSeconds,
		},
	}
}

type matchSessionPreferencesState struct {
	MinAgeCriteria     *minAgeCriteria      `json:"min_age"`
	MaxAgeCriteria     *maxAgeCriteria      `json:"max_age"`
	GenderCriteria     *genderCriteria      `json:"gender"`
	InterestsCriteria  *interestsCriteria   `json:"interests"`
	DistanceCriteria   *maxDistanceCriteria `json:"max_distance_km"`
	MaxWaitTimeSeconds int32                `json:"max_wait_time_seconds"`
}

type MatchPreferences struct {
	state *matchSessionPreferencesState
}

func (p *MatchPreferences) UnmarshalJSON(data []byte) error {
	tmpState := new(matchSessionPreferencesState)
	if err := json.Unmarshal(data, tmpState); err != nil {
		return err
	}

	if tmpState.MinAgeCriteria == nil {
		tmpState.MinAgeCriteria = &minAgeCriteria{minAllowedAge}
	}

	if tmpState.MaxAgeCriteria == nil {
		tmpState.MaxAgeCriteria = &maxAgeCriteria{maxAllowedAge}
	}

	if tmpState.MaxWaitTimeSeconds == 0 {
		tmpState.MaxWaitTimeSeconds = defaultMaxWaitSeconds
	}

	p.state = tmpState
	return nil
}

func (p MatchPreferences) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.state)
}

func (p *MatchPreferences) WithMaxWaitTimeSeconds(maxWaitTimeInSeconds int32) *MatchPreferences {
	if maxWaitTimeInSeconds <= 0 {
		maxWaitTimeInSeconds = defaultMaxWaitSeconds
	}

	p.state.MaxWaitTimeSeconds = maxWaitTimeInSeconds
	return p
}

func (p *MatchPreferences) WithMinAge(minAge int) *MatchPreferences {
	if minAge < minAllowedAge {
		minAge = minAllowedAge
	}

	p.state.MinAgeCriteria = &minAgeCriteria{minAge: minAge}
	return p
}

func (p *MatchPreferences) WithMaxAge(maxAge int) *MatchPreferences {
	if maxAge <= 0 || maxAge > maxAllowedAge {
		maxAge = maxAllowedAge
	}

	p.state.MaxAgeCriteria = &maxAgeCriteria{maxAge: maxAge}
	return p
}

func (p *MatchPreferences) WithGender(g gender.Gender) *MatchPreferences {
	if g.IsUnspecified() {
		return p
	}

	p.state.GenderCriteria = &genderCriteria{g}
	return p
}

func (p *MatchPreferences) WithInterests(interests []string) *MatchPreferences {
	if len(interests) == 0 {
		return p
	}

	p.state.InterestsCriteria = &interestsCriteria{interests}
	return p
}

func (p *MatchPreferences) WithMaxDistanceKm(loc *location.Location, maxDistanceKm float64) *MatchPreferences {
	if maxDistanceKm <= 0 {
		return p
	}

	if loc == nil {
		return p
	}

	p.state.DistanceCriteria = &maxDistanceCriteria{
		loc,
		maxDistanceKm,
	}
	return p
}

func (p MatchPreferences) MaxWaitTime() time.Duration {
	return time.Duration(p.state.MaxWaitTimeSeconds) * time.Second
}

func (p MatchPreferences) MinAge() int {
	if p.state.MinAgeCriteria == nil {
		return minAllowedAge
	}
	return p.state.MinAgeCriteria.minAge
}

func (p MatchPreferences) MaxAge() int {
	if p.state.MaxAgeCriteria == nil {
		return maxAllowedAge
	}
	return p.state.MaxAgeCriteria.maxAge
}

func (p MatchPreferences) Gender() gender.Gender {
	if p.state.GenderCriteria == nil {
		return gender.GenderUnspecified
	}
	return p.state.GenderCriteria.gender
}

func (p MatchPreferences) Location() *location.Location {
	if p.state.DistanceCriteria == nil {
		return nil
	}

	return p.state.DistanceCriteria.location
}

func (p MatchPreferences) Interests() []string {
	if p.state.InterestsCriteria == nil {
		return []string{}
	}
	return p.state.InterestsCriteria.interests
}

func (p MatchPreferences) MaxDistanceKm() float64 {
	if p.state.DistanceCriteria == nil {
		return 0
	}
	return p.state.DistanceCriteria.maxDistanceKm
}

func (p MatchPreferences) IsUserCompatible(matchUser *User) bool {
	if p.state.MinAgeCriteria != nil && !p.state.MinAgeCriteria.IsSatisfiedBy(matchUser) {
		return false
	}

	if p.state.MaxAgeCriteria != nil && !p.state.MaxAgeCriteria.IsSatisfiedBy(matchUser) {
		return false
	}

	if p.state.GenderCriteria != nil && !p.state.GenderCriteria.IsSatisfiedBy(matchUser) {
		return false
	}

	if p.state.InterestsCriteria != nil && !p.state.InterestsCriteria.IsSatisfiedBy(matchUser) {
		return false
	}

	if p.state.DistanceCriteria != nil && !p.state.DistanceCriteria.IsSatisfiedBy(matchUser) {
		return false
	}

	return true
}

func (p MatchPreferences) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	if p.state.MinAgeCriteria != nil {
		sb.WriteString(p.state.MinAgeCriteria.String())
		sb.WriteString(", ")
	}

	if p.state.MaxAgeCriteria != nil {
		sb.WriteString(p.state.MaxAgeCriteria.String())
		sb.WriteString(", ")
	}

	if p.state.GenderCriteria != nil {
		sb.WriteString(p.state.GenderCriteria.String())
		sb.WriteString(", ")
	}

	if p.state.InterestsCriteria != nil {
		sb.WriteString(p.state.InterestsCriteria.String())
		sb.WriteString(", ")
	}

	if p.state.DistanceCriteria != nil {
		sb.WriteString(p.state.DistanceCriteria.String())
		sb.WriteString(", ")
	}

	sb.WriteString("MaxWaitTime: ")
	sb.WriteString(strconv.Itoa(int(p.state.MaxWaitTimeSeconds)))
	sb.WriteString("}")
	return sb.String()
}

// minAgeCriteria represents the minimum age criteria for a match.
type minAgeCriteria struct {
	minAge int
}

func (m minAgeCriteria) Value() int {
	return m.minAge
}

func (m minAgeCriteria) IsSatisfiedBy(matchUser *User) bool {
	return m.minAge == 0 || matchUser.age >= m.minAge
}

func (m minAgeCriteria) String() string {
	return "MinAge: " + strconv.Itoa(m.minAge)
}

func (m *minAgeCriteria) UnmarshalJSON(data []byte) error {
	var minAge int
	if err := json.Unmarshal(data, &minAge); err != nil {
		return err
	}

	m.minAge = minAge
	return nil
}

func (m *minAgeCriteria) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.minAge)
}

// maxAgeCriteria represents the maximum age criteria for a match
type maxAgeCriteria struct {
	maxAge int
}

func (m maxAgeCriteria) Value() int {
	return m.maxAge
}

func (m maxAgeCriteria) IsSatisfiedBy(matchUser *User) bool {
	return m.maxAge == 0 || matchUser.age <= m.maxAge
}

func (m maxAgeCriteria) String() string {
	return "MaxAge: " + strconv.Itoa(m.maxAge)
}

func (m *maxAgeCriteria) UnmarshalJSON(data []byte) error {
	var maxAge int
	if err := json.Unmarshal(data, &maxAge); err != nil {
		return err
	}

	m.maxAge = maxAge
	return nil
}

func (m *maxAgeCriteria) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.maxAge)
}

// genderCriteria represents the filter criteria for matching mss.
type genderCriteria struct {
	gender gender.Gender
}

func (gc genderCriteria) Value() gender.Gender {
	return gc.gender
}

func (gc genderCriteria) IsSatisfiedBy(matchUser *User) bool {
	if gc.gender.IsUnspecified() {
		return true
	}

	return gc.gender.Is(matchUser.gender)
}

func (gc genderCriteria) String() string {
	return "Gender: " + gc.gender.String()
}

func (gc *genderCriteria) UnmarshalJSON(data []byte) error {
	var genderStr string
	if err := json.Unmarshal(data, &genderStr); err != nil {
		return err
	}

	gc.gender = gender.ParseString(genderStr)
	return nil
}

func (gc *genderCriteria) MarshalJSON() ([]byte, error) {
	return json.Marshal(gc.gender.String())
}

// interestsCriteria represents the filter criteria for matching mss.
type interestsCriteria struct {
	interests []string
}

func (ic interestsCriteria) Value() []string {
	return ic.interests
}

func (ic interestsCriteria) IsSatisfiedBy(matchUser *User) bool {
	if len(ic.interests) == 0 {
		return true
	}

	if len(matchUser.Preferences().Interests()) == 0 {
		return false
	}

	for _, interest := range ic.interests {
		for _, msInterest := range matchUser.Preferences().Interests() {
			if interest == msInterest {
				return true
			}
		}
	}

	return false
}

func (ic interestsCriteria) String() string {
	return "Interests: " + strings.Join(ic.interests, ", ")
}

func (ic *interestsCriteria) UnmarshalJSON(data []byte) error {
	var interests []string
	if err := json.Unmarshal(data, &interests); err != nil {
		return err
	}

	ic.interests = interests
	return nil
}

func (ic *interestsCriteria) MarshalJSON() ([]byte, error) {
	return json.Marshal(ic.interests)
}

// maxDistanceCriteria represents the filter criteria for matching mss.
type maxDistanceCriteria struct {
	location      *location.Location
	maxDistanceKm float64
}

func (mdc *maxDistanceCriteria) Value() float64 {
	return mdc.maxDistanceKm
}

func (mdc *maxDistanceCriteria) IsSatisfiedBy(matchUser *User) bool {
	if mdc.location == nil || mdc.maxDistanceKm == 0 {
		return true
	}

	if matchUser.Preferences().Location() == nil {
		return false
	}

	return withinDistance(
		*mdc.location,
		*matchUser.preferences.Location(),
		mdc.maxDistanceKm,
	)
}

func (mdc *maxDistanceCriteria) String() string {
	return "MaxDistanceKm: " + strconv.FormatFloat(mdc.maxDistanceKm, 'f', -1, 64)
}

func (mdc *maxDistanceCriteria) UnmarshalJSON(data []byte) error {
	var maxDistanceKm float64
	if err := json.Unmarshal(data, &maxDistanceKm); err != nil {
		return err
	}

	mdc.maxDistanceKm = maxDistanceKm
	return nil
}

func (mdc *maxDistanceCriteria) MarshalJSON() ([]byte, error) {
	return json.Marshal(mdc.maxDistanceKm)
}

func withinDistance(
	currentMatchSessionLocation, matchedSessionLocation location.Location,
	currentMatchSessionMaxDistanceKm float64,
) bool {
	distance, ok := currentMatchSessionLocation.DistanceTo(matchedSessionLocation)
	if !ok {
		return false
	}

	return distance <= currentMatchSessionMaxDistanceKm
}
