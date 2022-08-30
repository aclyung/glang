package compiler

import (
	buitin "almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"almeng.com/glang/global"
	"fmt"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

var errh = func(err error) { println(err.Error()) }

type Compiler struct {
	Target Target
	Module *ir.Module
	Spaces general.List[*Space]
	Global *Space
}

func NewCompiler(t Target) (c Compiler) {
	mod := ir.NewModule()
	c = Compiler{t, mod, nil, nil}
	return
}

func (c *Compiler) InitGlobal() {
	c.Global = &Space{
		Name: &syntax.Name{Value: "#global"},
		Decl: general.EmptyList[syntax.Decl](),
	}
	buitin.InitTypes(c.Global.Decl)
	buitin.InitConsts(c.Global.Decl)
	buitin.InitModule(c.Module)

}

type Space struct {
	Name *syntax.Name
	Decl *general.List[syntax.Decl]
}

func Compile(filename string) {
	f, _ := os.Open(filename)
	// Node
	t := NewTarget(AARCH64, APPLE, DARWIN)
	c := NewCompiler(t)
	global.Init(c.Module)
	c.InitGlobal()
	if true {
		var src string
		f, _ := os.Open(filename)
		if b, err := io.ReadAll(f); err != nil {
			panic(err)
		} else {
			src = string(b)
		}
		syntax.TokenizingTest(filename, src)
	}
	file := syntax.Parse(filename, f, errh)
	// TODO: Node to llvm IR
	c.CompileFile(file)
	// TODO: link write file
	compiled := c.GetIR()

	println(compiled)
	tmpDir, err := ioutil.TempDir("", "glang")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(tmpDir+"/main.ll", []byte(compiled), 0644)
	if err != nil {
		panic(err)
	}

	out, err := os.UserHomeDir()
	if err != nil {
		return
	}

	clangArgs := []string{
		t.String(),
		"-Wno-override-module",
		tmpDir + "/main.ll",
		"-o", out + "/Desktop/exec", "-O3",
	}

	cmd := exec.Command("clang", clangArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(string(output))
		panic(err)
	}
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return
}

func (c *Compiler) GetIR() string {
	return c.Module.String()
}

func CodeGen(node *syntax.File) {
	s := &Space{Name: node.SpaceName, Decl: general.NewList(node.DeclList...)}
	m := ir.NewModule()
	decl := func(d syntax.Decl) {
		ParseDecl(s, m, d)
	}
	s.Decl.Each(decl)
}
