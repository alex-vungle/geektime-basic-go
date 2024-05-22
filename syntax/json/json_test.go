package json

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestJSON(t *testing.T) {
	u := User{
		Name:     "Tom",
		age:      18,
		Password: "my-password",
		Phone:    "15223456789",
		PhoneV2:  "15223456789",
	}
	data, err := json.Marshal(u)
	require.NoError(t, err)
	t.Log(string(data))

	// 假如说你在日志找到这个 1d8839703df511bf191df2f318c07b51c4266283eb8827bc535a4b914ee8a7a254bb394b0269af
	// 然后你要反推出来 phone number 是哪个，然后去数据查数据
	data, err = hex.DecodeString("1d8839703df511bf191df2f318c07b51c4266283eb8827bc535a4b914ee8a7a254bb394b0269af")
	require.NoError(t, err)
	val, err := aesDecrypt(data)
	require.NoError(t, err)
	t.Log("我手工解析出来手机号码", string(val))

	// 假如说你要解析这个日志文件，或者这是前端传过来的
	// 你希望把 phoneV2 还原回来
	input := `{"name":"Tom","phone":"152****6789","phoneV2":"03db922a73527bdf1043a66e46244367f3490874d3eb8661526e14db113d468b7538565189054b"}`
	var u1 User
	err = json.Unmarshal([]byte(input), &u1)
	require.NoError(t, err)
	t.Log("u1", u1)
}

type User struct {
	// 使用 json 这个标签来指定序列化之后的 JSON 里面的字段的名字
	Name string `json:"name"`
	// 我希望打印成 135****1234
	// 就是中间这四位被遮住了
	Phone   PhoneNumber   `json:"phone"`
	PhoneV2 PhoneNumberV2 `json:"phoneV2"`

	// 我希望把某个字段加密打印出来，加密之后，我能够解密得到

	// 我能不能忽略掉这个字段？
	Password string `json:"-"`
	// 私有字段会被忽略掉
	age int `json:"age"`
}

type PhoneNumberV2 string

const key = "0123456789abcdef"

// PhoneNumberV2 是一个支持加密序列化，并且支持解密反序列化的类型
func (p PhoneNumberV2) MarshalJSON() ([]byte, error) {
	data, err := p.aesEncrypt([]byte(p))
	if err != nil {
		return nil, err
	}
	data = []byte(hex.EncodeToString(data))
	data = append([]byte{'"'}, data...)
	data = append(data, '"')
	return data, err
}

func (p *PhoneNumberV2) UnmarshalJSON(input []byte) error {
	// 去掉双引号
	input = input[1 : len(input)-1]
	input, err := hex.DecodeString(string(input))
	if err != nil {
		return err
	}
	val, err := aesDecrypt(input)
	if err != nil {
		return err
	}
	*p = PhoneNumberV2(val)
	return nil
}

// 模拟加密
func (p PhoneNumberV2) aesEncrypt(data []byte) ([]byte, error) {
	newCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

func aesDecrypt(data []byte) ([]byte, error) {
	newCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	nonce, cipherData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, cipherData, nil)
}

type PhoneNumber string

func (p PhoneNumber) MarshalJSON() ([]byte, error) {
	// 返回一个 "123****5678" 这样的东西
	res := make([]byte, 0, len(p)+2)
	res = append(res, '"')
	res = append(res, p[:3]...)
	res = append(res, '*')
	res = append(res, '*')
	res = append(res, '*')
	res = append(res, '*')
	res = append(res, p[7:]...)
	res = append(res, '"')
	return res, nil
}
