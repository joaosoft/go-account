package gomoney

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	goerror "github.com/joaosoft/go-error/app"
	"github.com/labstack/echo"
)

type downloadData struct {
	FileName string
	Data     bytes.Buffer
}

func download(key string, ctx echo.Context) ([]*downloadData, *goerror.ErrorData) {
	downloads := make([]*downloadData, 0)

	// multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	files := form.File[key]

	for _, file := range files {
		// source
		src, err := file.Open()
		if err != nil {
			return nil, goerror.NewError(err)
		}
		defer src.Close()

		// destination
		var dst bytes.Buffer

		// copy
		dst.ReadFrom(src)

		downloads = append(downloads,
			&downloadData{
				FileName: file.Filename,
				Data:     dst,
			})
	}
	log.Infof("uploaded successfully %d %s", len(files), key)

	return downloads, nil
}

func upload(client *http.Client, url string, values map[string]io.Reader) *goerror.ErrorData {
	var b bytes.Buffer
	var err error
	downloads := multipart.NewWriter(&b)

	for key, value := range values {
		var fw io.Writer
		if x, ok := value.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if file, ok := value.(*os.File); ok {
			if fw, err = downloads.CreateFormFile(key, file.Name()); err != nil {
				return goerror.NewError(err)
			}
		} else {
			// Add other fields
			if fw, err = downloads.CreateFormField(key); err != nil {
				return goerror.NewError(err)
			}
		}
		if _, err = io.Copy(fw, value); err != nil {
			return goerror.NewError(err)
		}

	}
	// close the multipart
	downloads.Close()

	// now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return goerror.NewError(err)
	}
	// don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", downloads.FormDataContentType())

	// submit the request
	res, err := client.Do(req)
	if err != nil {
		return goerror.NewError(err)
	}

	// check the response
	if res.StatusCode != http.StatusOK {
		return goerror.FromString(fmt.Sprintf("bad status: %s", res.Status))
	}

	return nil
}
