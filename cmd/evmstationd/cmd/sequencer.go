package cmd

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/airchains-network/evm-station/junction"
	"github.com/airchains-network/evm-station/shared"
	"github.com/airchains-network/evm-station/types"
	"github.com/airchains-network/evm-station/zk"
	junctionTypes "github.com/airchains-network/junction/x/junction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"reflect"
)

func ConnectJunctionClient(cmd *cobra.Command, homeDir string) error {
	// connect to junction if sequencer is initialised
	seqConfig, err := getSequencerData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Info("Sequencer is not initialised, skipping junction client connection")
		return err
	} else {
		junctionRpc := seqConfig.JunctionRpc
		if homeDir == "" {
			log.NewLogger(os.Stderr).Warn("Failed to get user home directory")
			return err
		}
		keyHomePath := filepath.Join(homeDir, JunctionKeysFolder)
		j := shared.JunctionNewConfig(junctionRpc, keyHomePath)
		shared.SetJunctionClient(j)
	}
	return nil
}

var sequencerCommands = func() *cobra.Command {
	sequencerCmd := &cobra.Command{
		Use:   "sequencer",
		Short: "Interact with the sequencer operations",
	}

	// Add commands
	// init sequencer command
	InitSequencerCmd.Flags().String("daRpc", "", "Description of DA RPC URL")
	InitSequencerCmd.Flags().String("daKey", "", "Description of DA key")
	InitSequencerCmd.Flags().String("daType", "", "Description of daType")
	InitSequencerCmd.Flags().String("junctionRpc", "", "Description of Junction RPC URL")
	InitSequencerCmd.Flags().String("junctionKeyName", "", "Description of Junction Key Name")
	InitSequencerCmd.MarkFlagRequired("daRpc")
	InitSequencerCmd.MarkFlagRequired("daKey")
	InitSequencerCmd.MarkFlagRequired("daType")
	InitSequencerCmd.MarkFlagRequired("junctionRpc")
	InitSequencerCmd.MarkFlagRequired("junctionKeyName")
	sequencerCmd.AddCommand(InitSequencerCmd)

	// get sequencer details: no extra flags required
	sequencerCmd.AddCommand(sequencerDetailsCmd)

	// Main command
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Balance-related operations",
	}

	// Subcommand "junction"
	JunctionBalanceCmd := &cobra.Command{
		Use:   "junction",
		Short: "Check balance of a junction account",
		Run: func(sequencerCmd *cobra.Command, args []string) {
			CheckBalanceJunction(sequencerCmd) // Pass jClient here
		},
	}

	// Subcommand "da"
	DaCmd := &cobra.Command{
		Use:   "da",
		Short: "To be implemented",
		Run: func(cmd *cobra.Command, args []string) {
			log.NewLogger(os.Stderr).Warn("Da Balance Command not implemented yet")
		},
	}

	// Add the subcommands to the main command
	balanceCmd.AddCommand(JunctionBalanceCmd)
	balanceCmd.AddCommand(DaCmd)

	// Add the main command to the root command
	sequencerCmd.AddCommand(balanceCmd)

	// create station
	CreateStationCmd := &cobra.Command{
		Use:   "create-station",
		Short: "Create station code call to junction",
		Run:   CreateStationOnJunction,
	}
	CreateStationCmd.Flags().String("info", "", "Station information")
	CreateStationCmd.Flags().StringSlice("tracks", []string{}, "Array for tracks for new station")
	CreateStationCmd.MarkFlagRequired("info")
	sequencerCmd.AddCommand(CreateStationCmd)

	return sequencerCmd
}

