package cmd

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/airchains-network/evm-station/junction"
	"github.com/airchains-network/evm-station/shared"
	"github.com/airchains-network/evm-station/station"
	"github.com/airchains-network/evm-station/types"
	"github.com/airchains-network/evm-station/zk"
	junctionTypes "github.com/airchains-network/junction/x/junction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

func ConnectJunctionClient(cmd *cobra.Command, homeDir string) error {
	// connect to junction if tracks is initialised
	seqConfig, err := getTracksData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Info("Tracks is not initialised, skipping junction client connection")
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

var TracksCommands = func() *cobra.Command {

	TracksCmd := &cobra.Command{
		Use:   "tracks",
		Short: "Interact with the tracks operations",
	}

	// Add commands
	// init tracks command
	InitTracksCmd.Flags().String("daRpc", "", "Description of DA RPC URL")
	InitTracksCmd.Flags().String("daKey", "", "Description of DA key")
	InitTracksCmd.Flags().String("daType", "", "Description of daType")
	InitTracksCmd.Flags().String("junctionRpc", "", "Description of Junction RPC URL")
	InitTracksCmd.Flags().String("junctionKeyName", "", "Description of Junction Key Name")
	err := InitTracksCmd.MarkFlagRequired("daRpc")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "daRpc", "error", err)
		return nil
	}
	err = InitTracksCmd.MarkFlagRequired("daKey")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "daKey", "error", err)
		return nil
	}
	err = InitTracksCmd.MarkFlagRequired("daType")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "daType", "error", err)
		return nil
	}
	err = InitTracksCmd.MarkFlagRequired("junctionRpc")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "junctionRpc", "error", err)
		return nil
	}
	err = InitTracksCmd.MarkFlagRequired("junctionKeyName")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "junctionKeyName", "error", err)
		return nil
	}
	TracksCmd.AddCommand(InitTracksCmd)

	// get tracks details: no extra flags required
	TracksCmd.AddCommand(TracksDetailsCmd)

	// Main command
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Balance-related operations",
	}

	// Subcommand "junction"
	JunctionBalanceCmd := &cobra.Command{
		Use:   "junction",
		Short: "Check balance of a junction account",
		Run: func(TracksCmd *cobra.Command, args []string) {
			CheckBalanceJunction(TracksCmd) // Pass jClient here
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
	TracksCmd.AddCommand(balanceCmd)

	// create station
	CreateStationCmd := &cobra.Command{
		Use:   "create-station",
		Short: "Create station code call to junction",
		Run:   CreateStationOnJunction,
	}
	CreateStationCmd.Flags().String("info", "", "Station information")
	CreateStationCmd.Flags().StringSlice("tracks", []string{}, "Array for tracks for new station")
	err = CreateStationCmd.MarkFlagRequired("info")
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in marking flag required", "flag", "info", "error", err)
		return nil
	}
	TracksCmd.AddCommand(CreateStationCmd)

	// start tracks. take data from 26657
	StartTracksCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the tracks",
		Run:   StartTracks,
	}
	TracksCmd.AddCommand(StartTracksCmd)

	return TracksCmd
}

func InitTracksConfigsCmd(cmd *cobra.Command) (*types.TracksConfigs, error) {
	var configs types.TracksConfigs
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

var InitTracksCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the tracks nodes",
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

		configs, err := InitTracksConfigsCmd(cmd)
		if err != nil {
			log.NewLogger(os.Stderr).Warn(err.Error())
			return
		}

		// create proving and verification key
		log.NewLogger(os.Stderr).Info("Creating proving and verification keys")
		err = zk.CreateVkPkNew(homeDir)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Error in creating vk and pk: " + err.Error())
			return
		}
		//log.NewLogger(os.Stderr).Info("Proving and Verification keys created successfully")

		// create wallet
		jKeyPath := filepath.Join(homeDir, JunctionKeysFolder)
		log.NewLogger(os.Stderr).Info("Creating account at", "Path", jKeyPath)
		err = junction.CreateAccount(configs.JunctionKeyName, jKeyPath, addressPrefix)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Error creating account", "err", err)
			return
		}

		tracksDir := filepath.Join(homeDir, EvmStationDir)
		sFile := filepath.Join(tracksDir, TracksFileName)
		file, err := os.Create(sFile)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Error creating file", "err", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.NewLogger(os.Stderr).Error("Error closing file", "err", err)
			}
		}(file)

		// Use the TOML encoder to write directly to the file
		if err := toml.NewEncoder(file).Encode(configs); err != nil {
			log.NewLogger(os.Stderr).Error("Error encoding configuration", "err", err)
			return
		} else {
			log.NewLogger(os.Stderr).Info("Tracks configuration saved successfully", "Path", sFile)
			return
		}
		// create proving and verification key
	},
}

