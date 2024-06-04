/**
RellwNote 扩展要放到本软件或文档库中的 extensions 文件夹，扩展文件夹以扩展自身的名字命名。
文件夹中可以包含以 built.js 或 built.css 为名称的文件，这俩文件会被直接嵌入到 html 页面中。

具体关于 js 中写扩展的方式可以参照 templates/js/common.js 中的方法，
全部扩展会被管理在 window.note.extensions 这个 js 变量中
*/

package extensions

import (
	"os"
	"path"
	"path/filepath"
	"rellwnote/core/config"
)

type Extension struct {
	Name       string
	BuiltinJS  string
	BuiltinCSS string
}

// LoadAll 会读取全部可安装扩展的目录下的全部扩展
func LoadAll() (res []Extension) {
	load := func(dir string) {
		dir, err := filepath.Abs(dir)
		if err != nil {
			return
		}
		open, err := os.Open(dir)
		if err != nil {
			return
		}
		readDir, err := open.ReadDir(-1)
		if err != nil {
			return
		}
		for _, i := range readDir {
			if !i.IsDir() {
				continue
			}
			fullPath := filepath.ToSlash(path.Join(dir, i.Name()))
			e, err := Load(fullPath)
			if err != nil {
				println(err.Error())
				continue
			}
			res = append(res, e)
		}
	}
	load(path.Join(config.LibraryPath, config.ExtensionDir))
	load(path.Join(config.BaseDir, config.ExtensionDir))
	return
}

func Load(extPath string) (ext Extension, err error) {
	extPath, _ = filepath.Abs(extPath)

	ext.Name = filepath.Base(extPath)
	js, jsErr := os.ReadFile(path.Join(extPath, config.BuiltinExtensionFileName+".js"))
	if jsErr == nil {
		ext.BuiltinJS = string(js)
	}
	css, cssErr := os.ReadFile(path.Join(extPath, config.BuiltinExtensionFileName+".css"))
	if cssErr == nil {
		ext.BuiltinCSS = string(css)
	}
	return
}
