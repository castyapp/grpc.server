package subtitle

import (
	"bytes"
	"fmt"
	"github.com/CastyLab/grpc.server/config"
	"github.com/CastyLab/grpc.server/services"
	"github.com/asticode/go-astisub"
	"io/ioutil"
	"mime/multipart"
	"os"
)

// Convert and return subtitle files to WebVTT
func ConvertToVTT(file multipart.File) (buffer *bytes.Buffer, err error) {

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	vttSubtitle, err := astisub.ReadFromSRT(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	buffer = new(bytes.Buffer)
	if err := vttSubtitle.WriteToWebVTT(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func Save(sFile *multipart.FileHeader) (string, error) {

	subtitleName := services.RandomNumber(20)

	subtitle, err := sFile.Open()
	if err != nil {
		return "", err
	}

	buf, err := ConvertToVTT(subtitle)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf("%s/uploads/subtitles/%s.vtt", config.Map.StoragePath, subtitleName))
	if err != nil {
		return "", err
	}

	if _, err := file.Write(buf.Bytes()); err != nil {
		return "", err
	}

	return subtitleName, nil
}