// get tracks details
// getTracksData function gets the tracks configuration data
func getTracksData(cmd *cobra.Command) (*types.TracksConfigs, error) {

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}

	homeDir := clientCtx.HomeDir
	if homeDir == "" {
		return nil, fmt.Errorf("failed to get user home directory")
	}

	tracksDir := filepath.Join(homeDir, EvmStationDir)
	sFile := filepath.Join(tracksDir, TracksFileName)

	if _, err = os.Stat(sFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("tracks configuration file does not exist")
	} else {
		f, err := os.Open(sFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.NewLogger(os.Stderr).Error("Error closing file", "err", err)
			}
		}(f)

		var configs types.TracksConfigs

		if _, err = toml.DecodeReader(f, &configs); err != nil {
			return nil, fmt.Errorf("error decoding tracks configuration file: %w", err)
		} else {
			return &configs, nil
		}
	}
}

func GetTracksDetails(cmd *cobra.Command, args []string) {

	configs, err := getTracksData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	log.NewLogger(os.Stderr).Info("Tracks configuration Data")
	value := reflect.ValueOf(*configs) // Dereference the pointer
	var field reflect.StructField
	for i := 0; i < value.NumField(); i++ {
		field = value.Type().Field(i)
		log.NewLogger(os.Stderr).Info("TracksInfo", field.Name, fmt.Sprint(value.Field(i).Interface()))
	}
}

var TracksDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get the details of the tracks",
	Run:   GetTracksDetails,
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

	configs, err := getTracksData(cmd)
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

	TracksConfigs, err := getTracksData(cmd)
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

	accountName := TracksConfigs.JunctionKeyName
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
		DaType:    TracksConfigs.DaType,
		Prover:    Prover,
	}

	stationId := uuid.New().String()
	success := junction.CreateStation(*jClient, keyringDir, accountName, addr, extraArg, stationId, stationInfo, verificationKey, tracks, homeDir)
	if !success {
		log.NewLogger(os.Stderr).Error("Failed to create station")
		return
	}

	log.NewLogger(os.Stderr).Info("Successfully created station")
	// save stationId in tracks configuration
	TracksConfigs.StationId = stationId

	tracksDir := filepath.Join(homeDir, EvmStationDir)
	sFile := filepath.Join(tracksDir, TracksFileName)
	log.NewLogger(os.Stderr).Info("Writing configuration", "Path", sFile)
	file, err := os.Create(sFile)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error creating file", "err", err)
		return
	}
	defer file.Close()

	// Use the TOML encoder to write directly to the file
	if err := toml.NewEncoder(file).Encode(TracksConfigs); err != nil {
		log.NewLogger(os.Stderr).Error("Error encoding configuration", "err", err)
		return
	} else {
		log.NewLogger(os.Stderr).Info("Included Station Id in Tracks configuration successfully.")
		return
	}
}

