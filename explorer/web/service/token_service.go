package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"encoding/base64"
	"bytes"
	"time"
	"os"
	"image"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"image/gif"
	"io"
	"errors"
	"github.com/wlcy/tron/explorer/web/buffer"
)

// QueryCommonTokenListBuffer
func QueryCommonTokenListBuffer() ([]*entity.TokenInfo, error) {
	tokenBuffer := buffer.GetTokenBuffer()
	commonTokenList := tokenBuffer.GetCommonTokenList()
	return commonTokenList, nil

}

// QueryIcoTokenListBuffer
func QueryIcoTokenListBuffer() ([]*entity.TokenInfo, error) {
	tokenBuffer := buffer.GetTokenBuffer()
	icoTokenList := tokenBuffer.GetIcoTokenList()
	return icoTokenList, nil

}

// QueryTokenDetailListBuffer
func QueryTokenDetailListBuffer() ([]*entity.TokenInfo, error) {
	tokenBuffer := buffer.GetTokenBuffer()
	tokenDetailList := tokenBuffer.GetTokensDetailList()
	return tokenDetailList, nil

}

//UploadTokenLogo 保存图片
func UploadTokenLogo(defaultPath, imgURL, imageData, address string) (string, error) {
	imgData, err := base64.StdEncoding.DecodeString(imageData)
	if nil != err {
		return "", err
	}
	buffer := bytes.NewBuffer(imgData)
	tempFileName := fmt.Sprintf("tokenLogo_%v.jpeg", time.Now().Format("20060102150405.000000"))

	dist, err := os.Create(defaultPath + "/" + tempFileName)
	if err != nil {
		log.Error(err)
	}
	defer dist.Close()

	err = scale(buffer, dist, 0, 0, 0)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("save file %v ", dist)
	dst := fmt.Sprintf("%v/%v", imgURL, tempFileName)
	err = InsertOrUpdateLogo(address, dst)

	return dst, err
}


/*
* 缩略图生成
* 入参:
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度为0则精度保持不变
* 矩形坐标系起点是左上
* 返回:error
 */
func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		log.Error(err)
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Resize(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	log.Debugf("fm:%v", fm)
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
		/*case "bmp":
		return bmp.Encode(out, canvas)*/
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}


//InsertOrUpdateLogo 更新或插入logo
func InsertOrUpdateLogo(address, url string) error {
	if address == "" || url == "" {
		log.Error("address or url is nil")
		return util.NewErrorMsg(util.Error_common_internal_error)
	}

	addressInfo, err := module.IsAddressExist(address)
	if err != nil {
		log.Errorf("check address:[%v] isExist err:[%v]", address, err)
		return err
	}
	log.Errorf("addressInfo:[%#v] ", addressInfo)
	if addressInfo {
		err = module.UpdateLogoInfo(address, url)
	} else {
		err = module.InsertLogoInfo(address, url)
	}
	return err
}

// SyncAssetIssueParticipated
func SyncAssetIssueParticipated() {
	assetIssues, _ := module.QueryAllAssetIssue()
	if len(assetIssues) == 0 {
		log.Info("SyncAssetIssueParticipated len(assetIssues) == 0")
		return
	}
	for index := range assetIssues {
		assetIssue := assetIssues[index]
		participateAsset, _ := module.QueryParticipateAsset(assetIssue.OwnerAddress, assetIssue.AssetName)
		if participateAsset.AssetName != "" && participateAsset.TotalAmount > assetIssue.Participated {
			module.UpdateAssetIssue(assetIssue.OwnerAddress, assetIssue.AssetName, participateAsset.TotalAmount)
		}
	}

}

// QueryAssetBalances
func QueryAssetBalances(req *entity.Token) (*entity.AssetBalanceResp, error){
	assetBalanceResp, err := module.QueryAssetBalances(req)
	if err != nil {
		log.Errorf("QueryAssetBalances list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	return assetBalanceResp, nil
}
