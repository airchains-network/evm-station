package junction

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/airchains-network/evm-station/types"
	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"io"
	"math/big"
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

func SerializeRequestCommitmentV2Plus(rc types.RequestCommitmentV2Plus) ([]byte, error) {
	var buf bytes.Buffer

	// Encode the blockNum
	err := binary.Write(&buf, binary.BigEndian, rc.BlockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to encode blockNum: %w", err)
	}

	// Encode the stationId as a fixed size or prefixed with its length
	// Here, we choose to prefix with length for simplicity
	if err := binary.Write(&buf, binary.BigEndian, uint64(len(rc.StationId))); err != nil {
		return nil, fmt.Errorf("failed to encode stationId length: %w", err)
	}
	buf.WriteString(rc.StationId)

	// Encode the upperBound
	err = binary.Write(&buf, binary.BigEndian, rc.UpperBound)
	if err != nil {
		return nil, fmt.Errorf("failed to encode upperBound: %w", err)
	}

	// Encode the requesterAddress as a fixed size or prefixed with its length
	if err := binary.Write(&buf, binary.BigEndian, uint64(len(rc.RequesterAddress))); err != nil {
		return nil, fmt.Errorf("failed to encode requesterAddress length: %w", err)
	}
	buf.WriteString(rc.RequesterAddress)

	// Encode the extraArgs
	//buf.WriteByte(rc.ExtraArgs)

	return buf.Bytes(), nil
}
func GenerateVRFProof(suite kyber.Group, privateKey kyber.Scalar, data []byte, nonce int64) ([]byte, []byte, error) {
	// Convert nonce to a deterministic scalar
	nonceBytes := big.NewInt(nonce).Bytes()
	nonceScalar := suite.Scalar().SetBytes(nonceBytes)

	// Generate proof like in a Schnorr signature: R = g^k, s = k + e*x
	R := suite.Point().Mul(nonceScalar, nil) // R = g^k
	hash := sha256.New()
	rBytes, _ := R.MarshalBinary()
	hash.Write(rBytes)
	hash.Write(data)
	e := suite.Scalar().SetBytes(hash.Sum(nil))                             // e = H(R||data)
	s := suite.Scalar().Add(nonceScalar, suite.Scalar().Mul(e, privateKey)) // s = k + e*x

	// The VRF output (pseudo-random value) is hash of R combined with data
	vrfHash := sha256.New()
	vrfHash.Write(rBytes)         // Incorporate R
	vrfHash.Write(data)           // Incorporate input data
	vrfOutput := vrfHash.Sum(nil) // This is the deterministic "random" output

	// Serialize R and s into the proof
	sBytes, _ := s.MarshalBinary()
	proof := append(rBytes, sBytes...)

	return proof, vrfOutput, nil
}