func InitSequencerConfigsCmd(cmd *cobra.Command) (*types.SequencerConfigs, error) {
	var configs types.SequencerConfigs
	var err error

	configs.DaRPC, err = cmd.Flags().GetString("daRpc")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daRpc': %w", err)
	}

	configs.DaKey, err = cmd.Flags().GetString("daKey")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daKey': %w", err)
	}

	configs.DaType, err = cmd.Flags().GetString("daType")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daType': %w", err)
	}

	configs.JunctionRpc, err = cmd.Flags().GetString("junctionRpc")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'junctionRpc': %w", err)
	}

	configs.JunctionKeyName, err = cmd.Flags().GetString("junctionKeyName")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'junctionKeyName': %w", err)
	}

	configs.StationId = "" // will assign while creating station in further steps

	return &configs, nil
}

var InitSequencerCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the sequencer nodes",
	Run: func(cmd *cobra.Command, args []string) {

		clientCtx, err := client.GetClientQueryContext(cmd)
		if err != nil {
			log.NewLogger(os.Stderr).Error(err.Error())
			return
		}

		homeDir := clientCtx.HomeDir
		if homeDir == "" {
			log.NewLogger(os.Stderr).Warn("Failed to get user home directory")
			return
		}

		configs, err := InitSequencerConfigsCmd(cmd)
		if err != nil {
			log.NewLogger(os.Stderr).Warn(err.Error())
			return
		}

		jKeyPath := filepath.Join(homeDir, JunctionKeysFolder)
		log.NewLogger(os.Stderr).Info("Creating account at", "Path", jKeyPath)

		// create wallet
		err = junction.CreateAccount(configs.JunctionKeyName, jKeyPath, addressPrefix)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Error creating account", "err", err)
			return
		}

		tracksDir := filepath.Join(homeDir, EvmStationDir)
		sFile := filepath.Join(tracksDir, SequencerFileName)
		log.NewLogger(os.Stderr).Info("Writing configuration", "Path", sFile)
		file, err := os.Create(sFile)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Error creating file", "err", err)
			return
		}
		defer file.Close()

		// Use the TOML encoder to write directly to the file
		if err := toml.NewEncoder(file).Encode(configs); err != nil {
			log.NewLogger(os.Stderr).Error("Error encoding configuration", "err", err)
			return
		} else {
			err = zk.CreateVkPkNew(homeDir)
			if err != nil {
				log.NewLogger(os.Stderr).Error("Error in creating vk and pk: " + err.Error())
				return
			}

			log.NewLogger(os.Stderr).Info("Sequencer configuration saved successfully")
			return
		}
		// create proving and verification key
	},
}

// get sequencer details
// getSequencerData function gets the sequencer configuration data
func getSequencerData(cmd *cobra.Command) (*types.SequencerConfigs, error) {

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}

	homeDir := clientCtx.HomeDir
	if homeDir == "" {
		return nil, fmt.Errorf("failed to get user home directory")
	}

	tracksDir := filepath.Join(homeDir, EvmStationDir)
	sFile := filepath.Join(tracksDir, SequencerFileName)

	if _, err = os.Stat(sFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("sequencer configuration file does not exist")
	} else {
		f, err := os.Open(sFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		var configs types.SequencerConfigs

		if _, err = toml.DecodeReader(f, &configs); err != nil {
			return nil, fmt.Errorf("error decoding sequencer configuration file: %w", err)
		} else {
			return &configs, nil
		}
	}
}

func GetSequencerDetails(cmd *cobra.Command, args []string) {

	configs, err := getSequencerData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	log.NewLogger(os.Stderr).Info("Sequencer configuration Data")
	value := reflect.ValueOf(*configs) // Dereference the pointer
	var field reflect.StructField
	for i := 0; i < value.NumField(); i++ {
		field = value.Type().Field(i)
		log.NewLogger(os.Stderr).Info("SequencerInfo", field.Name, fmt.Sprint(value.Field(i).Interface()))
	}
}

var sequencerDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get the details of the sequencer",
	Run:   GetSequencerDetails,
}

