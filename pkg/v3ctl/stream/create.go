package stream

import (
	"fmt"
	"strings"

	"github.com/v3io/v3ctl/pkg/v3ctl"

	"github.com/nuclio/errors"
	"github.com/spf13/cobra"
	v3io "github.com/v3io/v3io-go/pkg/dataplane"
)

type createStreamCommandeer struct {
	*v3ctl.CreateCommandeer
	shardCount           int
	retentionPeriodHours int
}

func newCreateStreamCommandeer(createCommandeer *v3ctl.CreateCommandeer) (*createStreamCommandeer, error) {
	commandeer := &createStreamCommandeer{
		CreateCommandeer: createCommandeer,
	}

	cmd := &cobra.Command{
		Use:   "stream name",
		Short: "Create a data stream",
		RunE: func(cmd *cobra.Command, args []string) error {

			// if we got positional arguments
			if len(args) != 1 {
				return errors.New("Stream create requires a stream path")
			}

			streamPath := args[0]

			// must end with "/"
			if !strings.HasSuffix(streamPath, "/") {
				streamPath += "/"
			}

			// initialize root
			if err := createCommandeer.Initialize(); err != nil {
				return errors.Wrap(err, "Failed to initialize root")
			}

			createStreamInput := v3io.CreateStreamInput{}
			createStreamInput.Path = streamPath
			createStreamInput.ShardCount = commandeer.shardCount
			createStreamInput.RetentionPeriodHours = commandeer.retentionPeriodHours

			err := createCommandeer.RootCommandeer.Container.CreateStreamSync(&createStreamInput)
			if err != nil {
				return errors.Wrap(err, "Failed to get create stream")
			}

			fmt.Printf("Stream %s created successfully\n", streamPath)
			return nil
		},
	}

	cmd.Flags().IntVar(&commandeer.shardCount, "shard-count", 1, "Number of shards in the stream")
	cmd.Flags().IntVar(&commandeer.retentionPeriodHours, "retention-period", 1, "Retention period of the stream, in hours")

	createRecordCommandeer, err := newCreateRecordCommandeer(commandeer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create create record")
	}

	cmd.AddCommand(createRecordCommandeer.Cmd)

	commandeer.Cmd = cmd

	return commandeer, nil
}

// register to factory
func init() {
	v3ctl.CreateCommandeerRegistrySingleton.Register("stream",
		func(createCommandeer *v3ctl.CreateCommandeer) (*cobra.Command, error) {
			newCreateStreamCommandeer, err := newCreateStreamCommandeer(createCommandeer)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to create commandeer")
			}

			return newCreateStreamCommandeer.Cmd, nil
		})
}
