package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

// IsEthPrivateKey validate if ti's a valid privateKey
func IsEthPrivateKey(hexString string) bool {
	privateKeyBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return false
	}

	curve := elliptic.P256()
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.Curve = curve
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyBytes)

	return privateKey != nil
}

// PublicKeyBytesToAddress ...
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}

// SigRSV signatures R S V returned as arrays
func SigRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := sigstr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}

// GetAddAndKey 从字符串中匹配私钥和地址
func GetAddAndKey(line string) (address, key string, err error) {
	re := regexp.MustCompile(`(?:0x)?[0-9a-fA-F]+`)

	// 以太坊地址正则表达式
	//addressRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	//// 以太坊私钥正则表达式
	//privateKeyRegex := regexp.MustCompile("^[0-9a-fA-F]{64}|^0x[0-9a-fA-F]{64}")
	// 测试地址验证
	//if re.MatchString(line) {
	//	//fmt.Println("地址格式正确", addressRegex.FindString(arr[0]))
	//	address = re.FindString(line)
	//} else {
	//	fmt.Println("地址格式错误")
	//	return "", "", errors.New("地址格式错误!!!")
	//}
	//if re.MatchString(line) {
	//	//fmt.Println("私钥格式正确", privateKeyRegex.FindString(arr[1]))
	//	key = re.FindString(line)
	//} else {
	//	fmt.Println("私钥格式错误")
	//	return "", "", errors.New("私钥格式错误!!!")
	//}

	matches := re.FindAllString(line, -1)
	for _, match := range matches {
		if len(match) == 42 {
			address = match
		} else if len(match) == 66 || len(match) == 64 {
			key = match
		} else {
			return "", "", errors.New("格式有误!!!")
		}
	}

	return address, key, nil
}

// GetAddressFromStr 从字符串中匹配地址
func GetAddressFromStr(line string) (address string, err error) {
	//以太坊地址正则表达式
	addressRegex := regexp.MustCompile("0x[0-9a-fA-F]+")

	if addressRegex.MatchString(line) {
		subMatch := addressRegex.FindAllString(line, -1)
		//fmt.Println(subMatch)
		for _, match := range subMatch {
			if len(match) == 42 {
				address = match
				break
			}
		}
		if address == "" {
			return "", errors.New("地址格式错误")
		}
		return address, nil
	} else {
		return "", errors.New("地址格式错误")
	}
}

// GetPrivateKeyFromStr 从字符串中匹配地址
func GetPrivateKeyFromStr(line string) (address string, err error) {

	// 以太坊私钥正则表达式
	privateKeyRegex := regexp.MustCompile(`[0-9a-fA-F]{64}`)
	// 测试地址验证
	if privateKeyRegex.MatchString(line) {
		//fmt.Println("私钥格式正确", privateKeyRegex.FindString(arr[1]))
		privateKey := privateKeyRegex.FindString(line)
		return privateKey, nil
	} else {
		return "", errors.New("私钥格式错误!!!")
	}

}

// FormatAddress 显示缩略地址
func FormatAddress(address string) string {
	if len(address) < 10 {
		return address
	}

	head := address[:8]
	tail := address[len(address)-5:]
	formattedAddress := fmt.Sprintf("%s...%s", head, tail)

	return formattedAddress
}