func CheckBalanceJunction(cmd *cobra.Command) {

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	homeDir := clientCtx.HomeDir
	if homeDir == "" {
		log.NewLogger(os.Stderr).Error("Failed to get user home directory")
		return
	}

	err = ConnectJunctionClient(cmd, homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	configs, err := getSequencerData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	accountName := configs.JunctionKeyName
	keyringDir := filepath.Join(homeDir, JunctionKeysFolder)
	addr, err := junction.CheckIfAccountExists(accountName, keyringDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in getting account address: " + err.Error())
		return
	}

	jClient, jConnected := shared.GetJunctionClient()
	if !jConnected {
		log.NewLogger(os.Stderr).Error("Junction client not connected")
		return
	}

	haveBalance, value, err := junction.CheckBalance(jClient, addr)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in checking balance", "Address", addr, "Error", err.Error())
		return
	} else if !haveBalance {
		log.NewLogger(os.Stderr).Warn("No Balance", "Address", addr)
		return
	}

	log.NewLogger(os.Stderr).Info("Balance", "Address", addr, "Amount(amf)", value)
	return
}

func CreateStationOnJunction(cmd *cobra.Command, args []string) {

	sequencerConfigs, err := getSequencerData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
	}
	homeDir := clientCtx.HomeDir

	err = ConnectJunctionClient(cmd, homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	//get values from all flags
	info, _ := cmd.Flags().GetString("info")
	if info == "" {
		log.NewLogger(os.Stderr).Error("info is required")
		return
	}

	accountName := sequencerConfigs.JunctionKeyName
	keyringDir := filepath.Join(homeDir, JunctionKeysFolder)
	addr, err := junction.CheckIfAccountExists(accountName, keyringDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in getting account address: " + err.Error())
		return
	}

	tracks, _ := cmd.Flags().GetStringSlice("tracks")
	if len(tracks) == 0 {
		log.NewLogger(os.Stderr).Info("--tracks is not provided so taking default account as track")
		tracks = append(tracks, addr)
	}

	jClient, jConnected := shared.GetJunctionClient()
	if !jConnected {
		log.NewLogger(os.Stderr).Error("Junction client not connected")
		return
	}

	haveBalance, value, err := junction.CheckBalance(jClient, addr)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in checking balance of"+addr, "Error", err.Error())
		return
	} else if !haveBalance {
		log.NewLogger(os.Stderr).Warn("Not have balance in " + addr)
		return
	}

	log.NewLogger(os.Stderr).Info("Junction Account Balance (in amf):", "account", addr, "balance", value)

	// get verification key
	provingKey, verificationKey, err := zk.GetVkPk(homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
	}
	_ = provingKey // its not required to create junction

	stationInfo := types.StationInfo{
		StationType: StationType,
	}

	extraArg := junctionTypes.StationArg{
		TrackType: TrackType,
		DaType:    sequencerConfigs.DaType,
		Prover:    Prover,
	}

	fmt.Println(keyringDir)

	stationId := uuid.New().String()
	success := junction.CreateStation(*jClient, keyringDir, accountName, addr, extraArg, stationId, stationInfo, verificationKey, tracks)
	if !success {
		log.NewLogger(os.Stderr).Error("Failed to create station")
		return
	}

	log.NewLogger(os.Stderr).Info("Successfully created station")
	// save stationId in sequencer configuration
	sequencerConfigs.StationId = stationId

	tracksDir := filepath.Join(homeDir, EvmStationDir)
	sFile := filepath.Join(tracksDir, SequencerFileName)
	log.NewLogger(os.Stderr).Info("Writing configuration", "Path", sFile)
	file, err := os.Create(sFile)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error creating file", "err", err)
		return
	}
	defer file.Close()

	// Use the TOML encoder to write directly to the file
	if err := toml.NewEncoder(file).Encode(sequencerConfigs); err != nil {
		log.NewLogger(os.Stderr).Error("Error encoding configuration", "err", err)
		return
	} else {
		log.NewLogger(os.Stderr).Info("Included Station Id in Sequencer configuration successfully.")
		return
	}
}
