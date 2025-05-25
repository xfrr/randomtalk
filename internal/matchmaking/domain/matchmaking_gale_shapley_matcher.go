package matchdomain

import (
	"sort"
)

var _ StableMatchFinder = (*GaleShapleyService)(nil)

// GaleShapleyService is the implementation of the Gale-Shapley stable matching algorithm.
type GaleShapleyService struct{}

func NewGaleShapleyStableMatcher() *GaleShapleyService {
	return &GaleShapleyService{}
}

// GaleShapley builds preference lists using User's MatchPreferences
// and runs the Gale-Shapley algorithm, returning matches from setA to setB.
func (s *GaleShapleyService) FindStableMatches(setA, setB []*User) []int {
	nA, nB := len(setA), len(setB)
	if nA == 0 || nB == 0 {
		return nil
	}

	// build preference lists in terms of indexes:
	// preferencesA[aIndex] = [] of bIndexes in order of preference
	preferencesA := make([][]int, nA)
	for aIndex, aUser := range setA {
		candidates := s.findCompatible(aUser, setB)
		// Sort candidate IDs or apply other logic to define the order
		sort.Slice(candidates, func(i, j int) bool {
			return setB[candidates[i]].ID() < setB[candidates[j]].ID()
		})
		preferencesA[aIndex] = candidates
	}

	// build preference lists for B
	preferencesB := make([][]int, nB)
	for bIndex, bUser := range setB {
		candidates := s.findCompatible(bUser, setA)
		// Sort candidate IDs or apply other logic for preference ordering
		sort.Slice(candidates, func(i, j int) bool {
			return setA[candidates[i]].ID() < setA[candidates[j]].ID()
		})
		preferencesB[bIndex] = candidates
	}

	// we need "inverse ranking" for B, so we can quickly
	// see which A is preferred if B gets multiple proposals.
	rankB := make([][]int, nB)
	for bIndex := 0; bIndex < nB; bIndex++ {
		rankB[bIndex] = make([]int, nA)
		// initialize with some large value
		for aIndex := 0; aIndex < nA; aIndex++ {
			rankB[bIndex][aIndex] = nA + 1 // sentinel
		}
		// fill in actual ranks
		for rank, aIdx := range preferencesB[bIndex] {
			rankB[bIndex][aIdx] = rank
		}
	}

	// run the standard Gale-Shapley over sets A and B using these preference lists.
	//
	// matches[aIndex] = bIndex. Start unmatched = -1
	matches := make([]int, nA)
	for i := range matches {
		matches[i] = -1
	}

	// engagedTo[bIndex] = aIndex
	engagedTo := make([]int, nB)
	for i := range engagedTo {
		engagedTo[i] = -1
	}

	// track next proposal index for each a
	nextProposal := make([]int, nA)
	freeCount := nA

	for freeCount > 0 {
		var aIndex int
		// find a free man aIndex who still has candidates to propose
		for aIndex = 0; aIndex < nA; aIndex++ {
			if matches[aIndex] == -1 && nextProposal[aIndex] < len(preferencesA[aIndex]) {
				break
			}
		}

		// if we found none, break
		if aIndex == nA {
			break // no more proposals possible
		}

		// get next candidate bIndex to propose to
		bIndex := preferencesA[aIndex][nextProposal[aIndex]]
		nextProposal[aIndex]++

		if engagedTo[bIndex] == -1 {
			// B is free, engage
			engagedTo[bIndex] = aIndex
			matches[aIndex] = bIndex
			freeCount--
		} else {
			// B currently engaged
			currentA := engagedTo[bIndex]
			// If B prefers this new A over current
			if rankB[bIndex][aIndex] < rankB[bIndex][currentA] {
				// B ditches currentA
				matches[currentA] = -1
				engagedTo[bIndex] = aIndex
				matches[aIndex] = bIndex
			}
			// else reject aIndex (remain free), do nothing
		}
	}

	return matches
}

// findCompatible returns the indices of users in `others` that are
// mutually compatible with `user`. The result is a slice of indices
// referencing positions in `others`.
func (s *GaleShapleyService) findCompatible(user *User, others []*User) []int {
	var compatibleIndexes []int
	for idx, other := range others {
		if isMutuallyCompatible(user, other) {
			compatibleIndexes = append(compatibleIndexes, idx)
		}
	}
	return compatibleIndexes
}

// isMutuallyCompatible checks if 'user1' passes 'user2' preferences and vice versa.
func isMutuallyCompatible(u1, u2 *User) bool {
	if u1.ID() == u2.ID() {
		return false
	}

	return u1.Preferences().IsSatisfiedBy(u2) &&
		u2.Preferences().IsSatisfiedBy(u1)
}
