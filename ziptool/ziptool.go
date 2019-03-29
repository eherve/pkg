package ziptool

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//Progress progress listener interface
type Progress interface {
	SetSize(s int64)
	Tick(n int) (int64, error)
}

//Unzip unzip archive
func Unzip(src, dest string, progress Progress) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, os.ModePerm)

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if progress != nil {
		progress.SetSize(int64(len(r.File)))
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if progress != nil {
			progress.Tick(1)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

//ZipFolder create zip archive from folder
func ZipFolder(filename, base, folder string, progress Progress) error {
	files, err := listFolderRec(folder, make([]string, 0))
	if err != nil {
		return err
	}
	return ZipFiles(filename, base, files, progress)
}

//ZipFiles create zip archive from files
func ZipFiles(filename, base string, files []string, progress Progress) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	if progress != nil {
		progress.SetSize(int64(len(files)))
	}
	for _, file := range files {
		if err = addFileToZip(zipWriter, base, file); err != nil {
			return err
		}
		if progress != nil {
			progress.Tick(1)
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, base, filename string) error {
	filename = strings.TrimLeft(filename, base)
	fileToZip, err := os.Open(filepath.Join(base, filename))
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func listFolderRec(base string, data []string) ([]string, error) {
	files, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			subfiles, err := listFolderRec(filepath.Join(base, f.Name()), make([]string, 0))
			if err != nil {
				return nil, err
			}
			data = append(data, subfiles...)
		} else {
			data = append(data, filepath.Join(base, f.Name()))
		}
	}

	return data, nil
}
