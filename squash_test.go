package squashfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	goappimage "github.com/CalebQ42/GoAppImage"
)

const (
	downloadURL  = "https://github.com/zilti/code-oss.AppImage/releases/download/continuous/Code_OSS-x86_64.AppImage"
	appImageName = "Code_OSS.AppImage"
	squashfsName = "testing.squashfs"
)

func TestMain(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	squashFil, err := os.Open(wd + "/testing/" + squashfsName)
	if os.IsNotExist(err) {
		TestCreateSquashFromAppImage(t)
		squashFil, err = os.Open(wd + "/testing/" + squashfsName)
		if err != nil {
			t.Fatal(err)
		}
	}
	defer squashFil.Close()
	stat, _ := squashFil.Stat()
	rdr, err := NewSquashfsReader(io.NewSectionReader(squashFil, 0, stat.Size()))
	if err != nil {
		t.Fatal(err)
	}
	//testing code to print out the directory structure
	// rdr.GetFileStructure()
	extractionFil := "Proton-5.9-GE-8-ST.tar.gz"
	os.Remove(wd + "/testing/" + extractionFil)
	desk, err := os.Create(wd + "/testing/" + extractionFil)
	if err != nil {
		t.Fatal(err)
	}
	ext, err := rdr.ReadFile(extractionFil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Size!", ext.FileSize)
	n, err := io.CopyBuffer(desk, ext, make([]byte, rdr.super.BlockSize/2))
	if err != nil {
		fmt.Println("Read", n)
		t.Fatal(err)
	}
	t.Fatal("No problems here!")
}

func TestCreateSquashFromAppImage(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(wd+"/testing", 0777)
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}
	_, err = os.Open(wd + "/testing/" + appImageName)
	if os.IsNotExist(err) {
		downloadTestAppImage(t, wd+"/testing")
		_, err = os.Open(wd + "/testing/" + appImageName)
		if err != nil {
			t.Fatal(err)
		}
	} else if err != nil {
		t.Fatal(err)
	}
	ai := goappimage.NewAppImage(wd + "/testing/" + appImageName)
	aiFil, err := os.Open(wd + "/testing/" + appImageName)
	if err != nil {
		t.Fatal(err)
	}
	defer aiFil.Close()
	aiFil.Seek(ai.Offset, 0)
	os.Remove(wd + "/testing/" + squashfsName)
	aiSquash, err := os.Create(wd + "/testing/" + squashfsName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(aiSquash, aiFil)
	if err != nil {
		t.Fatal(err)
	}
}

func downloadTestAppImage(t *testing.T, dir string) {
	//seems to time out. Need to fix that at some point
	appImage, err := os.Create(dir + "/" + appImageName)
	if err != nil {
		t.Fatal(err)
	}
	defer appImage.Close()
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := check.Get(downloadURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(appImage, resp.Body)
	if err != nil {
		t.Fatal(err)
	}
}
