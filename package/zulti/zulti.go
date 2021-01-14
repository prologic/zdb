package zulti

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// AddCookie :
func AddCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: name, Value: value, Expires: expire}
	http.SetCookie(w, &cookie)
}

// Hash :
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// FileToFolder :
func FileToFolder(filepath string) (string, string) {
	tempArr := strings.Split(filepath, "/")
	tempIndex := (len(tempArr) - 1)
	return strings.Join(tempArr[:tempIndex], "/"), strings.Join(tempArr[tempIndex:], "/")
}

// AddProtocol :
func AddProtocol(w http.ResponseWriter, r *http.Request, host string) string {
	_host := r.Host
	if host != "" {
		_host = host
	}
	if r.TLS == nil {
		_host = "http://" + _host
	} else {
		_host = "https://" + _host
	}
	return _host
}

// WriteLines :
func WriteLines(file string, lines []string) string {
	f, err := os.Create(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return ""
		}
	}
	return ""
}

var rxtimezone *time.Location

// Map : custom type
type Map map[string]interface{}

func init() {
}

// IsRun :
func IsRun(isRun bool, f func()) {
	if isRun {
		f()
	}
}

// IsRunViaParams :
func IsRunViaParams(isRun bool, f func(Map), params Map) (result Map) {
	if isRun {
		f(params)
	}
	return
}

// IsNil :
func IsNil(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return ""
	}
	return value
}

// Wrap :
func Wrap(vs ...interface{}) interface{} {
	return vs[0]
}

// Get :
func Get(data map[string]interface{}, params ...interface{}) (returnData interface{}) {
	var indexName interface{} = "default"
	var defaultValue interface{} = ""

	if len(params) >= 1 {
		indexName = params[0]
	}

	if len(params) >= 2 {
		defaultValue = params[1]
	}

	if val, ok := data[indexName.(string)]; ok {
		returnData = val
	} else {
		returnData = defaultValue
	}

	return
}

// GetMapMix :
func GetMapMix(_data map[string][]string, _param string, _default []string) (returnData []string) {
	indexName := _param
	defaultValue := []string{"", ""}

	if len(_default) > 0 {
		defaultValue = _default
	}

	if val, ok := _data[indexName]; ok {
		returnData = val
	} else {
		returnData = defaultValue
	}

	return
}

// GetDate
func GetDateStr() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

// IsIn :
func IsIn(a interface{}, list []interface{}) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// Find :
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if strings.HasPrefix(val, item) {
			return i, true
		}
	}
	return -1, false
}

// IfElse :
func IfElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}

	return b
}

// Timestamp :
func Timestamp() int64 {
	rxtimezone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(rxtimezone).UnixNano()
}

// TimestampMicrosec :
func TimestampMicrosec() int64 {
	rxtimezone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(rxtimezone).UnixNano() / int64(1000)
}

// TimestampMilisec :
func TimestampMilisec() int64 {
	rxtimezone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(rxtimezone).UnixNano() / int64(time.Millisecond)
}

// TimestampSec :
func TimestampSec() int64 {
	rxtimezone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(rxtimezone).Unix()
}

// CleanSERVER :
func CleanSERVER() string {
	return ""
}

// Call :
func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	if m[name] != nil {
		f := reflect.ValueOf(m[name])
		if len(params) != f.Type().NumIn() {
			err = errors.New("The number of params is not adapted")
			return
		}
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		result = f.Call(in)
	}
	return
}

// ReflectShow :
func ReflectShow(f interface{}) {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}

/////////////
// F I L E //
/////////////

// FileEnsureExists :
func FileEnsureExists(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

// FileExists :
func FileExists(filename string) bool {
	result := true
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return result
}

// FileWrite :
func FileWrite(filename string, data string) (err error) {
	err = FileEnsureExists(filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

// CheckErr L
func CheckErr(err error, errStr string) {
	if err != nil {
		fmt.Println(errStr)
	}
}

// JSONDecode :
func JSONDecode(jsonstr string) []interface{} {
	returnData := []interface{}{}
	uj := []interface{}{}

	if err := json.Unmarshal([]byte(jsonstr), &uj); err == nil {
		returnData = uj
	} else {
		fmt.Println(err)
	}

	return returnData
}

// JSONDecodeObj :
func JSONDecodeObj(jsonstr string) map[string]interface{} {
	returnData := map[string]interface{}{}
	uj := map[string]interface{}{}

	if err := json.Unmarshal([]byte(jsonstr), &uj); err == nil {
		returnData = uj
	} else {
		fmt.Println(err)
	}

	return returnData
}

// JSON :
func JSON(data Map) string {
	returnData := "{}"
	if uj, err := json.Marshal(data); err == nil {
		returnData = string(uj)
	}

	return returnData
}

// ToJSON :
func (om Map) ToJSON(order ...string) string {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	l := len(order)
	for i, k := range order {
		if om[k] != nil {

			switch v := om[k].(type) {
			case int:
				fmt.Fprintf(buf, "\"%s\":%v", k, v)
			case string:
				fmt.Fprintf(buf, "\"%s\":\"%v\"", k, v)
			default:
				fmt.Fprintf(buf, "\"%s\":\"%v\"", k, v)
			}

			if i < l-1 {
				buf.WriteByte(',')
			}
		}
	}
	buf.Write([]byte{'}'})

	// replace ,}
	return strings.Replace(buf.String(), ",}", "}", -1)
}

// ToLower :
func ToLower(value string) string {
	return strings.ToLower(value)
}

// ToString :
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", IsNil(value, ""))
}

// ToInt :
func ToInt(value interface{}) int {
	returnValue := 0

	switch value.(type) {
	case int:
		returnValue = value.(int)

	case float32:
		returnValue = int(value.(float32))

	case float64:
		returnValue = int(value.(float64))

	default:
		returnValue, _ = strconv.Atoi(fmt.Sprintf("%v", value))
	}

	return returnValue
}

// ToIntX :
func ToIntX(value interface{}) int64 {
	returnValue, _ := strconv.ParseInt(fmt.Sprintf("%v", IsNil(value, 0)), 10, 32)
	return returnValue
}

// ToFloatX :
func ToFloatX(value interface{}) float64 {
	returnValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", IsNil(value, "")), 10)
	return returnValue
}

// SendToClickhouse :
func SendToClickhouse(items []map[string]interface{}) {
	if len(items) > 0 {
	}
}

// DumpMap :
func DumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			DumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}
