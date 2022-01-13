package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateFileMetadata struct {
	Name string `json:"name"`
}

type UpdateFileResponse struct {
	*FileMetadata
	Error *ErrorResponse `json:"error,omitempty"`
}

func updateFileParseRequest(ctx *gin.Context) (fileData, *APIError) {
	res := fileData{
		ID: ctx.Param("id"),
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return fileData{}, InternalServerError(fmt.Errorf("problem reading multipart form: %w", err))
	}

	file := form.File["file"]
	if len(file) != 1 {
		return fileData{}, ErrMultipartFileWrong
	}

	res.header = file[0]

	metadata, ok := form.Value["metadata"]
	if ok {
		if len(metadata) != len(file) {
			return fileData{}, ErrMetadataLength
		}

		d := updateFileMetadata{}
		if err := json.Unmarshal([]byte(metadata[0]), &d); err != nil {
			return fileData{}, WrongMetadataFormatError(err)
		}

		res.Name = d.Name
	} else {
		res.Name = res.header.Filename
	}

	return res, nil
}

func (ctrl *Controller) updateFile(ctx *gin.Context) (FileMetadata, *APIError) {
	file, apiErr := updateFileParseRequest(ctx)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	originalMetadata, apiErr := ctrl.metadataStorage.GetFileByID(ctx.Request.Context(), file.ID, ctx.Request.Header)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	if apiErr = checkFileSize(
		file.header,
		originalMetadata.Bucket.MinUploadFile,
		originalMetadata.Bucket.MaxUploadFile,
	); apiErr != nil {
		return FileMetadata{}, InternalServerError(fmt.Errorf("problem checking file size %s: %w", file.Name, apiErr))
	}

	if apiErr := ctrl.metadataStorage.SetIsUploaded(ctx, file.ID, false, ctx.Request.Header); apiErr != nil {
		return FileMetadata{}, apiErr.ExtendError(
			fmt.Sprintf("problem flagging file as pending upload %s: %s", file.Name, apiErr.Error()),
		)
	}

	filepath := originalMetadata.BucketID + "/" + file.ID
	etag, contentType, apiErr := ctrl.uploadSingleFile(file, filepath)
	if apiErr != nil {
		// let's revert the change to isUploaded
		_ = ctrl.metadataStorage.SetIsUploaded(ctx, file.ID, true, ctx.Request.Header)

		return FileMetadata{}, InternalServerError(fmt.Errorf("problem processing file %s: %w", file.Name, apiErr))
	}

	newMetadata, apiErr := ctrl.metadataStorage.PopulateMetadata(
		ctx,
		file.ID, file.Name, file.header.Size, originalMetadata.BucketID, etag, true, contentType,
		ctx.Request.Header,
	)
	if apiErr != nil {
		return FileMetadata{}, apiErr.ExtendError(fmt.Sprintf("problem populating file metadata for file %s", file.Name))
	}

	return newMetadata, nil
}

func (ctrl *Controller) UpdateFile(ctx *gin.Context) {
	metadata, apiErr := ctrl.updateFile(ctx)
	if apiErr != nil {
		_ = ctx.Error(fmt.Errorf("problem parsing request: %w", apiErr))

		ctx.Header("X-Error", apiErr.publicMessage)
		ctx.AbortWithStatus(apiErr.statusCode)

		ctx.JSON(apiErr.statusCode, UpdateFileResponse{nil, apiErr.PublicResponse()})

		return
	}

	ctx.JSON(http.StatusOK, UpdateFileResponse{&metadata, nil})
}