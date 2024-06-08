package tracksdb

import (
	"cosmossdk.io/log"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/evm-station/types"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path/filepath"
)

var (
	txDbInstance            *leveldb.DB
	blockDbInstance         *leveldb.DB
	staticDbInstance        *leveldb.DB
	stateDbInstance         *leveldb.DB
	batchesDbInstance       *leveldb.DB
	proofDbInstance         *leveldb.DB
	publicWitnessDbInstance *leveldb.DB
	daDbInstance            *leveldb.DB
	mockDbInstance          *leveldb.DB
	TracksDir               string
)

func ConfigDb(tracksDir string) bool {
	TracksDir = tracksDir
	success := InitDb()
	if !success {
		log.NewLogger(os.Stderr).Error("Failed to initialize database")
		return false
	} else {
		log.NewLogger(os.Stderr).Info("Tracks Database Config Checked")
		return true
	}
}

// InitDb This function  initializes different databases and returns true if all of them are successfully initialized, otherwise it returns false.
func InitDb() bool {
	if !InitTxDb() {
		return false
	}
	if !InitBlockDb() {
		return false
	}
	if !InitStaticDb() {
		return false
	}
	if !InitStateDb() {
		return false
	}
	if !InitBatchesDb() {
		return false
	}
	if !InitProofDb() {
		return false
	}
	if !InitPublicWitnessDb() {
		return false
	}
	if !InitDaDb() {
		return false
	}
	if !InitMockDb() {
		return false
	}
	return true
}

// InitTxDb This function initializes a LevelDB database for transactions and returns a boolean indicating
// whether the initialization was successful.
func InitTxDb() bool {

	filePath := filepath.Join(TracksDir, TxDbName)

	txDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open transaction LevelDB:", err)
		return false
	}
	txDbInstance = txDB

	txnNumberByte, err := txDbInstance.Get([]byte("txnCount"), nil)
	if txnNumberByte == nil || err != nil {
		err = txDbInstance.Put([]byte("txnCount"), []byte("0"), nil)
		if err != nil {
			log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in saving txnCount in txnDb : %s", err.Error()))
			return false
		}
	}
	return true
}

// InitBlockDb This function initializes a LevelDB database for storing blocks and returns a boolean indicating
// whether the initialization was successful.
func InitBlockDb() bool {

	filePath := filepath.Join(TracksDir, BlockDbName)
	blockDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open block LevelDB:", err)
		return false
	}
	blockDbInstance = blockDB

	blockNumberByte, err := blockDB.Get([]byte("blockCount"), nil)
	if blockNumberByte == nil || err != nil {
		err = blockDB.Put([]byte("blockCount"), []byte("0"), nil)
		if err != nil {
			log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in saving blockCount in blockDatabase : %s", err.Error()))
			//return false
			os.Exit(0)
		}
	}

	return true
}

// InitStaticDb This function initializes a static LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitStaticDb() bool {

	filePath := filepath.Join(TracksDir, StaticDbName)
	staticDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open static LevelDB:", err)
		return false
	}
	staticDbInstance = staticDB
	return true
}

func InitStateDb() bool {

	filePath := filepath.Join(TracksDir, StateDbName)
	stateDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open state LevelDB:", err)
		return false
	}

	stateDbInstance = stateDB

	podStateByte, err := stateDB.Get([]byte("podState"), nil)
	if podStateByte == nil || err != nil {

		emptyPodState := types.PodState{
			LatestPodHeight:     1,
			LatestTxState:       "InitVRF",
			LatestPodHash:       nil,
			PreviousPodHash:     nil,
			LatestPodProof:      nil,
			LatestPublicWitness: nil,
			Votes:               make(map[string]types.Votes),
			TracksAppHash:       nil,
			Batch:               nil,
			MasterTrackAppHash:  nil,
		}
		byteEmptyPodState, err := json.Marshal(emptyPodState)
		if err != nil {

			log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in marshalling emptyPodState : %s", err.Error()))
			return false
		}

		err = stateDB.Put([]byte("podState"), byteEmptyPodState, nil)
		if err != nil {
			log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in saving podState in pod database : %s", err.Error()))
			return false
		}

	}

	return true
}

// InitBatchesDb This function initializes a batches LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitBatchesDb() bool {

	filePath := filepath.Join(TracksDir, BatchesDbName)
	batchesDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open batches LevelDB:", err)
		return false
	}
	batchesDbInstance = batchesDB
	return true
}

// InitProofDb This function initializes a proof LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitProofDb() bool {

	filePath := filepath.Join(TracksDir, ProofDbName)
	proofDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open proof LevelDB:", err)
		return false
	}
	proofDbInstance = proofDB
	return true
}

func InitPublicWitnessDb() bool {

	filePath := filepath.Join(TracksDir, PublicWitness)
	publicWitnessDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open publicWitness LevelDB:", err)
		return false
	}
	publicWitnessDbInstance = publicWitnessDB
	return true
}

func InitDaDb() bool {

	filePath := filepath.Join(TracksDir, DaDbName)
	daDB, err := leveldb.OpenFile(filePath, nil)
	da := types.DAStruct{
		DAKey:             "0",
		DAClientName:      "0",
		BatchNumber:       "0",
		PreviousStateHash: "0",
		CurrentStateHash:  "0",
	}

	daBytes, err := json.Marshal(da)
	if err != nil {
		log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in marshalling da : %s", err.Error()))
		return false
	}

	daDbInstance = daDB
	daBytes, err = daDbInstance.Get([]byte("batch_0"), nil)
	if daBytes == nil || err != nil {
		err = daDbInstance.Put([]byte("batch_0"), daBytes, nil)
		if err != nil {
			log.NewLogger(os.Stderr).Error(fmt.Sprintf("Error in saving daBytes in da Database : %s", err.Error()))
			return false
		}
	}

	return true
}
func InitMockDb() bool {

	filePath := filepath.Join(TracksDir, MockDbName)
	mockDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to open mock LevelDB:", err)
		return false
	}
	mockDbInstance = mockDB
	return true
}

// GetTxDbInstance This function returns the instance of the air-leveldb database.
func GetTxDbInstance() *leveldb.DB {
	return txDbInstance
}

// GetBlockDbInstance This function returns the instance of the block database.
func GetBlockDbInstance() *leveldb.DB {
	return blockDbInstance
}

// GetStaticDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitStaticDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetStaticDbInstance() *leveldb.DB {
	return staticDbInstance
}

func GetStateDbInstance() *leveldb.DB {
	return stateDbInstance
}

// GetBatchesDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitBatchesDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetBatchesDbInstance() *leveldb.DB {
	return batchesDbInstance
}

// GetProofDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitProofDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetProofDbInstance() *leveldb.DB {
	return proofDbInstance
}

func GetPublicWitnessDbInstance() *leveldb.DB {
	return publicWitnessDbInstance
}

func GetDaDbInstance() *leveldb.DB {
	return daDbInstance
}

func GetMockDbInstance() *leveldb.DB {
	return mockDbInstance
}
