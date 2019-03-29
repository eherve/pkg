package httptool

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

//Progress progress listener interface
type Progress interface {
	SetSize(s int64)
	SetWidth(w int64)
	Write(p []byte) (n int, err error)
}

//Download download file
func Download(client *http.Client, dest string, url string, progress Progress) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	destTmp := fmt.Sprintf("%s.tmp", dest)

	write := func() error {
		tmp, err := os.Create(destTmp)
		if err != nil {
			return err
		}
		defer tmp.Close()

		if progress != nil {
			size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
			if err != nil {
				return err
			}
			progress.SetSize(size)
			if _, err := io.Copy(tmp, io.TeeReader(resp.Body, progress)); err != nil {
				return err
			}
			return nil
		}
		if _, err := io.Copy(tmp, resp.Body); err != nil {
			return err
		}
		return nil
	}
	if err := write(); err != nil {
		return err
	}

	err = os.Rename(destTmp, dest)
	return err
}
