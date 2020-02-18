package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Round 四舍五入
func Round(val float64, places int) float64 {
	f := math.Pow10(places)
	return float64(int64(val*f+0.5)) / f
}

func FloatMultiToInt64(args ...float64) int64 {
	var r float64 = 1
	for _, v := range args {
		r *= v
	}
	r = Round(r, 0)
	return int64(r)
}

func FloatDivisionToInt64(args ...float64) int64 {
	if len(args) == 0 {
		return 0
	}
	r := args[0]
	for _, v := range args[1:] {
		r /= v
	}
	r = Round(r, 0)
	return int64(r)
}

// LRead 向上读取文件
func LRead(name string, level int) (raw []byte, err error) {
	var file *os.File
	for i := 0; i <= level; i++ {
		filePath := fmt.Sprintf("%s%s", strings.Repeat("../", i), name)
		file, err = os.OpenFile(filePath, os.O_RDONLY, 0600)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return
	}
	raw, err = ioutil.ReadAll(file)
	return
}

func TimeLimit(start, end int) bool {
	h := time.Now().Hour()
	if h >= start && h < end {
		return true
	}
	return false
}

func RandomString(n int) string {
	var original = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	random := make([]rune, n)
	for i := range random {
		random[i] = original[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(62)]
	}
	return string(random)
}

type Attachment struct {
	Name string
	Body []byte
}

func SendToMail(user, password, name, addr, to, subject, body string, isHtml bool, attachments ...Attachment) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}
	// create new SMTP client
	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", user, password, host)
	err = smtpClient.Auth(auth)
	if err != nil {
		return err
	}
	from := mail.Address{Name: name, Address: user}
	if err := smtpClient.Mail(from.Address); err != nil {
		return err
	}
	for _, v := range strings.Split(to, ";") {
		if err := smtpClient.Rcpt(strings.TrimSpace(v)); err != nil {
			return err
		}
	}

	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}
	var contentType string
	if isHtml {
		contentType = "text/html;\r\n\tcharset=utf-8"
	} else {
		contentType = "text/plain;\r\n\tcharset=utf-8"
	}

	boundary := "----THIS_IS_BOUNDARY_JUST_MAKE_YOURS_MIXED"

	buffer := bytes.NewBuffer(nil)

	header := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: multipart/mixed;\r\n\tBoundary=\"%s\"\r\n"+
		"Mime-Version: 1.0\r\n"+
		"Date: %s\r\n\r\n", to, user, subject, boundary, time.Now().String())
	buffer.WriteString(header)
	buffer.WriteString("This is a multi-part message in MIME format.\r\n\r\n")

	// 正文
	if len(body) > 0 {
		bodyBoundary := "----THIS_IS_BOUNDARY_JUST_MAKE_YOURS_BODY"
		buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buffer.WriteString(fmt.Sprintf("Content-Type: multipart/alternative;\r\n\tBoundary=\"%s\"\r\n\r\n", bodyBoundary))

		buffer.WriteString(fmt.Sprintf("--%s\r\n", bodyBoundary))
		buffer.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
		buffer.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n"))
		buffer.WriteString(fmt.Sprintf("%s\r\n\r\n", base64.StdEncoding.EncodeToString([]byte(body))))
		buffer.WriteString(fmt.Sprintf("--%s--\r\n", bodyBoundary))

	}
	for _, attachment := range attachments {
		t := mime.TypeByExtension(filepath.Ext(attachment.Name))
		if t == "" {
			t = "application/octet-stream"
		}
		buffer.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		buffer.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n"))
		buffer.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n\r\n", t, attachment.Name))
		buffer.WriteString(fmt.Sprintf("%s\r\n\r\n", base64.StdEncoding.EncodeToString(attachment.Body)))
	}

	buffer.WriteString("\r\n\r\n--" + boundary + "--")
	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return smtpClient.Quit()
}

func JSON(data interface{}) {
	bts, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(bts))
}
func Tracks() []string {
	var list []string
	var i int
	for {
		if i > 14 {
			break
		}
		_, file, line, ok := runtime.Caller(i)
		// dir, filename := filepath.Split(file)
		// file = fmt.Sprintf("%s/%s", filepath.Base(dir), filename)
		if !ok {
			break
		}
		i++
		list = append(list, fmt.Sprintf("%s:%d", file, line))
	}
	return list
}

func Gzip(data []byte) ([]byte, error) {
	var res bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&res, 7)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	} else {
		gz.Close()
	}
	return res.Bytes(), nil
}

func GetImageSrc(domain, fieldID string) string {
	if strings.HasPrefix(fieldID, "data:image") {
		return fieldID
	}
	domain = strings.TrimRight(domain, "/") + "/upload/images/"
	return domain + fieldID
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
