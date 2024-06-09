package station

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type PodsRes struct {
	Witness        string `json:"witness"`          // `protobuf:"bytes,1,opt,name=witness,proto3" json:"witness,omitempty"`
	MerkleRootHash string `json:"merkle_root_hash"` // `protobuf:"bytes,2,opt,name=merkle_root_hash,json=merkleRootHash,proto3" json:"merkle_root_hash,omitempty"`
	Proof          string `json:"proof"`            // `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	PodNumber      string `json:"pod_number"`       // uint64 `protobuf:"varint,4,opt,name=pod_number,json=podNumber,proto3" json:"pod_number,omitempty"`
}

// GetPodResponse represents the structure of the JSON response from the API
type GetPodResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  PodsRes `json:"result"`
}

type PodDetails struct {
	Witness        []byte
	MerkleRootHash string
	Proof          []byte
}

func QueryPodNumber(podNumber uint64) (*PodDetails, error) {

	// API URL
	url := fmt.Sprintf("http://localhost:26667/tracks_get_pod?podNumber=%d", podNumber)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal response into Response struct
	var response GetPodResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	WitnessByte, err := StrToByte(response.Result.Witness)
	if err != nil {
		return nil, fmt.Errorf("failed to convert witness to byte: %w", err)
	}

	ProofByte, err := StrToByte(response.Result.Proof)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proof to byte: %w", err)
	}

	MerkleRootByte, err := StrToByte(response.Result.MerkleRootHash)
	if err != nil {
		return nil, fmt.Errorf("failed to convert merkle root hash to byte: %w", err)
	}

	MerkleRootHash := string(MerkleRootByte)

	pod := &PodDetails{
		Witness:        WitnessByte,
		MerkleRootHash: MerkleRootHash,
		Proof:          ProofByte,
	}

	return pod, nil
}

// StrToByte converts a string representation of a byte array to an actual byte array
func StrToByte(str string) ([]byte, error) {
	// Remove the surrounding brackets and spaces
	str = strings.Trim(str, "[]")
	str = strings.ReplaceAll(str, " ", "")

	// Split the string by commas
	strArray := strings.Split(str, ",")

	// Create a byte slice with the same length as the string array
	byteArray := make([]byte, len(strArray))

	// Convert each string to a byte
	for i, s := range strArray {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		byteArray[i] = byte(num)
	}

	return byteArray, nil
}
