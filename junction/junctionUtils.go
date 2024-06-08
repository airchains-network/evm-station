package junction

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"io"
	"os"
	"path/filepath"
)

func NewKeyPair() (privateKeyX kyber.Scalar, publicKeyX kyber.Point) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey := suite.Scalar().Pick(suite.RandomStream())
	publicKey := suite.Point().Mul(privateKey, nil)
	return privateKey, publicKey
}

func SetVRFPubKey(pubKey, homeDir string) {
	VRFPubKeyPath := filepath.Join(homeDir, VrfPubKeyFile)
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

func SetVRFPrivKey(privateKey, homeDir string) {
	VRFPrivKeyPath := filepath.Join(homeDir, VRFPrivKeyFile)
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

func GetVRFPrivKey(homeDir string) (privateKey string, err error) {
	VRFPrivKeyPath := filepath.Join(homeDir, VRFPrivKeyFile)
	file, err := os.Open(VRFPrivKeyPath)
	if err != nil {
		// Handle the error if the file cannot be opened
		return "", fmt.Errorf("error opening vrfPrivKey.txt: %v", err)
	}
	defer file.Close()
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle the error if the file cannot be read
			log.Err(err).Msg("error reading vrfPrivKey.txt")
			return "", err
		}
		privateKey = string(buf[:n])
	}

	return privateKey, nil
}

func GetVRFPubKey(homeDir string) (publicKey string, err error) {
	VRFPubKeyPath := filepath.Join(homeDir, VrfPubKeyFile)
	file, err := os.Open(VRFPubKeyPath)
	if err != nil {
		// Handle the error if the file cannot be opened
		return "", fmt.Errorf("error opening vrfPubKey.txt: %v", err)
	}
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle the error if the file cannot be read
			log.Err(err).Msg("error reading vrfPrivKey.txt")
			return "", err
		}
		publicKey = string(buf[:n])
	}

	return publicKey, nil
}