func StartTracks(cmd *cobra.Command, args []string) {

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
	//TracksDir := filepath.Join(homeDir, TracksDbDir)

	// junction client connection
	err = ConnectJunctionClient(cmd, homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}

	TracksConfigs, err := getTracksData(cmd)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return
	}
	if TracksConfigs.StationId == "" {
		log.NewLogger(os.Stderr).Error("create station before stating sequencer")
		return
	}

	//// Initialise or Check database for Tracks
	//success := tracksdb.ConfigDb(TracksDir)
	//if !success {
	//	log.NewLogger(os.Stderr).Error("Failed to initialize Tracks Database")
	//	return
	//}

	// get VRF private and public keys
	VRFPrivateKeyStr, err := junction.GetVRFPrivKey(homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to get VRF private key")
		return
	}
	if VRFPrivateKeyStr == "" {
		log.NewLogger(os.Stderr).Error("VRF private key is empty")
		return
	}

	VRFPublicKeyStr, err := junction.GetVRFPubKey(homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to get VRF public key")
		return
	}
	if VRFPublicKeyStr == "" {
		log.NewLogger(os.Stderr).Error("VRF public key is empty")
		return
	}

	// Start remote signer (must start before node if running builtin).
	log.NewLogger(os.Stderr).Info("Starting Tracks", "ChainId", TracksConfigs.StationId, "JunctionRpc", TracksConfigs.JunctionRpc, "JunctionKeyName", TracksConfigs.JunctionKeyName)

	stationId := TracksConfigs.StationId
	accountName := TracksConfigs.JunctionKeyName
	keyringDir := filepath.Join(homeDir, JunctionKeysFolder)
	addr, err := junction.CheckIfAccountExists(accountName, keyringDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in getting account address: " + err.Error())
		return
	}
	account, err := junction.GetCosmosAccount(accountName, keyringDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in getting account: " + err.Error())
		return
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

	ctx := shared.GetContext()
	for {
		latestPodInStation := station.QueryLatestPodNumber() // latest pod number"
		log.NewLogger(os.Stderr).Info("latest pod on Station", "podNumber", latestPodInStation)

		latestVerifiedPodInJunction := junction.QueryLatestVerifiedBatch(*jClient, ctx, stationId)
		log.NewLogger(os.Stderr).Info("latest verified pod on junction", "podNumber", latestVerifiedPodInJunction)
		podNumberToProcess := latestVerifiedPodInJunction + 1

		if latestPodInStation == 0 {
			log.NewLogger(os.Stderr).Debug("No Pods at all on Station. Waiting for new pods", "waitTime", "10 seconds")
			time.Sleep(10 * time.Second)
		} else if latestPodInStation >= podNumberToProcess {
			log.NewLogger(os.Stderr).Info("Processing New Pod", "podNumber", podNumberToProcess)

			log.NewLogger(os.Stderr).Info("Processing InitVRf transaction")
			success := junction.InitVRF(podNumberToProcess, ctx, *jClient, account, addr, stationId, VRFPrivateKeyStr, VRFPublicKeyStr)
			if success == false {
				log.NewLogger(os.Stderr).Error("Failed to Init VRF due to above error")
				return
			}

			log.NewLogger(os.Stderr).Info("Processing ValidateVRF transaction")
			success = junction.ValidateVRF(podNumberToProcess, ctx, *jClient, account, addr, stationId, VRFPrivateKeyStr, VRFPublicKeyStr)
			if success == false {
				log.NewLogger(os.Stderr).Error("Failed to Validate VRF due to above error")
				return
			}

			// Query the pod from Station
			pod, err := station.QueryPodNumber(podNumberToProcess)
			if err != nil {
				log.NewLogger(os.Stderr).Error("Failed to get pod from Station", "podNumber", podNumberToProcess, "error", err)
				return
			}

			var previousMerkleRootHash string
			if podNumberToProcess > 1 {
				// Query the pod from Station
				previousPodNumber := podNumberToProcess - 1
				previousPod, err := station.QueryPodNumber(previousPodNumber)
				if err != nil {
					log.NewLogger(os.Stderr).Error("Failed to get pod from Station", "podNumber", previousPodNumber, "error", err)
					return
				}
				previousMerkleRootHash = previousPod.MerkleRootHash
			} else {
				previousMerkleRootHash = ""
			}

			log.NewLogger(os.Stderr).Info("Processing SubmitPod transaction")
			success = junction.SubmitPod(podNumberToProcess, ctx, *jClient, account, addr, stationId, previousMerkleRootHash, pod.MerkleRootHash, pod.Witness, pod.Proof)
			if success == false {
				log.NewLogger(os.Stderr).Error("Failed to Submit Pod due to above error")
				return
			}

			log.NewLogger(os.Stderr).Info("Processing VerifyPod transaction")
			success = junction.VerifyPod(podNumberToProcess, ctx, *jClient, account, addr, stationId, previousMerkleRootHash, pod.MerkleRootHash, pod.Proof)
			if success == false {
				log.NewLogger(os.Stderr).Error("Failed to Verify Pod due to above error")
				return
			}
		} else {
			log.NewLogger(os.Stderr).Debug("No New Pods on Station. Waiting for new pods", "waitTime", "10 seconds")
			time.Sleep(10 * time.Second)
		}
	}

}
