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
	plugin "github.com/hitzhangjie/protoc-gen-gorpc/plugin"
)

// GenerateTplFiles run go template engine to process template files
func (g *Generator) GenerateTplFiles() error {

	if len(g.allFiles) != 1 {
		return fmt.Errorf("only 1 *.proto is support, you specify %d", len(g.allFiles))
	}

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

	// run go template to generate template
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	root := filepath.Join(home, ".gorpc2/go")

	output, err := g.GetOutputDirectory()
	if err != nil {
		return err
	}
	os.MkdirAll(output, os.ModePerm)

	fn := func(path string, info os.FileInfo, err error) error {
		// 检查要不要处理当前文件
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 新生成文件目录结构，与模板路径保持一样的结构
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// 如果是文件，且为go模板文件，执行go模板引擎生成新文件
		// - 模板文件，执行模板处理引擎
		//target = filepath.Join(output, strings.TrimSuffix(rel, ".tpl"))
		target := strings.TrimSuffix(rel, ".tpl")
		return g.procTplFile(path, target, nfd)
	}

	return filepath.Walk(root, fn)
}

func (g *Generator) procTplFile(inFile, outFile string, nfd *gorpc.FileDescriptor) error {

	baseName := filepath.Base(inFile)

	instance, err := template.New(baseName).Funcs(gorpc.FuncMap).ParseFiles(inFile)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	p := TemplateParams{nfd, "whisper", "unspecified", 0}
	if err = instance.Execute(&buf, p); err != nil {
		return err
	}

	g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(outFile),
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
