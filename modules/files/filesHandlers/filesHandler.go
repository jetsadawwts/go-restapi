package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jetsadawwts/go-restapi/config"
	"github.com/jetsadawwts/go-restapi/modules/entities"
	"github.com/jetsadawwts/go-restapi/modules/files"
	"github.com/jetsadawwts/go-restapi/modules/files/filesUsecases"
	"github.com/jetsadawwts/go-restapi/pkg/utils"
)

type filesHandlersErrCode string

const (
	uploadErr filesHandlersErrCode = "files-001"
	deleteErr filesHandlersErrCode = "files-002"
)


type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFiles(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg: cfg,
		filesUsecase: filesUsecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	//Files ext validation
	extMap := map[string]string{
		"png": "png",
		"jpg": "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename),".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadGateway.Code,
				string(uploadErr),
				"extension is not acceptable.",
			).Res()
		}
		
		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadGateway.Code,
				string(uploadErr),
				fmt.Sprintf("file size must less than %d mib", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)

		req = append(req, &files.FileReq{
			File: file,
			Destination: destination + "/" + filename,
			FileName: filename,
			Extension: ext,
		})

	}

	res, err := h.filesUsecase.UploadToGCP(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}


func (h *filesHandler) DeleteFiles(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.filesUsecase.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}