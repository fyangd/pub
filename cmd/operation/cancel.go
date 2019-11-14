package operation

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
	cancelOperationsArgs struct {
		Publisher          string
		Offer              string
		NotificationEmails string
	}

	// Canceller provides the ability to cancel an Operation
	Canceller interface {
		CancelOperation(ctx context.Context, params partner.CancelOperationParams) (string, error)
	}
)

func newCancelCommand(clientFactory func() (Canceller, error)) (*cobra.Command, error) {
	var oArgs cancelOperationsArgs
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "cancel the active operation for a given offer and print the operations",
		Run: xcobra.RunWithCtx(func(ctx context.Context, cmd *cobra.Command, args []string) {
			client, err := clientFactory()
			if err != nil {
				log.Fatalf("unable to create Cloud Partner Portal client: %v", err)
			}

			opLocation, err := client.CancelOperation(ctx, partner.CancelOperationParams{
				PublisherID:        oArgs.Publisher,
				OfferID:            oArgs.Offer,
				NotificationEmails: oArgs.NotificationEmails,
			})

			if err != nil {
				xcobra.PrintfErrAndExit(1, "unable to cancel the active operation: %v", err)
			}

			fmt.Print(opLocation)

		}),
	}

	if err := args.BindPublisher(cmd, &oArgs.Publisher); err != nil {
		return cmd, err
	}

	if err := args.BindOffer(cmd, &oArgs.Offer); err != nil {
		return cmd, err
	}

	cmd.Flags().StringVarP(&oArgs.NotificationEmails, "notification-emails", "e", "", "Comma separated list of emails to notify.")
	return cmd, nil
}
