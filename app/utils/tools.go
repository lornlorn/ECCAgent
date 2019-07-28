package utils

import (
    "bufio"
    "bytes"
    "crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    mrand "math/rand"
    "os"
    "strings"
    "time"

    "github.com/cihub/seelog"
    "github.com/tidwall/gjson"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
)

/*
ReadRequestBody2JSON func(reqBody io.ReadCloser) []byte
*/
func ReadRequestBody2JSON(reqBody io.ReadCloser) []byte {

    body, err := ioutil.ReadAll(reqBody)
    if err != nil {
        seelog.Errorf("ioutil.ReadAll Error : %v", err)
        return []byte{}
    }

    return body

}

/*
GetJSONResultFromRequestBody func(reqBody []byte, path string) gjson.Result
*/
func GetJSONResultFromRequestBody(reqBody []byte, path string) gjson.Result {
    return gjson.Get(string(reqBody), path)
}

/*
ReadJSONData2Array func(reqBody []byte, path string) []gjson.Result
*/
func ReadJSONData2Array(reqBody []byte, path string) []gjson.Result {
    j := gjson.Get(string(reqBody), path)
    return j.Array()
}

/*
Convert2JSON 任意数据类型转JSON
*/
func Convert2JSON(data interface{}) []byte {

    switch data.(type) {
    case []byte:
        retdata := data.([]byte)
        return retdata
    default:
        // log.Println("Convert To JSON args not []byte")
        retdata, err := json.Marshal(data)
        if err != nil {
            seelog.Errorf("json.Marshal Error : %v", err)
            return []byte("")
        }
        return retdata
    }

}

/*
读文件 并 设置偏移量和行数
*/

// ReadLines reads contents from file and splits them by new line.
// A convenience wrapper to ReadLinesOffsetN(filename, 0, -1).
func ReadLines(filename string) ([]string, error) {
    return ReadLinesOffsetN(filename, 0, -1)
}

// ReadLinesOffsetN reads contents from file and splits them by new line.
// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
//   n >= 0: at most n lines
//   n < 0: whole file
func ReadLinesOffsetN(filename string, offset uint, n int) ([]string, error) {

    f, err := os.Open(filename)
    if err != nil {
        return []string{""}, err
    }
    defer f.Close()

    var ret []string

    r := bufio.NewReader(f)
    for i := 0; i < n+int(offset) || n < 0; i++ {
        line, err := r.ReadString('\n')
        if err != nil {
            break
        }
        if i < int(offset) {
            continue
        }
        ret = append(ret, strings.Trim(line, "\n"))
    }

    return ret, nil

}

/*
生成随机UID GetUniqueID()
*/

// GetMd5String 生成32位MD5字符串
func GetMd5String(s string) string {
    newmd5 := md5.New()
    newmd5.Write([]byte(s))
    return hex.EncodeToString(newmd5.Sum(nil))
}

// GetUniqueID 生成UID唯一标识
func GetUniqueID() string {

    newbyte := make([]byte, 48)

    _, err := io.ReadFull(rand.Reader, newbyte)
    if err != nil {
        // seelog.Errorf("io.ReadFull Error : %v", err)
        return "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    }

    return GetMd5String(base64.URLEncoding.EncodeToString(newbyte))

}

// DecodeGBK2UTF8 GBK转UTF8
func DecodeGBK2UTF8(in []byte) ([]byte, error) {
    I := bytes.NewReader(in)
    O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(O)
    if e != nil {
        return nil, e
    }
    return d, nil
}

// RandString 生成随机字符串
func RandString(len int) string {
    r := mrand.New(mrand.NewSource(time.Now().Unix()))
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        b := r.Intn(26) + 65
        bytes[i] = byte(b)
    }
    return string(bytes)
}

// RandomString 生成随机字符串 不使用time函数
func RandomString(len int) string {
    var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
    var container string

    b := bytes.NewBufferString(str)
    length := b.Len()
    bigInt := big.NewInt(int64(length))
    for i:=0;i<len;i++{
        randomInt,_ := rand.Int(rand.Reader,bigInt)
        container += string(str[randomInt.Int64()])
    }
    return container
}

// RandomNumber 生成随机数字字符串 不使用time函数
func RandomNumber(len int) string {
    var numbers = []byte{1,2,3,4,5,6,7,8,9,0}
    var container string

    length := bytes.NewReader(numbers).Len()
    for i:=1;i<=len;i++{
        random,_:=rand.Int(rand.Reader,big.NewInt(int64(length)))
        container += fmt.Sprintf("%d",numbers[random.Int64()])
    }
    return container
}