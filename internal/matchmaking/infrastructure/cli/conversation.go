package matchsessioncli

import (
	"strconv"

	"github.com/spf13/cobra"

	matchpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/matchmaking/v1"
)

var NewMatchSessionCobraCommand = func(grpcClient matchpb.MatchMakingServiceClient) *cobra.Command {
	return &cobra.Command{
		Use:       "find-match",
		Short:     "Try to match another User based on the preferences",
		ValidArgs: []string{"userID"},
		Args:      cobra.MinimumNArgs(3),
		Run:       handleRootCmd(grpcClient),
	}
}

func handleRootCmd(grpcClient matchpb.MatchMakingServiceClient) func(*cobra.Command, []string) {
	return func(cobraCmd *cobra.Command, args []string) {
		ctx := cobraCmd.Context()

		userID := args[0]
		userName := args[1]
		userAgeStr := args[2]
		userAge64, err := strconv.ParseInt(userAgeStr, 10, 32)
		if err != nil {
			cobraCmd.PrintErr(err)
			return
		}
		userAge := int32(userAge64)

		userGender := func() string {
			if len(args) < 4 || args[3] == "" {
				return "UNSPECIFIED"
			}
			return args[3]
		}()

		res, err := grpcClient.FindMatch(
			ctx,
			&matchpb.FindMatchRequest{
				UserId:   userID,
				UserName: userName,
				UserAge:  userAge,
				MatchPreferences: &matchpb.MatchPreferences{
					Gender:             matchpb.Gender(matchpb.Gender_value[userGender]),
					MinAge:             18,
					MaxAge:             40,
					MaxDistanceKm:      100,
					MaxWaitTimeSeconds: 20,
				},
			})

		if err != nil {
			cobraCmd.PrintErr(err)
			return
		}

		cobraCmd.Printf("Findd new Match Session: '%s'\n", res.GetMatchId())
	}
}
