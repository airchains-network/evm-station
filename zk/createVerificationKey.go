package zk

import (
	"cosmossdk.io/log"
	"encoding/json"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"os"
)

// CreateVkPkNew generates and saves a new Proving Key and Verification Key if either file doesn't exist
func CreateVkPkNew(homeDir string) error {

	provingKeyFile := homeDir + "/provingKey.txt"
	verificationKeyFile := homeDir + "/verificationKey.json"

	_, err1 := os.Stat(provingKeyFile)
	_, err2 := os.Stat(verificationKeyFile)
	// If either file doesn't exist, generate and save new keys
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		provingKey, verificationKey, err := GenerateVerificationKey()
		if err != nil {
			return err
		}

		// Save Proving Key
		pkFile, err := os.Create(provingKeyFile)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Unable to create Proving Key file" + err.Error())
			return err
		}
		_, err = provingKey.WriteTo(pkFile)
		pkFile.Close()
		if err != nil {
			log.NewLogger(os.Stderr).Error("Unable to write Proving Key" + err.Error())
			return err
		}

		// Save Verification Key
		file, _ := json.MarshalIndent(verificationKey, "", " ")
		err = os.WriteFile(verificationKeyFile, file, 0644)
		if err != nil {
			log.NewLogger(os.Stderr).Error("Unable to write Verification Key to file" + err.Error())
			return err
		}
		log.NewLogger(os.Stdout).Info("Proving key and Verification key generated and saved successfully", "Proving Key", provingKeyFile, "Verification Key", verificationKeyFile)

		return nil
	} else {
		log.NewLogger(os.Stdout).Info("Both Proving key and Verification key already exist. No action needed.")
		return nil
	}
}

func GetVkPk(homeDir string) (groth16.ProvingKey, groth16.VerifyingKey, error) {

	provingKeyFile := homeDir + "/provingKey.txt"
	verificationKeyFile := homeDir + "/verificationKey.json"

	// Read Proving Key
	pk, err := ReadProvingKeyFromFile2(provingKeyFile)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to read Proving Key")
		return nil, nil, err
	}

	vk, err := ReadVerificationKeyFromFile(verificationKeyFile)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to read Verification Key")
		return nil, nil, err
	}

	return pk, vk, nil
}

func ReadProvingKeyFromFile2(filename string) (groth16.ProvingKey, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pk := groth16.NewProvingKey(ecc.BLS12_381)
	_, err = pk.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func ReadVerificationKeyFromFile(filename string) (groth16.VerifyingKey, error) {

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	vk := groth16.NewVerifyingKey(ecc.BLS12_381)
	err = json.Unmarshal(file, vk)
	if err != nil {
		return nil, err
	}

	return vk, nil

}
