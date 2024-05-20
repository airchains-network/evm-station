package junction

import (
	"fmt"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"os"
	"path/filepath"
)

func NewKeyPair() (privateKeyX kyber.Scalar, publicKeyX kyber.Point) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey := suite.Scalar().Pick(suite.RandomStream())
	publicKey := suite.Point().Mul(privateKey, nil)
	return privateKey, publicKey
}

func SetVRFPubKey(pubKey string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error in getting home dir path: " + err.Error())
	}
	ConfigFilePath := filepath.Join(homeDir, BlockchainFolder)
	VRFPubKeyPath := filepath.Join(ConfigFilePath, VrfPubKeyFile)
	file, err := os.Create(VRFPubKeyPath)
	if err != nil {
		// Handle the error if the file cannot be created
		fmt.Println(fmt.Sprintf("error creating vrfPubKey.txt: %v", err))
		return
	}
	defer file.Close()

	// Write the stationId to the file
	_, err = file.WriteString(pubKey)
	if err != nil {
		// Handle the error if the file cannot be written to
		fmt.Println(fmt.Sprintf("error writing to vrfPubKey.txt: %v", err))
		return
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		// Handle the error if the file cannot be saved
		fmt.Println(fmt.Sprintf("error saving vrfPubKey.txt: %v", err))
		return
	}

	// Print the stationId
	fmt.Println(fmt.Sprintf("vrfPubKey ID: %s", pubKey))
}

func SetVRFPrivKey(privateKey string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error in getting home dir path: " + err.Error())
	}
	ConfigFilePath := filepath.Join(homeDir, BlockchainFolder)
	VRFPrivKeyPath := filepath.Join(ConfigFilePath, VRFPrivKeyFile)
	file, err := os.Create(VRFPrivKeyPath)
	if err != nil {
		// Handle the error if the file cannot be created
		fmt.Println(fmt.Sprintf("error creating vrfPrivKey.txt: %v", err))
		return
	}
	defer file.Close()

	// Write the stationId to the file
	_, err = file.WriteString(privateKey)
	if err != nil {
		// Handle the error if the file cannot be written to
		fmt.Println(fmt.Sprintf("error writing to vrfPrivKey.txt: %v", err))
		return
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		// Handle the error if the file cannot be saved
		fmt.Println(fmt.Sprintf("error saving vrfPrivKey.txt: %v", err))
		return
	}

	// Print the stationId
	fmt.Println(fmt.Sprintf("vrfPrivKey ID: %s", privateKey))
}
