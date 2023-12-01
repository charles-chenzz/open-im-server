package api

import (
	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"math/rand"
	"net/http"
	"strconv"
)

type Third rpcclient.Third

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry) Third {
	return Third(*rpcclient.NewThird(discov))
}

func (o *Third) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *Third) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

// s3

func (o *Third) PartLimit(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartLimit, o.Client, c)
}

func (o *Third) PartSize(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartSize, o.Client, c)
}

func (o *Third) InitiateMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.InitiateMultipartUpload, o.Client, c)
}

func (o *Third) AuthSign(c *gin.Context) {
	a2r.Call(third.ThirdClient.AuthSign, o.Client, c)
}

func (o *Third) CompleteMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.CompleteMultipartUpload, o.Client, c)
}

func (o *Third) AccessURL(c *gin.Context) {
	a2r.Call(third.ThirdClient.AccessURL, o.Client, c)
}

func (o *Third) ObjectRedirect(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	if name[0] == '/' {
		name = name[1:]
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = strconv.Itoa(rand.Int())
	}
	ctx := mcontext.SetOperationID(c, operationID)
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) == 0 {
			continue
		}
		query[key] = values[0]
	}
	resp, err := o.Client.AccessURL(ctx, &third.AccessURLReq{Name: name, Query: query})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, resp.Url)
}

// logs

func (o *Third) UploadLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.UploadLogs, o.Client, c)
}

func (o *Third) DeleteLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.DeleteLogs, o.Client, c)
}

func (o *Third) SearchLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.SearchLogs, o.Client, c)
}

func GetPrometheus(c *gin.Context) {
	c.Redirect(http.StatusFound, config.Config.Prometheus.PrometheusUrl)
}
