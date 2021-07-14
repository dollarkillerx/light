package cryptology

import (
	"fmt"
	"log"
	"testing"
)

func TestAES(t *testing.T) {
	key := []byte("12345678912345s612345678912345s6")
	val := []byte("hello world dsads1456484848948adsadsadsadsads1456484848948adsadsadsadsads1456484848948adsadsadsa")
	fmt.Println(len(val))
	encrypt, err := AESEncrypt(key, val)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(len(encrypt))

	decrypt, err := AESDecrypt(key, encrypt)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(string(decrypt))
}

func TestRSA(t *testing.T) {
	privateKey, pubKey, _ := GenRsaKey()
	fmt.Println("key: ", string(privateKey))
	fmt.Println("pubKey: ", string(pubKey))

	// 签名
	data1 := "tESt rsa 1212"
	sha256, _ := RsaSignWithSha256([]byte(data1), privateKey)
	fmt.Println("sh: ", string(sha256))

	// 验签
	fmt.Println(RsaVerySignWithSha256([]byte(data1), sha256, pubKey))

	// 加密
	encrypt, _ := RsaEncrypt([]byte(data1), pubKey)
	fmt.Println("en: ", encrypt, "  ", len(encrypt), "  org: ", len([]byte(data1)))

	// 解密
	decrypt, _ := RsaDecrypt(encrypt, privateKey)
	fmt.Println(string(decrypt))
}
