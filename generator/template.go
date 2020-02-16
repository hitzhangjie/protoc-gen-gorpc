package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc"
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc/gotpl"
	plugin "github.com/hitzhangjie/protoc-gen-gorpc/plugin"
)

// GenerateTplFiles run go template engine to process template files
func (g *Generator) GenerateTplFiles() error {

	for _, file := range g.allFiles {
		if err := g.generateTplFile(file); err != nil {
			return err
		}
	}
	return nil
}

// GetOutputDirectory get the output directory of code generation for current *.proto
func (g *Generator) GetOutputDirectory() (string, error) {
	if len(g.allFiles) != 1 {
		return "", fmt.Errorf("only 1 *.proto is support, you specify %d", len(g.allFiles))
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	file := g.allFiles[0]
	baseName := filepath.Base(file.GetName())
	fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	output := filepath.Join(wd, fileName)

	return output, nil
}

// generateTplFile process the go template files
func (g *Generator) generateTplFile(file *FileDescriptor) error {

	if len(file.FileDescriptorProto.Service) == 0 {
		return errors.New("No RPC Service defined")
	}

	// 准备模板变量信息
	nfd, err := convertFileDescriptor(file)
	if err != nil {
		return err
	}

	// 处理各模板文件
	for fp, tpl := range gotpl.GoRPCTemplates {
		err := g.procTplFile(fp, tpl, nfd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) procTplFile(tplFileName, tplFileContent string, nfd *gorpc.FileDescriptor) error {

	instance, err := template.New(tplFileName).Funcs(gorpc.FuncMap).Parse(tplFileContent)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	p := TemplateParams{nfd, "whisper", "unspecified", 0}
	if err = instance.Execute(&buf, p); err != nil {
		return err
	}

	g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(tplFileName),
		Content: proto.String(buf.String()),
	})

	return nil
}

type TemplateParams struct {
	*gorpc.FileDescriptor
	Protocol     string
	GoMod        string
	ServiceIndex int
}
