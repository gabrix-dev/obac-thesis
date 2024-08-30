package repository

import (
	"encoding/json"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type IpfsRepository struct {
	gateway    string
	imageCache map[string]string
}

func NewIpfsRepositoy() (*IpfsRepository, error) {
	imageCache := make(map[string]string)
	imageCache["QmY4eRbmnpgpa7Hg9TyxtDvg7d9MtX2ctpLViPofDJRJ8c"] = "https://i.ibb.co/kGRfYRW/fifa-basic-ticket.png"
	imageCache["QmP9DrACijtQPfxihVwcgR3LVTF1myFrKgLs1YEprhzTbC"] = "https://i.ibb.co/5FTWNk6/fifa-basicplus-ticket.png"
	imageCache["QmPsf1okzqbJKoT7jjNMAEE4JXVqk3fhmYQtF19wgAXxVn"] = "https://i.ibb.co/jbbtGKP/fifa-vip-ticket.png"
	imageCache["QmPLpE6XkDn25kS95jjtrDuAE6vUb1AhzxANN9jWRefwDx"] = "https://i.ibb.co/vQkHrbJ/qatar-2022-cocacola.jpg"
	return &IpfsRepository{
		gateway:    "https://ipfs.io/ipfs/",
		imageCache: imageCache,
	}, nil
}

func (i IpfsRepository) GetFilesBatch(ipfsLinks []string) ([][]byte, error) {
	var files [][]byte
	for _, ipfsLink := range ipfsLinks {
		file, err := i.GetFile(ipfsLink)
		if err == nil {
			files = append(files, file)
		}
	}
	return files, nil
}

func (i IpfsRepository) GetFile(ipfsLink string) ([]byte, error) {
	ipfsGwUri := strings.Replace(ipfsLink, "ipfs://", i.gateway, 1)
	content, err := i.retreiveFile(ipfsGwUri)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func (i IpfsRepository) GetImage(ipfsLink string) (string, error) {
	fileBytes, err := i.GetFile(ipfsLink)
	if err != nil {
		return "", err
	}
	var openseaMetadata models.OpenSeaMetadata
	if err = json.Unmarshal(fileBytes, &openseaMetadata); err != nil {
		return "", err
	}
	ipfsCidAndId := strings.Split(openseaMetadata.Image, "//")
	ipfsCid := strings.Split(ipfsCidAndId[1], "/")[0]
	imageUrl, exist := i.imageCache[ipfsCid]
	if !exist {
		return "", errors.New("image not cached")
	}
	return imageUrl, nil
}

func (i IpfsRepository) retreiveFile(ipfsGwUri string) ([]byte, error) {
	metadataBasePath := "resources/metadata"
	part2 := strings.Split(ipfsGwUri, "//")
	id := strings.Split(part2[1], "/")

	metadataPath := metadataBasePath + "/" + id[3] + ".json"

	jsonFile, err := os.Open(metadataPath)
	if err != nil {
		return []byte{}, fmt.Errorf("error opening the metadata file")
	}
	defer jsonFile.Close()
	content, _ := ioutil.ReadAll(jsonFile)
	return content, nil
}
