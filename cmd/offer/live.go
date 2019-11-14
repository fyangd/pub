package offer

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/devigned/pub/cmd/args"
	"github.com/devigned/pub/pkg/partner"
	"github.com/devigned/pub/pkg/xcobra"
)

type (
	goLiveOfferArgs struct {
		Publisher          string
		Offer              string
		NotificationEmails string
	}

	// Liver provides the ability to make an offer go live
	Liver interface {
		GoLiveWithOffer(ctx context.Context, params partner.GoLiveParams) (string, error)
	}
)

func newGoLiveCommand(clientFactory func() (Liver, error)) (*cobra.Command, error) {
	var oArgs goLiveOfferArgs
	cmd := &cobra.Command{
		Use:   "live",
		Short: "go live with an offer (make available to the world)",
		Run: xcobra.RunWithCtx(func(ctx context.Context, cmd *cobra.Command, args []string) {
			client, err := clientFactory()
			if err != nil {
				log.Fatalf("unable to create Cloud Partner Portal client: %v", err)
			}

			opLocation, err := client.GoLiveWithOffer(ctx, partner.GoLiveParams{
				NotificationEmails: oArgs.NotificationEmails,
				OfferID:            oArgs.Offer,
				PublisherID:        oArgs.Publisher,
			})

			if err != nil {
				xcobra.PrintfErrAndExit(1, "%v\n", err)
			}

			fmt.Println(opLocation)
		}),
	}

	if err := args.BindPublisher(cmd, &oArgs.Publisher); err != nil {
		return cmd, err
	}

	if err := args.BindOffer(cmd, &oArgs.Offer); err != nil {
		return cmd, err
	}

	cmd.Flags().StringVarP(&oArgs.NotificationEmails, "notification-emails", "e", "", "Comma separated list of emails to notify when publication completes.")
	return cmd, nil
}
