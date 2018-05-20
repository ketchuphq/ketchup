package admin

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/content/content"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/util/errors"
)

const previewTemplate = `<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.min.css" />
	</head>
	<body>
		<div class="markdown-body" style='padding: 15px;'>
			%s
		</div>
	</body>
</html>`

// handlePreview renders the given content object.
func (m *Module) handlePreview(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	buf := bytes.NewBufferString(req.PostFormValue("content"))

	req2 := &api.PreviewContentRequest{}
	err := jsonpb.Unmarshal(buf, req2)
	if err != nil {
		fmt.Println(err)
		return err
	}
	data, err := content.RenderContent(req2.GetContent())
	if err != nil {
		return errors.Wrap(err)
	}
	_, err = fmt.Fprintf(rw, previewTemplate, data)
	return errors.Wrap(err)
}
